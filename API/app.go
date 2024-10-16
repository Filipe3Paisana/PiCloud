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
    
    dbConn, err := db.Connect()
    if err != nil {
        fmt.Println("Erro ao conectar ao banco de dados:", err)
        return
    }
    defer dbConn.Close()

    http.HandleFunc("/users/add", handlers.CreateUserHandler(dbConn))      
    http.HandleFunc("/users/login", handlers.LoginHandler(dbConn))         
    http.HandleFunc("/user/upload", handlers.UploadHandler)       
    http.HandleFunc("/node/status/update", handlers.UpdateNodeStatusHandler(dbConn))

    http.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUsersHandler(dbConn))))   
    http.Handle("/user/", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUserHandler(dbConn))))    

    go handlers.MarkOfflineNodes(dbConn)

    fmt.Println("Servidor rodando em http://localhost:8081/")
    if err := http.ListenAndServe(":8080", nil); err != nil {  
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
