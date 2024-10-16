package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "net/http"
    "syscall"
    "time"
)

// Estrutura do status do nó
type NodeStatusRequest struct {
    NodeAddress       string `json:"node_address"`
    Location          string `json:"location"`
    Capacity          int    `json:"capacity"`
    AvailableCapacity int    `json:"available_capacity"`
    Status            string `json:"status"`
}

// Função para obter o endereço IP local
func getLocalIPAddress() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
            if ipNet.IP.To4() != nil {
                return ipNet.IP.String(), nil
            }
        }
    }
    return "", fmt.Errorf("não foi possível obter o IP")
}

// Função para obter o uso de disco
func getDiskUsage(path string) (total uint64, free uint64, err error) {
    var stat syscall.Statfs_t
    err = syscall.Statfs(path, &stat)
    if err != nil {
        return 0, 0, err
    }

    total = stat.Blocks * uint64(stat.Bsize)       // Capacidade total
    free = stat.Bavail * uint64(stat.Bsize)        // Capacidade disponível
    return total, free, nil
}

// Função para enviar o status do nó periodicamente
func SendNodeStatusPeriodically() {
    ticker := time.NewTicker(20 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        nodeAddress, err := getLocalIPAddress()
        if err != nil {
            fmt.Println("Erro ao obter o endereço IP:", err)
            continue
        }

        totalStorage, availableStorage, err := getDiskUsage("/app/fragments/")
        if err != nil {
			totalStorage, availableStorage, err = getDiskUsage("/")
			if err != nil {
				fmt.Println("Erro ao obter us de disco:", err)
				continue
			}
        }

        location := "Datacenter XYZ" // Defina isso dinamicamente, se necessário
        status := "online"           // Você pode mudar essa lógica com base em outras verificações

        nodeStatus := NodeStatusRequest{
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
            fmt.Println("Status enviado com sucesso")
        } else {
            fmt.Printf("Código de resposta: %d\n", resp.StatusCode)
        }
    }
}