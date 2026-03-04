package domain

import "time"

// PitstopAuthorisation represents a cached Pitstop API routing configuration
type PitstopAuthorisation struct {
	PitstopAuthID  string     `json:"pitstop_auth_id" db:"pitstop_auth_id"`
	DatasetID      string     `json:"dataset_id" db:"dataset_id"`
	DatasetName    string     `json:"dataset_name" db:"dataset_name"`
	UserID         *string    `json:"user_id" db:"user_id"`
	RegulatorID    string     `json:"regulator_id" db:"regulator_id"`
	RegulatorName  string     `json:"regulator_name" db:"regulator_name"`
	OnBehalfOfID   string     `json:"on_behalf_of_id" db:"on_behalf_of_id"`
	OnBehalfOfName string     `json:"on_behalf_of_name" db:"on_behalf_of_name"`
	Status         string     `json:"status" db:"status"`
	LastSyncedAt   *time.Time `json:"last_synced_at" db:"last_synced_at"`
}
