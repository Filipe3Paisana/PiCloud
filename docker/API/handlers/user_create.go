package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"api/models"
	"api/utils"
	"api/db"
)

// Validações para proteger contra XSS e entradas inválidas
func validateUserInput(user *models.User) error {
	// Validação do nome de usuário (apenas letras, números, sublinhados, e tamanho máximo de 50 caracteres)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{1,50}$`)
	if !usernameRegex.MatchString(user.Username) {
		return errors.New("Username inválido: deve conter apenas letras, números ou sublinhados e ter no máximo 50 caracteres")
	}

	// Validação de e-mail (regex básico para verificar formato)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return errors.New("E-mail inválido: formato de e-mail incorreto")
	}

	// Validação de senha (mínimo de 8 caracteres, incluindo letras e números)
	if len(user.Password) < 8 || !strings.ContainsAny(user.Password, "0123456789") {
		return errors.New("Senha inválida: deve conter pelo menos 8 caracteres e incluir pelo menos um número")
	}

	return nil
}

// Adiciona headers de segurança à resposta
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none';")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

// Handler para criar um novo utilizador com validação e proteção contra XSS
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(w, r) // Configuração de CORS

	// Apenas aceita requisições POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Configuração de headers de segurança
	setSecurityHeaders(w)

	// Decodifica a requisição
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		http.Error(w, "Requisição inválida", http.StatusBadRequest)
		return
	}

	// Validação dos dados do usuário
	if err := validateUserInput(&user); err != nil {
		log.Printf("Erro de validação: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Erro ao gerar hash da senha: %v", err)
		http.Error(w, "Erro interno ao processar os dados", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id"
	err = db.DB.QueryRow(query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
	if err != nil {
		log.Printf("Erro ao salvar no banco de dados: %v", err)
		http.Error(w, "Erro ao criar utilizador", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Utilizador criado com sucesso",
		"user_id": user.ID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta JSON: %v", err)
		http.Error(w, "Erro interno ao processar a resposta", http.StatusInternalServerError)
	}
}
