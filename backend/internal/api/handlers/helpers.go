package handlers

import (
	"net/http"
	"sgbuildex/internal/pkg/apperrors"
)

// writeError writes a typed HTTP error response.
// AppErrors are mapped to their defined status code; all other errors return 500
// with a generic message so internal details are never exposed to clients.
func writeError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		http.Error(w, appErr.Message, appErr.Code)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
