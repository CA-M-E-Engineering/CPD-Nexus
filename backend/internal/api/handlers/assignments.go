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

	// Step 0: Get Project User to enforce ownership
	var projectUserID string
	if err := h.DB.QueryRow("SELECT user_id FROM projects WHERE project_id = ?", projectId).Scan(&projectUserID); err != nil {
		log.Printf("[Assignments] ERROR fetching project user: %v", err)
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Step 1: Unassign workers currently on this project
	_, err := h.DB.Exec("UPDATE workers SET current_project_id = NULL WHERE current_project_id = ? AND user_id = ?", projectId, projectUserID)
	if err != nil {
		log.Printf("[Assignments] ERROR clearing old assignments: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Assign new workers
	stmt, err := h.DB.Prepare("UPDATE workers SET current_project_id = ? WHERE worker_id = ? AND user_id = ? AND status = 'active'")
	if err != nil {
		log.Printf("[Assignments] ERROR preparing stmt: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, workerId := range req.WorkerIDs {
		res, err := stmt.Exec(projectId, workerId, projectUserID)
		if err != nil {
			log.Printf("[Assignments] ERROR assigning worker %s: %v", workerId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			log.Printf("[Assignments] Warning: Worker %s not updated. Possible user mismatch or invalid ID.", workerId)
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

func (h *AssignmentsHandler) AssignDevicesToUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	var req AssignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Assignments] Assigning devices %v to user %s", req.DeviceIDs, userId)
	stmt, err := h.DB.Prepare("UPDATE devices SET user_id = ?, site_id = NULL, status = 'offline' WHERE device_id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, devId := range req.DeviceIDs {
		res, err := stmt.Exec(userId, devId)
		if err != nil {
			log.Printf("[Assignments] ERROR updating device %s: %v", devId, err)
			continue
		}
		rows, _ := res.RowsAffected()
		log.Printf("[Assignments] Device %s updated. Rows affected: %d", devId, rows)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "assigned", "user": userId})
}
