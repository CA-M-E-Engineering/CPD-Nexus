package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sgbuildex/internal/pkg/idgen"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	DB *sql.DB
}

func NewUsersHandler(db *sql.DB) *UsersHandler {
	return &UsersHandler{DB: db}
}

type User struct {
	ID          string  `json:"user_id"`
	UserName    string  `json:"user_name"`
	UserType    string  `json:"user_type"`
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

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
        SELECT 
            u.user_id, 
            COALESCE(u.user_name, 'Unknown'), 
            COALESCE(u.user_type, 'client'), 
            NULL as uen,
            u.contact_email,
            u.contact_phone,
            u.username,
            COALESCE(u.status, 'pending'),
            0 as is_bca,
            COALESCE(u.latitude, 0),
            COALESCE(u.longitude, 0),
            COALESCE(u.address, ''),
            (SELECT COUNT(*) FROM workers w WHERE w.user_id = u.user_id AND w.role = 'worker' AND w.status = 'active') as worker_count,
            (SELECT COUNT(*) FROM devices d WHERE d.user_id = u.user_id AND d.status != 'inactive') as device_count
        FROM users u
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		var uen, email, phone, username sql.NullString
		if err := rows.Scan(
			&u.ID,
			&u.UserName,
			&u.UserType,
			&uen,
			&email,
			&phone,
			&username,
			&u.Status,
			&u.IsBCA,
			&u.Latitude,
			&u.Longitude,
			&u.Address,
			&u.WorkerCount,
			&u.DeviceCount,
		); err != nil {
			log.Printf("Scan error in GetUsers: %v", err)
			continue
		}
		if uen.Valid {
			u.UEN = uen.String
		}
		if email.Valid {
			u.Email = email.String
		}
		if phone.Valid {
			u.Phone = phone.String
		}
		if username.Valid {
			u.Username = username.String
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UsersHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	var uen, email, phone, username sql.NullString
	err := h.DB.QueryRow(`
        SELECT 
            u.user_id, 
            u.user_name, 
            u.user_type, 
            NULL as uen,
            u.contact_email,
            u.contact_phone,
            u.username,
            u.status,
            0 as is_bca,
            u.latitude,
            u.longitude,
            u.address
        FROM users u
        WHERE u.user_id = ?`, id).
		Scan(&u.ID, &u.UserName, &u.UserType, &uen, &email, &phone, &username, &u.Status, &u.IsBCA, &u.Latitude, &u.Longitude, &u.Address)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if uen.Valid {
		u.UEN = uen.String
	}
	if email.Valid {
		u.Email = email.String
	}
	if phone.Valid {
		u.Phone = phone.String
	}
	if username.Valid {
		u.Username = username.String
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate sequential ID
	id, errID := idgen.GenerateNextID(h.DB, "users", "user_id", "user")
	if errID != nil {
		log.Printf("Failed to generate user ID: %v", errID)
		http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
		return
	}
	u.ID = id

	if u.UserName == "" || u.UserType == "" || u.Email == "" {
		http.Error(w, "Missing required user fields", http.StatusBadRequest)
		return
	}

	// Auto-generate username if missing
	if u.Username == "" {
		baseName := ""
		for _, r := range u.UserName {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				baseName += string(r)
			} else if r == ' ' {
				baseName += "_"
			}
		}
		if len(baseName) > 15 {
			baseName = baseName[:15]
		}
		u.Username = baseName + "_" + u.ID[len(u.ID)-4:]
	}

	// Handle password
	var finalHash string
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		finalHash = string(hash)
	} else {
		// Default: password123
		finalHash = "$2a$10$YGl2KZOrJ8oAtuyu5l59JuLCAeHZMfm15blSCSLwGAkfIU04c.F6G"
	}

	if u.Status == "" {
		u.Status = "active"
	}

	_, err := h.DB.Exec(`
		INSERT INTO users (
			user_id, user_name, user_type, contact_email, contact_phone, username, password_hash, status, address, latitude, longitude
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.ID, u.UserName, u.UserType, u.Email, u.Phone, u.Username, finalHash, u.Status, u.Address, u.Latitude, u.Longitude)

	if err != nil {
		log.Printf("CreateUser DB Error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Do not return password in response
	u.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if password update is requested
	var errUpdate error
	if u.Password != "" {
		hash, errGen := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if errGen != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		_, errUpdate = h.DB.Exec(`
			UPDATE users
			SET user_name=?, user_type=?, contact_email=?, contact_phone=?, username=?, status=?, latitude=?, longitude=?, address=?, password_hash=?
			WHERE user_id=?`,
			u.UserName, u.UserType, u.Email, u.Phone, u.Username, u.Status, u.Latitude, u.Longitude, u.Address, string(hash), id)
	} else {
		_, errUpdate = h.DB.Exec(`
			UPDATE users
			SET user_name=?, user_type=?, contact_email=?, contact_phone=?, username=?, status=?, latitude=?, longitude=?, address=?
			WHERE user_id=?`,
			u.UserName, u.UserType, u.Email, u.Phone, u.Username, u.Status, u.Latitude, u.Longitude, u.Address, id)
	}

	if errUpdate != nil {
		log.Printf("UpdateUser DB Error: %v", errUpdate)
		http.Error(w, errUpdate.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := h.DB.Exec("UPDATE users SET status='inactive' WHERE user_id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
