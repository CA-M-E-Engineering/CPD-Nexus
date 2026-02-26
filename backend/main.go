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
	"sgbuildex/internal/core/ports"
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
	settingsService := services.NewSettingsService(settingsRepo)

	// Handlers
	routerCfg := api.RouterConfig{
		AuthHandler:        apiHandlers.NewAuthHandler(authService),
		WorkersHandler:     apiHandlers.NewWorkersHandler(workerService),
		ProjectsHandler:    apiHandlers.NewProjectsHandler(projectService),
		SitesHandler:       apiHandlers.NewSitesHandler(siteService),
		DevicesHandler:     apiHandlers.NewDevicesHandler(deviceService),
		UsersHandler:       apiHandlers.NewUsersHandler(userService),
		AssignmentsHandler: apiHandlers.NewAssignmentsHandler(workerService, deviceService, projectService),
		AnalyticsHandler:   apiHandlers.NewAnalyticsHandler(analyticsService),
		AttendanceHandler:  apiHandlers.NewAttendanceHandler(attendanceService),
		SettingsHandler:    apiHandlers.NewSettingsHandler(settingsService),
	}

	// Bridge Integration
	transport := bridge.NewTransport(cfg.BridgeURL)
	userSyncBuilder := bridgeHandlers.NewUserSyncBuilder(workerService, workerRepo, deviceRepo)
	routerCfg.BridgeSyncHandler = apiHandlers.NewBridgeSyncHandler(userSyncBuilder, transport)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- 3. Component A: REST API ---
	go startAPI(cfg, routerCfg)

	// --- 4. Component B: Bridge Connector ---
	go startBridge(ctx, cfg, transport, db, attendanceService, userSyncBuilder)

	// --- 5. Component C: Submission Worker ---
	go startWorker(ctx, cfg, db, settingsRepo, attendanceRepo, submissionRepo)

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

func startBridge(ctx context.Context, cfg *config.Config, transport *bridge.Transport, db *sql.DB, attendanceService ports.AttendanceService, userSyncBuilder *bridgeHandlers.UserSyncBuilder) {
	requestMgr := bridge.NewRequestManager(transport, db)

	handler := bridgeHandlers.NewAttendanceHandler(attendanceService)
	requestMgr.RegisterHandler("GET_ATTENDANCE_RESPONSE", handler)

	// Connection maintenance loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				transport.Close()
				return
			default:
				if !transport.IsConnected() {
					logger.Infof("[Bridge] Connecting to %s...", cfg.BridgeURL)
					if err := transport.Connect(); err != nil {
						logger.Errorf("[Bridge] Connection failed: %v. Retrying in 5s...", err)
						time.Sleep(5 * time.Second)
						continue
					}
				}
				time.Sleep(10 * time.Second)
			}
		}
	}()

	// Start message listener
	go requestMgr.HandleIncomingMessages(ctx)

	// Command scheduler
	requestInterval := time.Duration(cfg.BridgeIntervalSeconds) * time.Second
	logger.Infof("[Bridge] Command scheduler started (Interval: %v)", requestInterval)

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if transport.IsConnected() {
				if err := requestMgr.RequestAttendance(); err != nil {
					logger.Errorf("[Bridge] Fetch request failed: %v", err)
				}
				// Sync pending user registrations/updates
				// if err := requestMgr.RequestUserSync(ctx, userSyncBuilder); err != nil {
				// 	logger.Errorf("[Bridge] User sync failed: %v", err)
				// }
			}
		case <-ctx.Done():
			return
		}
	}
}

func startWorker(ctx context.Context, cfg *config.Config, db *sql.DB, settingsRepo ports.SettingsRepository, attendanceRepo ports.AttendanceRepository, submissionRepo ports.SubmissionRepository) {
	client := sgbuildex.NewClient(cfg.IngressURL)

	task := func(taskCtx context.Context) {
		logger.Infof("[Worker] Starting submission check...")

		// 1. Submit Manpower Distribution (Monthly Aggregate - Runs every time)
		distRows, err := attendanceRepo.ExtractMonthlyDistributionData(taskCtx)
		if err != nil {
			logger.Errorf("[Worker] MD Extraction failed: %v", err)
		} else {
			mdPayloads := sgbuildex.MapAggregationToDistribution(distRows)
			mdSubmittables := make([]sgbuildex.ManpowerDistributionWrapper, len(mdPayloads))
			for i, p := range mdPayloads {
				mdSubmittables[i] = sgbuildex.ManpowerDistributionWrapper{ManpowerDistribution: p}
			}
			if err := sgbuildex.SubmitPayloads(taskCtx, submissionRepo, client, mdSubmittables); err != nil {
				logger.Errorf("[Worker] MD Submission failed: %v", err)
			} else {
				logger.Infof("[Worker] Successfully processed %d MD aggregate records", len(mdPayloads))
			}
		}

		// 2. Submit Manpower Utilization (Incremental - Runs only if pending)
		rows, err := attendanceRepo.ExtractPendingAttendance(taskCtx)
		if err != nil {
			logger.Errorf("[Worker] MU Extraction failed: %v", err)
			return
		}
		if len(rows) == 0 {
			logger.Debugf("[Worker] No pending records for MU")
			return
		}

		muPayloads := sgbuildex.MapAttendanceToManpower(rows)
		muSubmittables := make([]sgbuildex.ManpowerUtilizationWrapper, len(muPayloads))
		for i, p := range muPayloads {
			muSubmittables[i] = sgbuildex.ManpowerUtilizationWrapper{ManpowerUtilization: p}
		}
		if err := sgbuildex.SubmitPayloads(taskCtx, submissionRepo, client, muSubmittables); err != nil {
			logger.Errorf("[Worker] MU Submission failed: %v", err)
		} else {
			logger.Infof("[Worker] Successfully processed %d MU record(s)", len(muPayloads))
		}
	}

	scheduler := services.NewCPDScheduler(settingsRepo, task)
	logger.Infof("[Worker] Scheduled submission tasks started")

	// Execute immediate check
	task(ctx)

	scheduler.Start(ctx)
}
