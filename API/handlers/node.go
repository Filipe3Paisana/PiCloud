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


func UpdateNodeStatusHandler() http.HandlerFunc {
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
        err := db.DB.QueryRow("SELECT id FROM Nodes WHERE node_address = $1", req.NodeAddress).Scan(&nodeID)

        if err == sql.ErrNoRows {
            // Se o nó não existir, insere um novo registro
            _, err = db.DB.Exec(
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
            _, err = db.DB.Exec(
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

