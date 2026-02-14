package domain

type User struct {
	ID        string `json:"user_id"`
	Name      string `json:"user_name"`
	Username  string `json:"username"`
	UserType  string `json:"user_type"`
	Status    string `json:"status"`
	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
}
