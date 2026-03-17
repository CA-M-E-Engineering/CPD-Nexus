package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cpd-nexus/internal/adapters/external/sgbuildex"
	"cpd-nexus/internal/adapters/repository/mysql"
	"cpd-nexus/internal/api"
	apiHandlers "cpd-nexus/internal/api/handlers"
	"cpd-nexus/internal/api/middleware"
	"cpd-nexus/internal/bridge"
	bridgeHandlers "cpd-nexus/internal/bridge/handlers"
	"cpd-nexus/internal/core/domain"
	"cpd-nexus/internal/core/ports"
	"cpd-nexus/internal/core/services"
	"cpd-nexus/internal/pkg/config"
	"cpd-nexus/internal/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg := config.LoadConfig()
	logger.Infof("--- CPD Nexus Unified Backend Starting ---")

	// --- 0. Configure JWT middleware ---
	middleware.SetJWTSecret(cfg.JWTSecret)

	// --- 1. DB Connection ---
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		logger.Errorf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Configure connection pool (#23)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		logger.Errorf("Failed to ping DB: %v", err)
		os.Exit(1)
	}

	// --- 2. Shared Initialization ---
	// Repositories
	attendanceRepo := mysql.NewAttendanceRepository(db)
	workerRepo := mysql.NewWorkerRepository(db)
	deviceRepo := mysql.NewDeviceRepository(db)
	settingsRepo := mysql.NewMySQLSettingsRepository(db)
	submissionRepo := mysql.NewSubmissionRepository(db)
	userRepo := mysql.NewUserRepository(db)
	siteRepo := mysql.NewSiteRepository(db)
	projectRepo := mysql.NewProjectRepository(db)
	analyticsRepo := mysql.NewAnalyticsRepository(db)
	pitstopRepo := mysql.NewPitstopRepository(db)

	// Services
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	workerService := services.NewWorkerService(workerRepo, analyticsService)
	attendanceService := services.NewAttendanceService(attendanceRepo, workerRepo, deviceRepo, analyticsService)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, analyticsService)
	userService := services.NewUserService(userRepo, analyticsService, cfg.DefaultUserPassword)
	siteService := services.NewSiteService(siteRepo, analyticsService)
	projectService := services.NewProjectService(projectRepo, analyticsService)
	deviceService := services.NewDeviceService(deviceRepo, analyticsService)
	var settingsService ports.SettingsService

	// Internal client for external fetch
	sgClient := sgbuildex.NewClient(cfg.IngressURL, cfg.PitstopURL)
	pitstopService := services.NewPitstopService(pitstopRepo, sgClient, attendanceRepo, submissionRepo, settingsRepo, analyticsService)

	// Handlers
	routerCfg := api.RouterConfig{
		AuthHandler:        apiHandlers.NewAuthHandler(authService, userService),
		WorkersHandler:     apiHandlers.NewWorkersHandler(workerService),
		ProjectsHandler:    apiHandlers.NewProjectsHandler(projectService),
		SitesHandler:       apiHandlers.NewSitesHandler(siteService),
		DevicesHandler:     apiHandlers.NewDevicesHandler(deviceService),
		UsersHandler:       apiHandlers.NewUsersHandler(userService),
		AssignmentsHandler: apiHandlers.NewAssignmentsHandler(workerService, deviceService, projectService),
		AnalyticsHandler:   apiHandlers.NewAnalyticsHandler(analyticsService),
		AttendanceHandler:  apiHandlers.NewAttendanceHandler(attendanceService),
		PitstopHandler:     apiHandlers.NewPitstopHandler(pitstopService),
		// SettingsHandler will be added later after Schedulers are ready
	}

	// Bridge Integration
	bridgeRepo := mysql.NewBridgeRepository(db)
	requestMgr := bridge.NewRequestManager(bridgeRepo)
	userSyncBuilder := bridgeHandlers.NewUserSyncBuilder(workerService, workerRepo, deviceRepo)
	routerCfg.BridgeSyncHandler = apiHandlers.NewBridgeSyncHandler(userSyncBuilder, requestMgr)

	attendanceHandler := bridgeHandlers.NewAttendanceHandler(attendanceService)
	requestMgr.RegisterHandler("GET_ATTENDANCE_RESPONSE", attendanceHandler)

	userSyncResponseHandler := bridgeHandlers.NewUserSyncResponseHandler(workerRepo)
	requestMgr.RegisterHandler("REGISTER_USER_RESPONSE", userSyncResponseHandler)
	requestMgr.RegisterHandler("UPDATE_USER_RESPONSE", userSyncResponseHandler)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Task 1: Attendance Sync (Bridge -> Nexus)
	syncTask := func(taskCtx context.Context) {
		logger.Infof("[AttendanceSync] Starting scheduled bridge fetch...")
		if err := requestMgr.RequestAttendance(taskCtx); err != nil {
			logger.Errorf("[AttendanceSync] Bridge fetch failed: %v", err)
		} else {
			logger.Infof("[AttendanceSync] Fetch requests sent to bridge.")
		}
	}

	// Task 2: CPD Submission (Nexus → SGBuildex) — delegated to the service layer
	submitTask := func(taskCtx context.Context) {
		logger.Infof("[CPDSubmission] Starting scheduled submission cycle...")
		if err := pitstopService.SubmitPendingAttendance(taskCtx); err != nil {
			logger.Errorf("[CPDSubmission] Submission cycle failed: %v", err)
		} else {
			logger.Infof("[CPDSubmission] Submission cycle completed.")
		}
	}

	attendanceSyncScheduler := services.NewDailyScheduler(
		settingsRepo,
		"AttendanceSync",
		func(s *domain.SystemSettings) string { return s.AttendanceSyncTime },
		syncTask,
	)
	cpdSubmissionScheduler := services.NewDailyScheduler(
		settingsRepo,
		"CPDSubmission",
		func(s *domain.SystemSettings) string { return s.CPDSubmissionTime },
		submitTask,
	)

	// Finalized Settings Service with Scheduler injection for real-time updates
	settingsService = services.NewSettingsService(settingsRepo, attendanceSyncScheduler, cpdSubmissionScheduler, analyticsService)
	routerCfg.SettingsHandler = apiHandlers.NewSettingsHandler(settingsService)

	// --- 4. Component C: REST API ---
	server := startAPI(cfg, routerCfg)

	// --- 5. Component D: Core Loops ---
	go startBridge(ctx, cfg, bridgeRepo, requestMgr, userSyncBuilder)
	go attendanceSyncScheduler.Start(ctx)
	go cpdSubmissionScheduler.Start(ctx)

	logger.Infof("[System] Schedulers and API services fully operational")

	// --- 6. Wait for Shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Infof("Shutting down CPD Nexus unified backend...")
	cancel()

	// Shutdown HTTP Server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
	}

	// Give a small grace period for goroutines to clean up
	time.Sleep(1 * time.Second)
	logger.Infof("Final shutdown complete.")
}

