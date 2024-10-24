package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    
    "api/models"
    "api/utils"
)


func GetUserFilesHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w, r)
        
        if r.Method != http.MethodGet {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        userID, err := utils.ExtractUserIDFromJWT(r)
        if err != nil || userID == 0 {
            http.Error(w, "ID não fornecido", http.StatusBadRequest)
            return
        }

        rows, err := db.Query("SELECT id, name, size FROM files WHERE user_id=$1", userID)
        if err != nil {
            http.Error(w, "Erro ao buscar arquivos", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var files []models.File
        for rows.Next() {
            var file models.File
            if err := rows.Scan(&file.ID, &file.Name, &file.Size); err != nil {
                http.Error(w, "Erro ao escanear arquivo", http.StatusInternalServerError)
                return
            }
            files = append(files, file)
        }

        // Garantir que uma lista vazia seja retornada em vez de null
        if files == nil {
            files = []models.File{}
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(files); err != nil {
            http.Error(w, "Erro ao codificar resposta JSON", http.StatusInternalServerError)
        }
    }
}
