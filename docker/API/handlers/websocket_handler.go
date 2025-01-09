package handlers

import (
	"fmt"
	"net/http"
	"sync"

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
	// Atualiza a conexão HTTP para WebSocket
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
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			break
		}
		fmt.Println("Mensagem recebida:", msg)

		// Exemplo: responder ao node
		err = conn.WriteJSON(map[string]string{
			"message": "Mensagem recebida com sucesso!",
		})
		if err != nil {
			fmt.Println("Erro ao enviar resposta:", err)
			break
		}
	}

	connMutex.Lock()
	delete(connections, conn)
	connMutex.Unlo