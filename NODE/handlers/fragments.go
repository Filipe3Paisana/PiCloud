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

func DownloadFragmentHandler(w http.ResponseWriter, r *http.Request) { //TODO falta guardar o tipo de ficheiro para colocar a extensão no fim (PDF, TXT, ...)
    // Extrair `file_id` e `fragment_order` dos parâmetros da URL
    fileID := r.URL.Query().Get("file_id")
    fragmentOrder := r.URL.Query().Get("fragment_order")

    // Verificar se os parâmetros necessários foram fornecidos
    if fileID == "" || fragmentOrder == "" {
        http.Error(w, "File ID and fragment order are required", http.StatusBadRequest)
        return
    }

    // Formar o nome do fragmento com base no fileID e na ordem do fragmento
    fragmentFileName := fmt.Sprintf("file_%sfragment_%s", fileID, fragmentOrder)
    filePath := filepath.Join(fragmentStorageDir, fragmentFileName)

    // Adicionar logs para ajudar na depuração
    fmt.Printf("Tentando servir o fragmento: %s\n", filePath)

    // Verificar se o arquivo existe
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        fmt.Printf("Fragmento não encontrado: %s\n", filePath)
        http.Error(w, "Fragment not found", http.StatusNotFound)
        return
    }

    // Servir o arquivo como resposta
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fragmentFileName))
    w.Header().Set("Content-Type", "application/octet-stream")
    http.ServeFile(w, r, filePath)
}
