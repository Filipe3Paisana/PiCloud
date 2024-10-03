package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "sync"
    "time"

    _ "github.com/lib/pq"
)

var mu sync.Mutex

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password_hash"`
    Email    string `json:"email"`
}

var db *sql.DB

// Função para habilitar CORS
func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// Handler para criação de usuários
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(w) // Habilitar CORS

    if r.Method == http.MethodOptions {
        return // Retorna para as requisições OPTIONS necessárias para CORS
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    var user User
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

// Handler para obter usuários
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(w) // Habilitar CORS

    if r.Method == http.MethodOptions {
        return // Retorna para as requisições OPTIONS necessárias para CORS
    }

    if r.Method != http.MethodGet {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    rows, err := db.Query("SELECT id, username, email, password_hash FROM users")
    if err != nil {
        http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
            http.Error(w, "Erro ao escanear usuário", http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Erro ao processar usuários", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func main() {

    var err error
    connStr := "host=postgres-container port=5432 user=test password=test dbname=test sslmode=disable"

    for {
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            fmt.Println("Erro ao conectar ao banco de dados, tentando novamente em 5 segundos...")
            time.Sleep(5 * time.Second)
            continue
        }

        if err = db.Ping(); err == nil {
            break
        }
        fmt.Println("Banco de dados não está pronto, tentando novamente em 5 segundos...")
        time.Sleep(5 * time.Second)
    }

    rand.Seed(time.Now().UnixNano())

    // Definindo rotas e handlers
    http.HandleFunc("/users/add", createUserHandler)
    http.HandleFunc("/users", getUsersHandler)

    fmt.Println("Servidor rodando em http://localhost:8081/")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
