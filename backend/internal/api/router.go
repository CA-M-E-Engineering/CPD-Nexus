package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/api/handlers"
	"sgbuildex/internal/api/middleware"

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
}

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *mux.Router, cfg RouterConfig) {
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

	// --- Auth Routes ---
	api.HandleFunc("/auth/login", cfg.AuthHandler.Login).Methods("POST")
	api.HandleFunc("/auth/me", cfg.AuthHandler.Me).Methods("GET")

	// --- Workers Routes ---
	api.HandleFunc("/workers", cfg.WorkersHandler.GetWorkers).Methods("GET")
	api.HandleFunc("/workers", cfg.WorkersHandler.CreateWorker).Methods("POST")
	api.HandleFunc("/workers/{id}", cfg.WorkersHandler.GetWorkerById).Methods("GET")
	api.HandleFunc("/workers/{id}", cfg.WorkersHandler.UpdateWorker).Methods("PUT")
	api.HandleFunc("/workers/{id}", cfg.WorkersHandler.DeleteWorker).Methods("DELETE")

	// --- Projects Routes ---
	api.HandleFunc("/projects", cfg.ProjectsHandler.GetProjects).Methods("GET")
	api.HandleFunc("/projects", cfg.ProjectsHandler.CreateProject).Methods("POST")
	api.HandleFunc("/projects/{id}", cfg.ProjectsHandler.GetProjectById).Methods("GET")
	api.HandleFunc("/projects/{id}", cfg.ProjectsHandler.UpdateProject).Methods("PUT")
	api.HandleFunc("/projects/{id}", cfg.ProjectsHandler.DeleteProject).Methods("DELETE")

	// --- Sites Routes ---
	api.HandleFunc("/sites", cfg.SitesHandler.GetSites).Methods("GET")
	api.HandleFunc("/sites", cfg.SitesHandler.CreateSite).Methods("POST")
	api.HandleFunc("/sites/{id}", cfg.SitesHandler.GetSiteById).Methods("GET")
	api.HandleFunc("/sites/{id}", cfg.SitesHandler.UpdateSite).Methods("PUT")
	api.HandleFunc("/sites/{id}", cfg.SitesHandler.DeleteSite).Methods("DELETE")

	// --- Devices Routes ---
	api.HandleFunc("/devices", cfg.DevicesHandler.GetDevices).Methods("GET")
	api.HandleFunc("/devices", cfg.DevicesHandler.CreateDevice).Methods("POST")
	api.HandleFunc("/devices/{id}", cfg.DevicesHandler.GetDeviceById).Methods("GET")
	api.HandleFunc("/devices/{id}", cfg.DevicesHandler.UpdateDevice).Methods("PUT")
	api.HandleFunc("/devices/{id}", cfg.DevicesHandler.DeleteDevice).Methods("DELETE")

	// --- Users Routes ---
	api.HandleFunc("/users", cfg.UsersHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users", cfg.UsersHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", cfg.UsersHandler.GetUserById).Methods("GET")
	api.HandleFunc("/users/{id}", cfg.UsersHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", cfg.UsersHandler.DeleteUser).Methods("DELETE")

	// --- Attendance Routes ---
	api.HandleFunc("/attendance", cfg.AttendanceHandler.GetAttendance).Methods("GET")

	// --- Uploads ---
	api.HandleFunc("/upload/face", handlers.UploadFaceHandler).Methods("POST")

	// --- Assignments Routes ---
	api.HandleFunc("/projects/{projectId}/assign-workers", cfg.AssignmentsHandler.AssignWorkers).Methods("POST")
	api.HandleFunc("/sites/{siteId}/assign-devices", cfg.AssignmentsHandler.AssignDevices).Methods("POST")

	api.HandleFunc("/sites/{siteId}/assign-projects", cfg.AssignmentsHandler.AssignProjects).Methods("POST")

	api.HandleFunc("/users/{userId}/devices/bulk", cfg.AssignmentsHandler.AssignDevicesToUser).Methods("POST")

	// --- Analytics Routes ---
	api.HandleFunc("/analytics/dashboard", cfg.AnalyticsHandler.GetDashboardStats).Methods("GET")
	api.HandleFunc("/analytics/activity-log", cfg.AnalyticsHandler.GetActivityLog).Methods("GET")
	api.HandleFunc("/analytics/detailed", cfg.AnalyticsHandler.GetDetailedAnalytics).Methods("GET")

	// --- Settings Routes ---
	api.HandleFunc("/settings", cfg.SettingsHandler.GetSettings).Methods("GET")
	api.HandleFunc("/settings", cfg.SettingsHandler.UpdateSettings).Methods("PUT")

	// --- Bridge Sync Routes ---
	if cfg.BridgeSyncHandler != nil {
		api.HandleFunc("/bridge/sync-users", cfg.BridgeSyncHandler.SyncUsers).Methods("POST")
	}

	// Serve Static Files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
}
