package handlers

import (
    "net/http"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "api/models"
    "api/db"
)


func UpdateNodeStatusHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        var req models.Node
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
            return
        }

        
        currentTime := time.Now()

        var nodeID int
        err := db.QueryRow("SELECT id FROM Nodes WHERE node_address = $1", req.NodeAddress).Scan(&nodeID)

        if err == sql.ErrNoRows {
            // Se o nó não existir, insere um novo registro
            _, err = db.Exec(
                "INSERT INTO Nodes (node_address, location, capacity, available_capacity, status, last_updated) VALUES ($1, $2, $3, $4, $5, $6)",
                req.NodeAddress, req.Location, req.Capacity, req.AvailableCapacity, req.Status, currentTime,
            )
            if err != nil {
                http.Error(w, "Erro ao inserir o nó", http.StatusInternalServerError)
                return
            }
            fmt.Fprintln(w, "Nó adicionado com sucesso")
        } else if err == nil {
            // Se o nó já existir, atualiza o registro
            _, err = db.Exec(
                "UPDATE Nodes SET location = $1, capacity = $2, available_capacity = $3, status = $4, last_updated = $5 WHERE id = $6",
                req.Location, req.Capacity, req.AvailableCapacity, req.Status, currentTime, nodeID,
            )
            if err != nil {
                http.Error(w, "Erro ao atualizar o nó", http.StatusInternalServerError)
                return
            }
            fmt.Fprintln(w, "Nó atualizado com sucesso")
        } else {
            http.Error(w, "Erro ao verificar o nó", http.StatusInternalServerError)
        }
    }
}

func MarkOfflineNodes(db *sql.DB) {
    for {
        time.Sleep(30 * time.Second) 

        offlineThreshold := time.Now().Add(-25 * time.Second)

        // Atualiza os nós para "offline" onde last_updated está além do limite
        _, err := db.Exec("UPDATE Nodes SET status = 'offline' WHERE last_updated < $1", offlineThreshold)
        if err != nil {
            fmt.Println("Erro ao atualizar status dos nós para offline:", err)
        }
        numberOfNodes := GetNumberOfOnlineNodes()
        fmt.Println("Número de nós online:", numberOfNodes)
    }
}


func GetNumberOfOnlineNodes() int {
    var onlineNodes int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM Nodes WHERE status = 'online'").Scan(&onlineNodes)
    if err != nil {
        fmt.Println("Erro ao obter o número de nós online:", err)
    }
    return onlineNodes
}

// Lista todos os nós que estão online
func GetOnlineNodesList() []models.Node {
    var onlineNodesList []models.Node

    rows, err := db.DB.Query("SELECT id, node_address, location, capacity, available_capacity, status FROM Nodes WHERE status = 'online'")
    if err != nil {
        fmt.Println("Erro ao obter a lista de nós online:", err)
        return nil
    }
    defer rows.Close()

    for rows.Next() {
        var node models.Node
        if err := rows.Scan(&node.NodeID, &node.NodeAddress, &node.Location, &node.Capacity, &node.AvailableCapacity, &node.Status); err != nil {
            fmt.Println("Erro ao escanear dados do nó:", err)
            continue
        }
        onlineNodesList = append(onlineNodesList, node)
    }

    if err := rows.Err(); err != nil {
        fmt.Println("Erro ao iterar sobre os resultados:", err)
        return nil
    }

    return onlineNodesList
}