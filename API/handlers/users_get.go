package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    
    "api/models"
    "api/utils"
)


func GetUsersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w, r)
        if r.Method != http.MethodGet {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        rows, err := db.Query("SELECT id, username, email FROM users")
        if err != nil {
            http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []models.User
        for rows.Next() {
            var user models.User
            if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
                http.Error(w, "Erro ao escanear usuário", http.StatusInternalServerError)
                return
            }
            users = append(users, user)
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}