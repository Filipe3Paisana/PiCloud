package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strings"
    
    "api/models"
    "api/utils"
	"api/db"
)


func GetUserHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w, r)
        if r.Method != http.MethodGet {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        id := strings.TrimPrefix(r.URL.Path, "/user/")
        if id == "" {
            http.Error(w, "ID não fornecido", http.StatusBadRequest)
            return
        }

        var user models.User
        query := "SELECT id, username, email FROM users WHERE id=$1"
        err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Utilizador não encontrado", http.StatusNotFound)
            } else {
                http.Error(w, "Erro ao procurar utilizador", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}
