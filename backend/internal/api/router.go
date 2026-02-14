package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/adapters/repository/mysql"
	"sgbuildex/internal/api/handlers"
	"sgbuildex/internal/core/services"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *mux.Router, db *sql.DB) {
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

	// --- Handlers ---
	// Auth Module Wiring
	tenantRepo := mysql.NewTenantRepository(db)
	authService := services.NewAuthService(tenantRepo)
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

	assignmentsHandler := handlers.NewAssignmentsHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)
	tenantsHandler := handlers.NewTenantsHandler(db)
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
	// api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST") // Stateless JWT usually handles logout on client

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

	// --- Tenants Routes ---
	api.HandleFunc("/tenants", tenantsHandler.GetTenants).Methods("GET")
	api.HandleFunc("/tenants", tenantsHandler.CreateTenant).Methods("POST")
	api.HandleFunc("/tenants/{id}", tenantsHandler.GetTenantById).Methods("GET")
	api.HandleFunc("/tenants/{id}", tenantsHandler.UpdateTenant).Methods("PUT")
	api.HandleFunc("/tenants/{id}", tenantsHandler.DeleteTenant).Methods("DELETE")

	// --- Attendance Routes ---
	api.HandleFunc("/attendance", attendanceHandler.GetAttendance).Methods("GET")

	// --- Assignments Routes ---
	api.HandleFunc("/projects/{projectId}/assign-workers", assignmentsHandler.AssignWorkers).Methods("POST")
	api.HandleFunc("/sites/{siteId}/assign-devices", assignmentsHandler.AssignDevices).Methods("POST")

	api.HandleFunc("/sites/{siteId}/assign-projects", assignmentsHandler.AssignProjects).Methods("POST")

	api.HandleFunc("/tenants/{tenantId}/devices/bulk", assignmentsHandler.AssignDevicesToTenant).Methods("POST")

	// --- Analytics Routes ---
	api.HandleFunc("/analytics/dashboard", analyticsHandler.GetDashboardStats).Methods("GET")
	api.HandleFunc("/analytics/activity-log", analyticsHandler.GetActivityLog).Methods("GET")
	api.HandleFunc("/analytics/detailed", analyticsHandler.GetDetailedAnalytics).Methods("GET")

	// --- Settings Routes ---
	api.HandleFunc("/settings", settingsHandler.GetSettings).Methods("GET")
	api.HandleFunc("/settings", settingsHandler.UpdateSettings).Methods("PUT")
}
