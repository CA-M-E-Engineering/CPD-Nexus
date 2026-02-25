package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadFaceHandler(w http.ResponseWriter, r *http.Request) {
	// 10 MB max memory limit
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to parse image from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure upload directory exists
	uploadDir := "./uploads/faces"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Create unique file name
	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	filepath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write file contents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Return the relative path
	json.NewEncoder(w).Encode(map[string]string{
		"path": "/uploads/faces/" + filename,
	})
}
