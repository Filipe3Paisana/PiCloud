package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strings"
    "log"
    
    "api/models"
    "api/utils"
    "api/db"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

    utils.EnableCors(w, r)
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Requisição inválida", http.StatusBadRequest)
        return
    }

    HashPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        http.Error(w, "Erro ao gerar hash da senha", http.StatusInternalServerError)
        return
    }

    query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
    err = db.DB.QueryRow(query, user.Username, user.Email, HashPassword).Scan(&user.ID) 
    if err != nil {
        http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Requisição inválida", http.StatusBadRequest)
        return
    }

    var user models.User
    query := "SELECT id, email, password_hash, username FROM users WHERE email=$1"
    err := db.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Email, &user.Password, &user.Username)
    if err != nil {
        http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
        return
    }

    if err := utils.ComparePassword(credentials.Password, user.Password); err != nil {
        http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
        return
    }

    
    token, err := utils.GenerateJWT(user.ID, user.Username , user.Email)
    if err != nil {
        log.Printf("Erro ao gerar o token: %v", err)
        http.Error(w, "Erro interno", http.StatusInternalServerError)
        return
    }

    // Retornar o token JWT no corpo da resposta
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token":   token,
        "message": "Login bem-sucedido",
    })
}

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

func GetUserHandler(db *sql.DB) http.HandlerFunc {
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
        err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Usuário não encontrado", http.StatusNotFound)
            } else {
                http.Error(w, "Erro ao buscar usuário", http.StatusInternalServerError)
            }
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}