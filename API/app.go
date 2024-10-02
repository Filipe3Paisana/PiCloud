package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=postgres-container port=5432 user=test password=test dbname=test sslmode=disable"
	
	// Tentativa de conexão com o banco de dados
	for {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Error connecting to the database, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Tenta fazer o ping no banco de dados
		if err = db.Ping(); err == nil {
			break
		}
		log.Println("Database not ready, retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	http.HandleFunc("/users", createUserHandler)

	fmt.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, user.Username, user.Password, user.Email).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
