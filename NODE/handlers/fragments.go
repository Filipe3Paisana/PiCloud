package handlers

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
)

const fragmentStorageDir = "./fragments"

// UploadFragmentHandler handles the upload of file fragments
func UploadFragmentHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(10 << 20) // Limit the size to 10MB
    if err != nil {
        http.Error(w, "File too large", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("fragment")
    if err != nil {
        http.Error(w, "Error retrieving the file fragment", http.StatusBadRequest)
        return
    }
    defer file.Close()

    if _, err := os.Stat(fragmentStorageDir); os.IsNotExist(err) {
        os.Mkdir(fragmentStorageDir, os.ModePerm)
    }

    filePath := filepath.Join(fragmentStorageDir, handler.Filename)
    dst, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Error creating file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Fragment uploaded successfully")
}

// DownloadFragmentHandler handles serving file fragments
func DownloadFragmentHandler(w http.ResponseWriter, r *http.Request) {
    fragmentID := r.URL.Query().Get("id") // Getting the fragment ID from query params

    if fragmentID == "" {
        http.Error(w, "Fragment ID is required", http.StatusBadRequest)
        return
    }

    filePath := filepath.Join(fragmentStorageDir, fragmentID)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.Error(w, "Fragment not found", http.StatusNotFound)
        return
    }

    http.ServeFile(w, r, filePath)
}
