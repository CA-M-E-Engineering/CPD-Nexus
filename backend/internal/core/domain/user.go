package domain

type User struct {
	ID           string  `json:"user_id"`
	Name         string  `json:"user_name"`
	Username     string  `json:"username"`
	UserType     string  `json:"user_type"`
	Status       string  `json:"status"`
	Latitude     float64 `json:"lat"`
	Longitude    float64 `json:"lng"`
	ContactEmail string  `json:"email"`
	ContactPhone string  `json:"phone"`
	Address      string  `json:"address"`
	PasswordHash string  `json:"-"`
	WorkerCount  int     `json:"worker_count,omitempty"`
	DeviceCount  int     `json:"device_count,omitempty"`
}
