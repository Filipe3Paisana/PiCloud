package models

type Node struct {
    NodeID            int    `json:"node_id"`           // Adicionado NodeID para identificar nós existentes
    NodeAddress       string `json:"node_address"`      // Endereço do nó (IP ou nome de domínio)
    Location          string `json:"location"`          // Localização do nó
    Capacity          int    `json:"capacity"`          // Capacidade total do nó
    AvailableCapacity int    `json:"available_capacity"` // Capacidade disponível
    Status            string `json:"status"`            // Status do nó (ex: "online", "offline")
}
