package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "node/models"
    "node/helpers"
)


// Função para enviar o status do nó periodicamente
func SendNodeStatusHandler() {
    ticker := time.NewTicker(20 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        nodeAddress, err := helpers.GetLocalIPAddress()
        if err != nil {
            fmt.Println("Erro ao obter o endereço IP:", err)
            continue
        }

        totalStorage, availableStorage, err := helpers.GetDiskUsage("/app/fragments/")
        if err != nil {
			totalStorage, availableStorage, err = helpers.GetDiskUsage("/")
			if err != nil {
				fmt.Println("Erro ao obter us de disco:", err)
				continue
			}
        }

        location := "Datacenter XYZ" // Defina isso dinamicamente, se necessário
        status := "online"           // Você pode mudar essa lógica com base em outras verificações

        nodeStatus := models.NodeStatusRequest{
            NodeAddress:       nodeAddress,
            Location:          location,
            Capacity:          int(totalStorage),
            AvailableCapacity: int(availableStorage),
            Status:            status,
        }

        statusJSON, err := json.Marshal(nodeStatus)
        if err != nil {
            fmt.Println("Erro ao serializar o status:", err)
            continue
        }

        fmt.Println("Enviando o seguinte status:", string(statusJSON))

        endpoint := "http://api-container:8080/node/status/update" 
        resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(statusJSON))
        if err != nil {
            fmt.Println("Erro ao enviar o status:", err)
            continue
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            fmt.Println("Status enviado com sucesso, Node: ", nodeAddress)
        } else {
            fmt.Printf("Código de resposta: %d\n", resp.StatusCode)
        }
    }
}