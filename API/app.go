package main

import (
    "fmt"
    "net/http"

    "api/db"
    "api/handlers"
    "api/utils"

    _ "github.com/lib/pq"
)

func main() {
    // Conectar ao banco de dados
    dbConn, err := db.Connect()
    if err != nil {
        fmt.Println("Erro ao conectar ao banco de dados:", err)
        return
    }
    defer dbConn.Close()

    // Definir rotas e handlers
    http.HandleFunc("/users/add", handlers.CreateUserHandler(dbConn))      // Rota pública
    http.HandleFunc("/users/login", handlers.LoginHandler(dbConn))          // Rota pública
    http.HandleFunc("/node/status", handlers.CheckNodeStatusHandler)        // Nova rota
    
    // Aplicar middleware de autenticação JWT nas rotas protegidas
    http.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUsersHandler(dbConn))))   
    http.Handle("/user/", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUserHandler(dbConn))))    

    fmt.Println("Servidor rodando em http://localhost:8081/")
    if err := http.ListenAndServe(":8080", nil); err != nil {  
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
