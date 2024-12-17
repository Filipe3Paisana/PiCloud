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