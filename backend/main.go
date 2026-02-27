package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"sgbuildex/internal/adapters/external/sgbuildex"
	"sgbuildex/internal/adapters/repository/mysql"
	"sgbuildex/internal/api"
	apiHandlers "sgbuildex/internal/api/handlers"
	"sgbuildex/internal/bridge"
	bridgeHandlers "sgbuildex/internal/bridge/handlers"
	"sgbuildex/internal/core/domain"
	"sgbuildex/internal/core/services"
	"sgbuildex/internal/pkg/config"
	"sgbuildex/internal/pkg/logger"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg := config.LoadConfig()
	logger.Infof("--- CPD Nexus Unified Backend Starting ---")

	// --- 1. DB Connection ---
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		logger.Errorf("Failed to connect to DB: %v", err)
		os.Exit(1)
	}
	defer db.Close()

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

	// Services
	workerService := services.NewWorkerService(workerRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo, workerRepo, deviceRepo)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	siteService := services.NewSiteService(siteRepo)
	projectService := services.NewProjectService(projectRepo)
	deviceService := services.NewDeviceService(deviceRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	var settingsService *services.SettingsService

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
		// SettingsHandler will be added later after Schedulers are ready
	}

	// Bridge Integration
	requestMgr := bridge.NewRequestManager(db)
	userSyncBuilder := bridgeHandlers.NewUserSyncBuilder(workerService, workerRepo, deviceRepo)
	routerCfg.BridgeSyncHandler = apiHandlers.NewBridgeSyncHandler(userSyncBuilder, requestMgr)

	attendanceHandler := bridgeHandlers.NewAttendanceHandler(attendanceService)
	requestMgr.RegisterHandler("GET_ATTENDANCE_RESPONSE", attendanceHandler)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- 3. Component B: Submission Workers & Schedulers ---
	client := sgbuildex.NewClient(cfg.IngressURL)

	// Task 1: Attendance Sync (Bridge -> Nexus)
	syncTask := func(taskCtx context.Context) {
		logger.Infof("[AttendanceSync] Starting scheduled bridge fetch...")
		if err := requestMgr.RequestAttendance(); err != nil {
			logger.Errorf("[AttendanceSync] Bridge fetch failed: %v", err)
		} else {
			logger.Infof("[AttendanceSync] Fetch requests sent to bridge.")
		}
	}

	// Task 2: CPD Submission (Nexus -> SGBuildex)
	submitTask := func(taskCtx context.Context) {
		logger.Infof("[CPDSubmission] Starting scheduled submission cycle...")

		// 0. Fetch latest settings for limits
		settings, err := settingsRepo.GetSettings(taskCtx)
		if err != nil {
			logger.Errorf("[CPDSubmission] Failed to fetch settings: %v. Using defaults.", err)
			settings = &domain.SystemSettings{
				MaxWorkersPerRequest: 100,
				MaxPayloadSizeKB:     256,
				MaxRequestsPerMinute: 150,
			}
		}

		// 1. Manpower Distribution
		distRows, err := attendanceRepo.ExtractMonthlyDistributionData(taskCtx)
		if err == nil {
			mdPayloads := sgbuildex.MapAggregationToDistribution(distRows)
			mdSubmittables := make([]sgbuildex.ManpowerDistributionWrapper, len(mdPayloads))
			for i, p := range mdPayloads {
				mdSubmittables[i] = sgbuildex.ManpowerDistributionWrapper{ManpowerDistribution: p}
			}
			sgbuildex.SubmitPayloads(taskCtx, submissionRepo, client, settings, mdSubmittables)
		}

		// 2. Manpower Utilization
		rows, err := attendanceRepo.ExtractPendingAttendance(taskCtx)
		if err == nil && len(rows) > 0 {
			muPayloads := sgbuildex.MapAttendanceToManpower(rows)
			muSubmittables := make([]sgbuildex.ManpowerUtilizationWrapper, len(muPayloads))
			for i, p := range muPayloads {
				muSubmittables[i] = sgbuildex.ManpowerUtilizationWrapper{ManpowerUtilization: p}
			}
			sgbuildex.SubmitPayloads(taskCtx, submissionRepo, client, settings, muSubmittables)
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
	settingsService = services.NewSettingsService(settingsRepo, attendanceSyncScheduler, cpdSubmissionScheduler)
	routerCfg.SettingsHandler = apiHandlers.NewSettingsHandler(settingsService)

	// --- 4. Component C: REST API ---
	go startAPI(cfg, routerCfg)

	// --- 5. Component D: Core Loops ---
	go startBridge(ctx, cfg, db, requestMgr, userSyncBuilder)
	go attendanceSyncScheduler.Start(ctx)
	go cpdSubmissionScheduler.Start(ctx)

	logger.Infof("[System] Schedulers and API services fully operational")

	// --- 6. Wait for Shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	logger.Infof("Shutting down CPD Nexus unified backend...")
	cancel()

	// Give a small grace period for goroutines to clean up
	time.Sleep(1 * time.Second)
	logger.Infof("Final shutdown complete.")
}

func startAPI(cfg *config.Config, routerCfg api.RouterConfig) {
	router := mux.NewRouter()
	api.RegisterRoutes(router, routerCfg)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", "http://127.0.0.1:5173",
			"http://localhost:5174", "http://127.0.0.1:5174",
			"http://localhost:5175", "http://127.0.0.1:5175",
			"http://localhost:5176", "http://127.0.0.1:5176",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-User-ID"},
		AllowCredentials: true,
	})

	server := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: c.Handler(router),
	}

	logger.Infof("[API] Starting REST server on port %s", cfg.APIPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("[API] Server failed: %v", err)
	}
}

func startBridge(ctx context.Context, cfg *config.Config, db *sql.DB, requestMgr *bridge.RequestManager, userSyncBuilder *bridgeHandlers.UserSyncBuilder) {
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
				rows, err := db.QueryContext(ctx, "SELECT user_id, bridge_ws_url FROM users WHERE bridge_status = 'active' AND bridge_ws_url IS NOT NULL")
				if err != nil {
					logger.Errorf("[Bridge] Failed to fetch active bridges: %v", err)
				} else {
					activeIDs := make(map[string]bool)
					for rows.Next() {
						var userID string
						var wsURL string
						if err := rows.Scan(&userID, &wsURL); err == nil {
							activeIDs[userID] = true
							transport, exists := requestMgr.GetTransport(userID)

							if !exists {
								logger.Infof("[Bridge] Creating new connection for user %s to %s", userID, wsURL)
								t := bridge.NewTransport(wsURL)
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
					}
					rows.Close()

					// Remove any transports that are no longer active
					for id, t := range requestMgr.GetAllTransports() {
						if !activeIDs[id] {
							logger.Infof("[Bridge] Removing inactive connection for user %s", id)
							t.Close()
							requestMgr.RemoveTransport(id)
						}
					}
				}
				time.Sleep(10 * time.Second)
			}
		}
	}()

	// Command scheduler
	requestInterval := time.Duration(cfg.BridgeIntervalSeconds) * time.Second
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
