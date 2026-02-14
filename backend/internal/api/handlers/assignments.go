package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type AssignmentsHandler struct {
	DB *sql.DB
}

func NewAssignmentsHandler(db *sql.DB) *AssignmentsHandler {
	return &AssignmentsHandler{DB: db}
}

type AssignRequest struct {
	WorkerIDs  []string `json:"workerIds"`
	DeviceIDs  []string `json:"deviceIds"`
	ProjectIDs []string `json:"projectIds"`
}

func (h *AssignmentsHandler) AssignWorkers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId := vars["projectId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Assignments] ERROR decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Assignments] Syncing workers %v for project %s", req.WorkerIDs, projectId)

	// Step 0: Get Project Tenant to enforce ownership
	var projectTenantID string
	if err := h.DB.QueryRow("SELECT tenant_id FROM projects WHERE project_id = ?", projectId).Scan(&projectTenantID); err != nil {
		log.Printf("[Assignments] ERROR fetching project tenant: %v", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Step 1: Unassign workers currently on this project (Safe key constraint usually calls for this, but simplistic approach here)
	// We should technically only unassign workers causing conflicts, but the frontend sends the *full* new list usually?
	// Wait, the frontend sends *Ids to add*? Or the *complete state*?
	// AssignmentsHandler usually implies "Set these as the workers".
	// But `WorkerList` "Assign" button implies "Add this worker to project".
	// `WorkerAssignProject.vue` does a single update.
	// `ProjectAssignWorkers.vue` does a bulk assignment.
	// Let's assume this endpoint sets the list.
	// But to be safe against cross-tenant unassignment, we add tenant_id check.

	// Actually, clearing *all* workers with `current_project_id = projectId` is safe IF we assume `projectId` is owned by the caller.
	// But better to be explicit.
	_, err := h.DB.Exec("UPDATE users SET current_project_id = NULL WHERE current_project_id = ? AND tenant_id = ?", projectId, projectTenantID)
	if err != nil {
		log.Printf("[Assignments] ERROR clearing old assignments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Assign new workers
	// Enforce that the worker being assigned ALSO belongs to the same tenant as the project
	stmt, err := h.DB.Prepare("UPDATE users SET current_project_id = ? WHERE user_id = ? AND tenant_id = ?")
	if err != nil {
		log.Printf("[Assignments] ERROR preparing stmt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, workerId := range req.WorkerIDs {
		res, err := stmt.Exec(projectId, workerId, projectTenantID)
		if err != nil {
			log.Printf("[Assignments] ERROR assigning worker %s: %v", workerId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			log.Printf("[Assignments] Warning: Worker %s not updated. Possible tenant mismatch or invalid ID.", workerId)
		} else {
			log.Printf("[Assignments] Worker %s assigned to project %s", workerId, projectId)
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "project": projectId})
}

func (h *AssignmentsHandler) AssignDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["siteId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Assignments] ERROR decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Assignments] Syncing devices %v for site %s", req.DeviceIDs, siteId)

	// Step 1: Unassign devices currently on this site
	_, err := h.DB.Exec("UPDATE devices SET site_id = NULL WHERE site_id = ?", siteId)
	if err != nil {
		log.Printf("[Assignments] ERROR clearing old device assignments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Assign new devices
	stmt, err := h.DB.Prepare("UPDATE devices SET site_id = ? WHERE device_id = ?")
	if err != nil {
		log.Printf("[Assignments] ERROR preparing device stmt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, devId := range req.DeviceIDs {
		res, err := stmt.Exec(siteId, devId)
		if err != nil {
			log.Printf("[Assignments] ERROR assigning device %s: %v", devId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		log.Printf("[Assignments] Device %s assignment result: %d rows affected", devId, rows)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "site": siteId})
}

func (h *AssignmentsHandler) AssignProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["siteId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[Assignments] ERROR decoding request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Assignments] Syncing projects %v for site %s", req.ProjectIDs, siteId)

	// Step 1: Unassign projects currently on this site
	_, err := h.DB.Exec("UPDATE projects SET site_id = NULL WHERE site_id = ?", siteId)
	if err != nil {
		log.Printf("[Assignments] ERROR clearing old project assignments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Assign new projects
	stmt, err := h.DB.Prepare("UPDATE projects SET site_id = ? WHERE project_id = ?")
	if err != nil {
		log.Printf("[Assignments] ERROR preparing project stmt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, projId := range req.ProjectIDs {
		res, err := stmt.Exec(siteId, projId)
		if err != nil {
			log.Printf("[Assignments] ERROR assigning project %s: %v", projId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		log.Printf("[Assignments] Project %s assignment result: %d rows affected", projId, rows)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "site": siteId})
}

func (h *AssignmentsHandler) AssignDevicesToTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantId := vars["tenantId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update devices table. Set site_id to NULL to avoid invalid site references.
	// Also set status to 'offline' (or 'online') so it shows up in the list (since we filter out 'inactive')
	log.Printf("[Assignments] Assigning devices %v to tenant %s", req.DeviceIDs, tenantId)
	stmt, err := h.DB.Prepare("UPDATE devices SET tenant_id = ?, site_id = NULL, status = 'offline' WHERE device_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, devId := range req.DeviceIDs {
		res, err := stmt.Exec(tenantId, devId)
		if err != nil {
			log.Printf("[Assignments] ERROR updating device %s: %v", devId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		log.Printf("[Assignments] Device %s updated. Rows affected: %d", devId, rows)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "tenant": tenantId})
}
