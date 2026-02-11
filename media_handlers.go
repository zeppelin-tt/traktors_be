package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const maxFileSize = 10 << 20 // 10 MB

// POST /media — загрузка картинки.
// multipart/form-data, поле «image».
func uploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)

	file, header, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read image: "+err.Error())
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !isImageContentType(contentType) {
		writeError(w, http.StatusBadRequest, "unsupported file type '"+contentType+"', allowed: jpeg, png, gif, webp, heic")
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := primitive.NewObjectID().Hex() + ext

	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to write file")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"url": baseURL + "/media/" + filename,
	})
}

func isImageContentType(ct string) bool {
	switch strings.ToLower(ct) {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/heic", "image/heif":
		return true
	}
	return false
}
