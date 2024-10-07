package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "Node is online")
    })

    fmt.Println("Node rodando na porta 8082")
    if err := http.ListenAndServe(":8082", nil); err != nil {
        fmt.Println("Erro ao iniciar o Node:", err)
    }
}
