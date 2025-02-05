package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"database/sql"

	"api/db"
	"api/models"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO Permitir conexões de qualquer origem (melhorar para produção)
		return true
	},
}

var connections = make(map[*websocket.Conn]bool)
var connMutex = sync.Mutex{}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao criar WebSocket:", err)
		return
	}
	defer conn.Close()

	connMutex.Lock()
	connections[conn] = true
	connMutex.Unlock()
	
	fmt.Println("Nova conexão WebSocket")

	for {
		var msg models.Node
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			break
		}
		fmt.Println("Mensagem recebida:", msg)


		// Inserir/Atualizar status do Node na base de dados
		err = updateNodeStatusInDB(msg)
		if err != nil {
			fmt.Printf("Erro ao armazenar status do Node: %v\n", err)
			continue
		}

		// Confirmação ao Node
		err = conn.WriteJSON(map[string]string{
			"message": "Status recebido e armazenado com sucesso!",
		})
		if err != nil {
			fmt.Println("Erro ao enviar resposta:", err)
			break
		}
	}

	connMutex.Lock()
	delete(connections, conn)
	connMutex.Unlock()
}


// Função para enviar fragmentos aos Nodes conectados via WebSocket
func SendUploadCommandToNodes(nodeAddress string, fragmentData []byte, fragmentOrder int) error {
	connMutex.Lock()
	defer connMutex.Unlock()

	for conn := range connections {
		// Enviar apenas para o Node específico
		if err := conn.WriteJSON(map[string]interface{}{
			"type":           "upload_fragment",
			"fragment_order": fragmentOrder,
			"data":           fragmentData,
			"node_address":   nodeAddress,
		}); err != nil {
			fmt.Printf("Erro ao enviar fragmento ao Node %s: %v\n", nodeAddress, err)
			return err
		}
		fmt.Printf("Fragmento enviado ao Node %s, ordem %d\n", nodeAddress, fragmentOrder)
	}
	return nil
}

// Função para inserir/atualizar status do Node na base de dados
func updateNodeStatusInDB(status models.Node) error {
	// Verifica se o nó já existe
	var nodeID int
	err := db.DB.QueryRow("SELECT id FROM Nodes WHERE node_address = $1", status.NodeAddress).Scan(&nodeID)

	if err == sql.ErrNoRows {
		// Inserir novo registro se não existir
		_, err = db.DB.Exec(
			"INSERT INTO Nodes (node_address, location, capacity, available_capacity, status, last_updated) VALUES ($1, $2, $3, $4, $5, NOW())",
			status.NodeAddress, status.Location, status.Capacity, status.AvailableCapacity, status.Status,
		)
		if err != nil {
			return fmt.Errorf("erro ao inserir o nó: %w", err)
		}
		fmt.Println("Nó adicionado com sucesso:", status.NodeAddress)
	} else if err == nil {
		// Atualizar registro existente
		_, err = db.DB.Exec(
			"UPDATE Nodes SET location = $1, capacity = $2, available_capacity = $3, status = $4, last_updated = NOW() WHERE id = $5",
			status.Location, status.Capacity, status.AvailableCapacity, status.Status, nodeID,
		)
		if err != nil {
			return fmt.Errorf("erro ao atualizar o nó: %w", err)
		}
		fmt.Println("Nó atualizado com sucesso:", status.NodeAddress)
	} else {
		return fmt.Errorf("erro ao verificar o nó: %w", err)
	}

	return nil
}