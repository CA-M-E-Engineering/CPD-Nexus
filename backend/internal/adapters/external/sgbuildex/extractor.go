package sgbuildex

// This file used to contain direct SQL extraction logic.
// All extraction logic has been moved to the Repository layer (internal/adapters/repository/mysql/attendance_repo.go)
// to follow the modular architecture and avoid direct DB access in external adapters.
//
// The domain models for these extractions can be found in internal/core/domain/sgbuildex.go
