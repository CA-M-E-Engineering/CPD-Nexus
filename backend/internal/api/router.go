package api

import (
	"encoding/json"
	"net/http"
	"cpd-nexus/internal/api/handlers"
	"cpd-nexus/internal/api/middleware"
	"cpd-nexus/internal/core/ports"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	AuthHandler        *handlers.AuthHandler
	WorkersHandler     *handlers.WorkersHandler
	ProjectsHandler    *handlers.ProjectsHandler
	SitesHandler       *handlers.SitesHandler
	DevicesHandler     *handlers.DevicesHandler
	UsersHandler       *handlers.UsersHandler
	AssignmentsHandler *handlers.AssignmentsHandler
	AnalyticsHandler   *handlers.AnalyticsHandler
	AttendanceHandler  *handlers.AttendanceHandler
	SettingsHandler    *handlers.SettingsHandler
	BridgeSyncHandler  *handlers.BridgeSyncHandler
	PitstopHandler     *handlers.PitstopHandler
	UserRepo           ports.UserRepository
}

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *mux.Router, cfg RouterConfig) {

	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	// --- Auth Routes (Public) ---
	r.HandleFunc("/api/auth/login", cfg.AuthHandler.Login).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.UserScopeMiddleware(cfg.UserRepo))

	// --- Auth Routes (Protected) ---
	api.HandleFunc("/auth/me", cfg.AuthHandler.Me).Methods("GET")

	// --- Administrative Routes (Global Admin / Vendor Only) ---
	admin := api.PathPrefix("").Subrouter()
	admin.Use(middleware.RequireAdminScope)

	admin.HandleFunc("/users", cfg.UsersHandler.GetUsers).Methods("GET")
	admin.HandleFunc("/users", cfg.UsersHandler.CreateUser).Methods("POST")
	admin.HandleFunc("/users/{id}", cfg.UsersHandler.GetUserById).Methods("GET")
	admin.HandleFunc("/users/{id}", cfg.UsersHandler.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", cfg.UsersHandler.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/users/{id}/bridge", cfg.UsersHandler.UpdateBridgeConfig).Methods("PUT")

	admin.HandleFunc("/devices", cfg.DevicesHandler.CreateDevice).Methods("POST")

	if cfg.PitstopHandler != nil {

		admin.HandleFunc("/pitstop/authorisations/sync", cfg.PitstopHandler.SyncConfig).Methods("POST")
		admin.HandleFunc("/users/{id}/pitstop-on-behalf-of", cfg.PitstopHandler.AssignOnBehalfOf).Methods("POST")
	}

	// --- Scoped Routes (Project Isolation) ---
	scoped := api.PathPrefix("").Subrouter()
	scoped.Use(middleware.RequireUserScope)

	// --- Workers Routes ---
	scoped.HandleFunc("/workers", cfg.WorkersHandler.GetWorkers).Methods("GET")
	scoped.HandleFunc("/workers", cfg.WorkersHandler.CreateWorker).Methods("POST")
	scoped.HandleFunc("/workers/{id}", cfg.WorkersHandler.GetWorkerById).Methods("GET")
	scoped.HandleFunc("/workers/{id}", cfg.WorkersHandler.UpdateWorker).Methods("PUT")
	scoped.HandleFunc("/workers/{id}", cfg.WorkersHandler.DeleteWorker).Methods("DELETE")

	// --- Projects Routes ---
	scoped.HandleFunc("/projects", cfg.ProjectsHandler.GetProjects).Methods("GET")
	scoped.HandleFunc("/projects", cfg.ProjectsHandler.CreateProject).Methods("POST")
	scoped.HandleFunc("/projects/{id}", cfg.ProjectsHandler.GetProjectById).Methods("GET")
	scoped.HandleFunc("/projects/{id}", cfg.ProjectsHandler.UpdateProject).Methods("PUT")
	scoped.HandleFunc("/projects/{id}", cfg.ProjectsHandler.DeleteProject).Methods("DELETE")

	// --- Sites Routes ---
	scoped.HandleFunc("/sites", cfg.SitesHandler.GetSites).Methods("GET")
	scoped.HandleFunc("/sites", cfg.SitesHandler.CreateSite).Methods("POST")
	scoped.HandleFunc("/sites/{id}", cfg.SitesHandler.GetSiteById).Methods("GET")
	scoped.HandleFunc("/sites/{id}", cfg.SitesHandler.UpdateSite).Methods("PUT")
	scoped.HandleFunc("/sites/{id}", cfg.SitesHandler.DeleteSite).Methods("DELETE")

	// --- Devices Routes (Scoped) ---
	scoped.HandleFunc("/devices", cfg.DevicesHandler.GetDevices).Methods("GET")
	scoped.HandleFunc("/devices/{id}", cfg.DevicesHandler.GetDeviceById).Methods("GET")
	scoped.HandleFunc("/devices/{id}", cfg.DevicesHandler.UpdateDevice).Methods("PUT")
	scoped.HandleFunc("/devices/{id}", cfg.DevicesHandler.DeleteDevice).Methods("DELETE")

	// --- Attendance Routes ---
	scoped.HandleFunc("/attendance", cfg.AttendanceHandler.GetAttendance).Methods("GET")
	scoped.HandleFunc("/attendance/{id}", cfg.AttendanceHandler.UpdateAttendance).Methods("PUT")

	// --- Uploads ---
	scoped.HandleFunc("/upload/face", handlers.UploadFaceHandler).Methods("POST")

	// --- Assignments Routes ---
	scoped.HandleFunc("/projects/{projectId}/assign-workers", cfg.AssignmentsHandler.AssignWorkers).Methods("POST")
	scoped.HandleFunc("/sites/{siteId}/assign-devices", cfg.AssignmentsHandler.AssignDevices).Methods("POST")
	scoped.HandleFunc("/sites/{siteId}/assign-projects", cfg.AssignmentsHandler.AssignProjects).Methods("POST")
	scoped.HandleFunc("/users/{userId}/devices/bulk", cfg.AssignmentsHandler.AssignDevicesToUser).Methods("POST")

	// --- Analytics Routes ---
	scoped.HandleFunc("/analytics/dashboard", cfg.AnalyticsHandler.GetDashboardStats).Methods("GET")
	scoped.HandleFunc("/analytics/activity-log", cfg.AnalyticsHandler.GetActivityLog).Methods("GET")
	scoped.HandleFunc("/analytics/detailed", cfg.AnalyticsHandler.GetDetailedAnalytics).Methods("GET")

	// --- Settings Routes ---
	scoped.HandleFunc("/settings", cfg.SettingsHandler.GetSettings).Methods("GET")
	scoped.HandleFunc("/settings", cfg.SettingsHandler.UpdateSettings).Methods("PUT")

	// --- Bridge Sync Routes ---
	if cfg.BridgeSyncHandler != nil {
		scoped.HandleFunc("/bridge/sync-users", cfg.BridgeSyncHandler.SyncUsers).Methods("POST")
	}

	// --- Pitstop Test Endpoints (Scoped/Admin) ---
	if cfg.PitstopHandler != nil {
		scoped.HandleFunc("/pitstop/authorisations", cfg.PitstopHandler.GetAuthorisations).Methods("GET")
		scoped.HandleFunc("/pitstop/authorisations/testing-projects", cfg.PitstopHandler.GetTestingProjects).Methods("GET")
		scoped.HandleFunc("/pitstop/authorisations/test-submission/{project_id}", cfg.PitstopHandler.TestSubmission).Methods("POST")
	}

	// Serve Static Files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
}
