package handlers

import (
    "encoding/json"
    "net/http"
    "log"
    
    "api/models"
    "api/utils"
    "api/db"
)


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