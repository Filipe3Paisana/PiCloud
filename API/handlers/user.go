package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "api/models"
    "api/utils"
)

func CreateUserHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w)
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        var user models.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Requisição inválida", http.StatusBadRequest)
            return
        }

        query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
        err := db.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
        if err != nil {
            http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w)
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        var credentials struct {
            Email string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
            http.Error(w, "Requisição inválida", http.StatusBadRequest)
            return
        }

        var user models.User
        query := "SELECT id, email, password_hash FROM users WHERE email=$1"
        err := db.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Email, &user.Password)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
            } else {
                http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
            }
            return
        }

        if user.Password != credentials.Password {
            http.Error(w, "Senha incorreta", http.StatusUnauthorized)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"message": "Login bem-sucedido"})
    }
}

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        utils.EnableCors(w)
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
