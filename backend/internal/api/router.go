package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/adapters/repository/mysql"
	"sgbuildex/internal/api/handlers"
	"sgbuildex/internal/api/middleware"
	"sgbuildex/internal/core/services"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *mux.Router, db *sql.DB, bridgeSyncHandler *handlers.BridgeSyncHandler) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("ROUTER DEBUG: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})
	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.UserScopeMiddleware)

	// --- Handlers ---
	// Auth Module Wiring
	userRepo := mysql.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Workers Module Wiring
	workerRepo := mysql.NewWorkerRepository(db)
	workerService := services.NewWorkerService(workerRepo)
	workersHandler := handlers.NewWorkersHandler(workerService)

	// Projects Module Wiring
	projectRepo := mysql.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectsHandler := handlers.NewProjectsHandler(projectService)

	// Sites Module Wiring
	siteRepo := mysql.NewSiteRepository(db)
	siteService := services.NewSiteService(siteRepo)
	sitesHandler := handlers.NewSitesHandler(siteService)

	// Devices Module Wiring
	deviceRepo := mysql.NewDeviceRepository(db)
	deviceService := services.NewDeviceService(deviceRepo)
	devicesHandler := handlers.NewDevicesHandler(deviceService)

	// Users Module Wiring
	usersService := services.NewUserService(userRepo)
	usersHandler := handlers.NewUsersHandler(usersService, db)

	assignmentsHandler := handlers.NewAssignmentsHandler(workerService, deviceService, projectService)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	// Attendance Module Wiring
	attendanceRepo := mysql.NewAttendanceRepository(db)
	attendanceService := services.NewAttendanceService(attendanceRepo, workerRepo, deviceRepo)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)

	// Settings Module Wiring
	settingsRepo := mysql.NewMySQLSettingsRepository(db)
	settingsService := services.NewSettingsService(settingsRepo)
	settingsHandler := handlers.NewSettingsHandler(settingsService)

	// --- Auth Routes ---
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// --- Workers Routes ---
	api.HandleFunc("/workers", workersHandler.GetWorkers).Methods("GET")
	api.HandleFunc("/workers", workersHandler.CreateWorker).Methods("POST")
	api.HandleFunc("/workers/{id}", workersHandler.GetWorkerById).Methods("GET")
	api.HandleFunc("/workers/{id}", workersHandler.UpdateWorker).Methods("PUT")
	api.HandleFunc("/workers/{id}", workersHandler.DeleteWorker).Methods("DELETE")

	// --- Projects Routes ---
	api.HandleFunc("/projects", projectsHandler.GetProjects).Methods("GET")
	api.HandleFunc("/projects", projectsHandler.CreateProject).Methods("POST")
	api.HandleFunc("/projects/{id}", projectsHandler.GetProjectById).Methods("GET")
	api.HandleFunc("/projects/{id}", projectsHandler.UpdateProject).Methods("PUT")
	api.HandleFunc("/projects/{id}", projectsHandler.DeleteProject).Methods("DELETE")

	// --- Sites Routes ---
	api.HandleFunc("/sites", sitesHandler.GetSites).Methods("GET")
	api.HandleFunc("/sites", sitesHandler.CreateSite).Methods("POST")
	api.HandleFunc("/sites/{id}", sitesHandler.GetSiteById).Methods("GET")
	api.HandleFunc("/sites/{id}", sitesHandler.UpdateSite).Methods("PUT")
	api.HandleFunc("/sites/{id}", sitesHandler.DeleteSite).Methods("DELETE")

	// --- Devices Routes ---
	api.HandleFunc("/devices", devicesHandler.GetDevices).Methods("GET")
	api.HandleFunc("/devices", devicesHandler.CreateDevice).Methods("POST")
	api.HandleFunc("/devices/{id}", devicesHandler.GetDeviceById).Methods("GET")
	api.HandleFunc("/devices/{id}", devicesHandler.UpdateDevice).Methods("PUT")
	api.HandleFunc("/devices/{id}", devicesHandler.DeleteDevice).Methods("DELETE")

	// --- Users Routes ---
	api.HandleFunc("/users", usersHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users", usersHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", usersHandler.GetUserById).Methods("GET")
	api.HandleFunc("/users/{id}", usersHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", usersHandler.DeleteUser).Methods("DELETE")

	// --- Attendance Routes ---
	api.HandleFunc("/attendance", attendanceHandler.GetAttendance).Methods("GET")

	// --- Uploads ---
	api.HandleFunc("/upload/face", handlers.UploadFaceHandler).Methods("POST")

	// --- Assignments Routes ---
	api.HandleFunc("/projects/{projectId}/assign-workers", assignmentsHandler.AssignWorkers).Methods("POST")
	api.HandleFunc("/sites/{siteId}/assign-devices", assignmentsHandler.AssignDevices).Methods("POST")

	api.HandleFunc("/sites/{siteId}/assign-projects", assignmentsHandler.AssignProjects).Methods("POST")

	api.HandleFunc("/users/{userId}/devices/bulk", assignmentsHandler.AssignDevicesToUser).Methods("POST")

	// --- Analytics Routes ---
	api.HandleFunc("/analytics/dashboard", analyticsHandler.GetDashboardStats).Methods("GET")
	api.HandleFunc("/analytics/activity-log", analyticsHandler.GetActivityLog).Methods("GET")
	api.HandleFunc("/analytics/detailed", analyticsHandler.GetDetailedAnalytics).Methods("GET")

	// --- Settings Routes ---
	api.HandleFunc("/settings", settingsHandler.GetSettings).Methods("GET")
	api.HandleFunc("/settings", settingsHandler.UpdateSettings).Methods("PUT")

	// --- Bridge Sync Routes ---
	if bridgeSyncHandler != nil {
		api.HandleFunc("/bridge/sync-users", bridgeSyncHandler.SyncUsers).Methods("POST")
	}
}
