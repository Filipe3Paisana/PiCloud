package helpers

import (
    "fmt"
    "time"

    "api/models"
    "api/db"
)

func MarkOfflineNodes() {
    for {
        time.Sleep(10 * time.Second) 

        offlineThreshold := time.Now().Add(-5 * time.Second)

        // Atualiza os nós para "offline" onde last_updated está além do limite
        _, err := db.DB.Exec("UPDATE Nodes SET status = 'offline' WHERE last_updated < $1", offlineThreshold)
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

func SelectNodesForFragment(availableNodes []models.Node, replicationFactor int) []models.Node {
	connMutex.Lock()
	defer connMutex.Unlock()

	// Filtrar Nodes que possuem conexão WebSocket ativa
	var activeNodes []models.Node
	for _, node := range availableNodes {
		for conn := range connections {
			if node.NodeAddress == conn.RemoteAddr().String() {
				activeNodes = append(activeNodes, node)
				break
			}
		}
	}

	if replicationFactor >= len(activeNodes) {
		return activeNodes
	}

	// Selecionar Nodes aleatórios entre os ativos
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(activeNodes), func(i, j int) {
		activeNodes[i], activeNodes[j] = activeNodes[j], activeNodes[i]
	})
	return activeNodes[:replicationFactor]
}