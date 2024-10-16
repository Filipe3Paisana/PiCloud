package main

import (
    "fmt"
    "net/http"
    "node/handlers"
)


func main() {
    
    go handlers.SendNodeStatusPeriodically()

    //http.HandleFunc("/status", handlers.GetNodeStatusHandler)

    http.HandleFunc("/fragments/upload", handlers.UploadFragmentHandler)

    // Endpoint para download de fragmento de arquivo por ID
    http.HandleFunc("/fragments/download", handlers.DownloadFragmentHandler)


    fmt.Println("Node a bombar na porta 8082")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        fmt.Println("Erro ao iniciar o Node:", err)
    }
}
