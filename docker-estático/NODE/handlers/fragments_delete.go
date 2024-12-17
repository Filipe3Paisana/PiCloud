package handlers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
)

// DeleteFragmentHandler handles the deletion of file fragments
func DeleteFragmentHandler(w http.ResponseWriter, r *http.Request) {
    fileID := r.URL.Query().Get("file_id")
    fragmentOrder := r.URL.Query().Get("fragment_order")

    if fileID == "" || fragmentOrder == "" {
        http.Error(w, "File ID and fragment order are required", http.StatusBadRequest)
        return
    }

    // Formar o nome do fragmento com base no fileID e na ordem do fragmento
    fragmentFileName := fmt.Sprintf("file_%sfragment_%s", fileID, fragmentOrder)
    filePath := filepath.Join(fragmentStorageDir, fragmentFileName)

    // Verificar se o fragmento existe
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.Error(w, "Fragment not found", http.StatusNotFound)
        return
    }

    // Apagar o fragmento
    err := os.Remove(filePath)
    if err != nil {
        http.Error(w, "Failed to delete fragment", http.StatusInternalServerError)
        return
    }

    // Responder com sucesso
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Fragment deleted successfully")
}