package main

import (
    "fmt"
    "net/http"
    
    "node/handlers"
)

func main() {
    // Endpoint for Node Status
    http.HandleFunc("/status", handlers.GetNodeStatusHandler)

    // Endpoint to upload a file fragment
    http.HandleFunc("/fragments/upload", handlers.UploadFragmentHandler)

    // Endpoint to download a file fragment by ID
    http.HandleFunc("/fragments/download", handlers.DownloadFragmentHandler)

    // Start the node on port 8082
    fmt.Println("Node rodando na porta 8082")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        fmt.Println("Erro ao iniciar o Node:", err)
    }
}
