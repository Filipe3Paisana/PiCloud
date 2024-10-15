package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "io"
)

// Estrutura para enviar o status do nó
type NodeStatusRequest struct {
    NodeAddress       string `json:"node_address"`       
    Location          string `json:"location"`           
    Capacity          int    `json:"capacity"`           
    AvailableCapacity int    `json:"available_capacity"` 
    Status            string `json:"status"`             
}


func SendNodeStatusPeriodically() {
    ticker := time.NewTicker(20 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        var totalStorage uint64 = 100 * 1024 * 1024
        var availableStorage uint64 = 60 * 1024 * 1024

        nodeAddress := "127.0.0.1"
        location := "Datacenter XYZ"

        status := NodeStatusRequest{
            NodeAddress:       nodeAddress,
            Location:          location,
            Capacity:          int(totalStorage),
            AvailableCapacity: int(availableStorage),
            Status:            "online",
        }

        statusJSON, err := json.Marshal(status)
        if err != nil {
            fmt.Println("Erro ao serializar o status:", err)
            continue
        }

        fmt.Println("Enviando o seguinte status:", string(statusJSON)) // Log do JSON enviado

        endpoint := "http://localhost:8081/node/status/update"
        resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(statusJSON))
        if err != nil {
            fmt.Println("Erro ao enviar o status:", err)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            fmt.Println("Status enviado com sucesso")
        } else {
            body, _ := io.ReadAll(resp.Body)
            fmt.Printf("Código de resposta: %d, corpo: %s\n", resp.StatusCode, body) // Log do corpo da resposta
        }
    }
}

