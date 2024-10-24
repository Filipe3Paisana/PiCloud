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