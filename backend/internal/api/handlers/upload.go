package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const maxUploadSizeBytes = 5 * 1024 * 1024 // 5MB hard limit

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func UploadFaceHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce hard size limit before reading body
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSizeBytes)
	if err := r.ParseMultipartForm(2 << 20); err != nil { // 2MB memory buffer
		http.Error(w, "File too large (max 5MB)", http.StatusRequestEntityTooLarge)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to parse image from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Validate file extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if !allowedExtensions[ext] {
		http.Error(w, "Only JPEG, PNG, and WebP images are allowed", http.StatusBadRequest)
		return
	}

	// 3. Read first 512 bytes and detect MIME type via magic bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	contentType := http.DetectContentType(buf[:n])
	if !allowedMIMETypes[contentType] {
		http.Error(w, "File content is not a valid image (JPEG/PNG/WebP)", http.StatusBadRequest)
		return
	}
	// Seek back to start so we can copy the whole file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	// 4. Get Trade from request (fallback to 'general')
	trade := r.FormValue("trade")
	if trade == "" {
		trade = "general"
	}
	// Sanitize trade for folder name
	trade = filepath.Base(trade)

	// 5. Ensure upload directory exists: ./uploads/faces/<trade>
	uploadSubDir := filepath.Join("uploads", "faces", trade)
	if err := os.MkdirAll(uploadSubDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// 6. Use UUID-based filename to prevent path traversal and enumeration
	safeFilename := uuid.New().String() + ext
	savePath := filepath.Join(uploadSubDir, safeFilename)

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

	// 7. Construct the URL address
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	fileURL := fmt.Sprintf("%s/uploads/faces/%s/%s", baseURL, trade, safeFilename)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"url":  fileURL,
		"path": fmt.Sprintf("/uploads/faces/%s/%s", trade, safeFilename),
	})
}
