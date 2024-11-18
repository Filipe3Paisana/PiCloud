package handlers

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
)
//FIXME verificar que podemos mandar um node abaixo e fazer download na mesma. 

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
