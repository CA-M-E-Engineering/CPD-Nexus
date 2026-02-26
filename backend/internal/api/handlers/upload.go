package handlers

import (
	"encoding/json"
	"fmt"
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

	// 1. Get Trade from request (fallback to 'general')
	trade := r.FormValue("trade")
	if trade == "" {
		trade = "general"
	}
	// Sanitize trade for folder name
	trade = filepath.Base(trade)

	// 2. Ensure upload directory exists: ./uploads/faces/<trade>
	uploadSubDir := filepath.Join("uploads", "faces", trade)
	if err := os.MkdirAll(uploadSubDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// 3. Create unique file name
	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	savePath := filepath.Join(uploadSubDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write file contents", http.StatusInternalServerError)
		return
	}

	// 4. Construct the URL address
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	// Support for common proxy headers
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)
	fileURL := fmt.Sprintf("%s/uploads/faces/%s/%s", baseURL, trade, filename)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Return the image URL address
	json.NewEncoder(w).Encode(map[string]string{
		"url":  fileURL,
		"path": fmt.Sprintf("/uploads/faces/%s/%s", trade, filename),
	})
}
