package main

import (
    "fmt"
    "net/http"

    "api/db"
    "api/handlers"
    "api/helpers"
    "api/utils"

    _ "github.com/lib/pq"
)

func main() {
    
    dbConn, err := db.Connect()
    if err != nil {
        fmt.Println("Erro ao conectar à base de dados:", err)
        return
    }
    defer dbConn.Close()

    http.HandleFunc("/users/add", handlers.CreateUserHandler)      
    http.HandleFunc("/users/login", handlers.LoginHandler)         
    http.HandleFunc("/user/upload", handlers.UploadHandler)
    http.HandleFunc("/user/file/details", handlers.GetFileDetailsHandler)
    
    http.Handle("/user/", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUserHandler())))
    http.Handle("/users", utils.AuthMiddleware(http.HandlerFunc(handlers.GetUsersHandler())))   
    
    http.HandleFunc("/node/status/update", handlers.UpdateNodeStatusHandler())

    http.Handle("/user/files", (http.HandlerFunc(handlers.GetUserFilesHandler())))    
    http.Handle("/user/download", (http.HandlerFunc(handlers.DownloadHandler)))
    http.Handle("/user/delete", (http.HandlerFunc(handlers.DeleteFileHandler)))



    go helpers.MarkOfflineNodes()

    fmt.Println("Servidor a bombar em http://localhost:8081/")
    if err := http.ListenAndServe(":8080", nil); err != nil {  
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
