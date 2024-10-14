package handlers

import (
    "encoding/json"
    "net/http"
)

// NodeStatus contains information about the node's current state
type NodeStatus struct {
    Status           string `json:"status"` // "online", "offline"
    TotalStorage     uint64 `json:"total_storage"`
    AvailableStorage uint64 `json:"available_storage"`
    Active           bool   `json:"active"` // Is the node actively serving files?
}

// GetNodeStatusHandler handles the node status reporting
func GetNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
    var totalStorage uint64 = 100 * 1024 * 1024 // 100 MB for testing purposes
    var availableStorage uint64 = 60 * 1024 * 1024 // 60 MB available

    status := NodeStatus{
        Status:           "online",
        TotalStorage:     totalStorage,
        AvailableStorage: availableStorage,
        Active:           true,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}
