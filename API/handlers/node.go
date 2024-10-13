package handlers

import (
    "net/http"
    "io/ioutil"
    "api/utils"
)

func CheckNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
    utils.EnableCors(w, r)

    if r.Method != http.MethodGet {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    
    nodeURL := "http://node-container:8082/status"

    
    resp, err := http.Get(nodeURL)
    if err != nil {
        http.Error(w, "Erro ao se conectar com o Node", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Erro ao ler a resposta do Node", http.StatusInternalServerError)
        return
    }

    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}