func startAPI(cfg *config.Config, routerCfg api.RouterConfig) *http.Server {
	router := mux.NewRouter()
	api.RegisterRoutes(router, routerCfg)

	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(cfg.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: c.Handler(router),
	}

	go func() {
		logger.Infof("[API] Starting REST server on port %s", cfg.APIPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("[API] Server failed: %v", err)
		}
	}()

	return server
}

func startBridge(ctx context.Context, cfg *config.Config, bridgeRepo ports.BridgeRepository, requestMgr *bridge.RequestManager, userSyncBuilder *bridgeHandlers.UserSyncBuilder) {
	// Connection maintenance loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, t := range requestMgr.GetAllTransports() {
					t.Close()
				}
				return
			default:
				// Fetch active bridges
				configs, err := bridgeRepo.GetActiveBridges(ctx)
				if err != nil {
					logger.Errorf("[Bridge] Failed to fetch active bridges: %v", err)
				} else {
					activeIDs := make(map[string]bool)
					for _, c := range configs {
						userID := c.UserID
						wsURL := c.WSURL
						authToken := c.AuthToken

						activeIDs[userID] = true
						transport, exists := requestMgr.GetTransport(userID)

						if !exists {
							logger.Infof("[Bridge] Creating new connection for user %s to %s", userID, wsURL)
							t := bridge.NewTransport(wsURL, authToken)
							requestMgr.AddTransport(userID, t)

							go requestMgr.HandleIncomingMessages(ctx, userID, t)

							if err := t.Connect(); err != nil {
								logger.Errorf("[Bridge] Connection failed for %s: %v", userID, err)
							}
						} else if !transport.IsConnected() {
							logger.Infof("[Bridge] Reconnecting for user %s to %s", userID, wsURL)
							if err := transport.Connect(); err != nil {
								logger.Errorf("[Bridge] Reconnection failed for %s: %v", userID, err)
							}
						}
					}

					// Remove any transports that are no longer active
					for id, t := range requestMgr.GetAllTransports() {
						if !activeIDs[id] {
							logger.Infof("[Bridge] Removing inactive connection for user %s", id)
							t.Close()
							requestMgr.RemoveTransport(id)
						}
					}
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(10 * time.Second):
				}
			}
		}
	}()

	// Command scheduler
	requestInterval := 10 * time.Second
	logger.Infof("[Bridge] Command scheduler started (Interval: %v)", requestInterval)

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Sync pending user registrations/updates on interval
			if err := requestMgr.RequestUserSync(ctx, userSyncBuilder); err != nil {
				logger.Errorf("[Bridge] User sync failed: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}
