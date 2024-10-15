package handlers

import (
    "net/http"
    "io"
    "api/utils"
    "database/sql"
    "encoding/json"
    "fmt"

)


type NodeStatusRequest struct {
    NodeAddress      string `json:"node_address"`      // Endereço do nó (IP ou nome de domínio)
    Location         string `json:"location"`          // Localização do nó
    Capacity         int    `json:"capacity"`          // Capacidade total do nó
    AvailableCapacity int    `json:"available_capacity"` // Capacidade disponível
    Status           string `json:"status"`            // Status do nó (ex: "online", "offline")
}


func CheckNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)

    if r.Method != http.MethodGet {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    
    nodeURL := "http://node:8082/status"

    
    resp, err := http.Get(nodeURL)
    if err != nil {
        http.Error(w, "Erro ao se conectar com o Node", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Erro ao ler a resposta do Node", http.StatusInternalServerError)
        return
    }

    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}

func UpdateNodeStatusHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
            return
        }

        var req NodeStatusRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Erro ao decodificar a requisição", http.StatusBadRequest)
            return
        }

        var nodeID int
        err := db.QueryRow("SELECT id FROM Nodes WHERE node_address = $1", req.NodeAddress).Scan(&nodeID)

        if err == sql.ErrNoRows {
            // Se o nó não existir, insere um novo registro
            _, err = db.Exec(
                "INSERT INTO Nodes (node_address, location, capacity, available_capacity, status) VALUES ($1, $2, $3, $4, $5)",
                req.NodeAddress, req.Location, req.Capacity, req.AvailableCapacity, req.Status,
            )
            if err != nil {
                http.Error(w, "Erro ao inserir o nó", http.StatusInternalServerError)
                return
            }
            fmt.Fprintln(w, "Nó adicionado com sucesso")
        } else if err == nil {
            // Se o nó já existir, atualiza o registro
            _, err = db.Exec(
                "UPDATE Nodes SET location = $1, capacity = $2, available_capacity = $3, status = $4 WHERE id = $5",
                req.Location, req.Capacity, req.AvailableCapacity, req.Status, nodeID,
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
