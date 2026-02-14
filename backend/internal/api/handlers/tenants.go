package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type TenantsHandler struct {
	DB *sql.DB
}

func NewTenantsHandler(db *sql.DB) *TenantsHandler {
	return &TenantsHandler{DB: db}
}

type Tenant struct {
	ID          string  `json:"tenant_id"`
	TenantName  string  `json:"tenant_name"`
	TenantType  string  `json:"tenant_type"`
	UEN         string  `json:"uen"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Username    string  `json:"username"`
	Status      string  `json:"status"`
	IsBCA       bool    `json:"is_bca"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lng"`
	Address     string  `json:"address"`
	WorkerCount int     `json:"worker_count"`
	DeviceCount int     `json:"device_count"`
	Password    string  `json:"password,omitempty"`
}

func (h *TenantsHandler) GetTenants(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
        SELECT 
            t.tenant_id, 
            t.tenant_name, 
            t.tenant_type, 
            c.uen,
            t.contact_email,
            t.contact_phone,
            t.username,
            t.status,
            IFNULL(c.company_type = 'contractor', 0) as is_bca,
            t.latitude,
            t.longitude,
            t.address,
            (SELECT COUNT(*) FROM users u WHERE u.tenant_id = t.tenant_id AND u.role = 'worker' AND u.status = 'active') as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.tenant_id = t.tenant_id AND d.status != 'inactive') as device_count
        FROM tenants t
        LEFT JOIN companies c ON t.tenant_id = c.tenant_id AND c.company_type = 'contractor'
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tenants := []Tenant{}
	for rows.Next() {
		var t Tenant
		var uen, email, phone, username sql.NullString
		if err := rows.Scan(
			&t.ID,
			&t.TenantName,
			&t.TenantType,
			&uen,
			&email,
			&phone,
			&username,
			&t.Status,
			&t.IsBCA,
			&t.Latitude,
			&t.Longitude,
			&t.Address,
			&t.WorkerCount,
			&t.DeviceCount,
		); err != nil {
			log.Printf("Scan error in GetTenants: %v", err)
			continue
		}
		if uen.Valid {
			t.UEN = uen.String
		}
		if email.Valid {
			t.Email = email.String
		}
		if phone.Valid {
			t.Phone = phone.String
		}
		if username.Valid {
			t.Username = username.String
		}
		tenants = append(tenants, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

func (h *TenantsHandler) GetTenantById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var t Tenant
	var uen, email, phone, username sql.NullString
	err := h.DB.QueryRow(`
        SELECT 
            t.tenant_id, 
            t.tenant_name, 
            t.tenant_type, 
            c.uen,
            t.contact_email,
            t.contact_phone,
            t.username,
            t.status,
            IFNULL(c.company_type = 'contractor', 0) as is_bca,
            t.latitude,
            t.longitude,
            t.address
        FROM tenants t
        LEFT JOIN companies c ON t.tenant_id = c.tenant_id AND c.company_type = 'contractor'
        WHERE t.tenant_id = ?`, id).
		Scan(&t.ID, &t.TenantName, &t.TenantType, &uen, &email, &phone, &username, &t.Status, &t.IsBCA, &t.Latitude, &t.Longitude, &t.Address)

	if err == sql.ErrNoRows {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if uen.Valid {
		t.UEN = uen.String
	}
	if email.Valid {
		t.Email = email.String
	}
	if phone.Valid {
		t.Phone = phone.String
	}
	if username.Valid {
		t.Username = username.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TenantsHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var t Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate ID first
	t.ID = "tenant-" + uuid.New().String()

	if t.TenantName == "" || t.TenantType == "" || t.Email == "" {
		http.Error(w, "Missing required tenant fields", http.StatusBadRequest)
		return
	}

	// Auto-generate username if missing
	if t.Username == "" {
		// simple sanitization: replace spaces with _, lowercase
		baseName := ""
		for _, r := range t.TenantName {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				baseName += string(r)
			} else if r == ' ' {
				baseName += "_"
			}
		}
		if len(baseName) > 15 {
			baseName = baseName[:15]
		}
		t.Username = baseName + "_" + t.ID[len(t.ID)-4:]
	}

	// Handle password
	var finalHash string
	if t.Password != "" {
		// Use provided password
		hash, err := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		finalHash = string(hash)
	} else {
		// Default: password123
		finalHash = "$2a$10$YGl2KZOrJ8oAtuyu5l59JuLCAeHZMfm15blSCSLwGAkfIU04c.F6G"
	}

	if t.Status == "" {
		t.Status = "active"
	}

	_, err := h.DB.Exec(`
		INSERT INTO tenants (
			tenant_id, tenant_name, tenant_type, contact_email, contact_phone, username, password_hash, status, address, latitude, longitude
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.TenantName, t.TenantType, t.Email, t.Phone, t.Username, finalHash, t.Status, t.Address, t.Latitude, t.Longitude)

	if err != nil {
		log.Printf("CreateTenant DB Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Do not return password in response
	t.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func (h *TenantsHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var t Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if password update is requested
	var errUpdate error
	if t.Password != "" {
		hash, errGen := bcrypt.GenerateFromPassword([]byte(t.Password), bcrypt.DefaultCost)
		if errGen != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		_, errUpdate = h.DB.Exec(`
			UPDATE tenants
			SET tenant_name=?, tenant_type=?, contact_email=?, contact_phone=?, username=?, status=?, latitude=?, longitude=?, address=?, password_hash=?
			WHERE tenant_id=?`,
			t.TenantName, t.TenantType, t.Email, t.Phone, t.Username, t.Status, t.Latitude, t.Longitude, t.Address, string(hash), id)
	} else {
		_, errUpdate = h.DB.Exec(`
			UPDATE tenants
			SET tenant_name=?, tenant_type=?, contact_email=?, contact_phone=?, username=?, status=?, latitude=?, longitude=?, address=?
			WHERE tenant_id=?`,
			t.TenantName, t.TenantType, t.Email, t.Phone, t.Username, t.Status, t.Latitude, t.Longitude, t.Address, id)
	}

	if errUpdate != nil {
		log.Printf("UpdateTenant DB Error: %v", errUpdate)
		http.Error(w, errUpdate.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *TenantsHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := h.DB.Exec("UPDATE tenants SET status='inactive' WHERE tenant_id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
