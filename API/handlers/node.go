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

    // Endereço do Node
    nodeURL := "http://node-container:8082/status"

    // Realizar a requisição ao Node
    resp, err := http.Get(nodeURL)
    if err != nil {
        http.Error(w, "Erro ao se conectar com o Node", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Ler a resposta do Node
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "Erro ao ler a resposta do Node", http.StatusInternalServerError)
        return
    }

    // Enviar a resposta para o cliente
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}
