package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "sync"
    "log"
)

//FIXME as metricas não estão a funcionar a 100%, as vezes funciona, as vezes tem erro ao decodificar o payload. 

type Metric struct {
    Instance string  `json:"instance"`
    Value    float64 `json:"value"`
}

type MetricsData struct {
    NodeID    string                 `json:"node_id"`
    Timestamp int64                  `json:"timestamp"`
    Metrics   map[string][]Metric    `json:"metrics"`
}

var mu sync.Mutex  // Mutex para proteger o acesso ao ficheiro

func receiveMetricsHandler(w http.ResponseWriter, r *http.Request) {
    var data MetricsData
    data.Timestamp = time.Now().Unix()

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
        return
    }
    r.Body.Close()

    // Log do payload recebido
    log.Printf("Payload recebido: %s", string(bodyBytes))

    // Decodificar o payload
    if err := json.Unmarshal(bodyBytes, &data); err != nil {
        http.Error(w, "Erro ao decodificar payload", http.StatusBadRequest)
        log.Printf("Erro ao decodificar payload: %v\n", err)
        return
    }

    // guardar métricas no ficheiro all_metrics.json
    if err := appendAllMetricsData(data); err != nil {
        http.Error(w, "Erro ao guardar métricas", http.StatusInternalServerError)
        log.Printf("Erro ao guardar métricas: %v\n", err)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Métricas recebidas e armazenadas com sucesso"))
}

func appendAllMetricsData(data MetricsData) error {
    mu.Lock()
    defer mu.Unlock()

    filename := "prometheus_data/all_metrics.json"
    var existingData []MetricsData

    // Verificar se o ficheiro existe e tem conteúdo válido
    if _, err := os.Stat(filename); err == nil {
        file, err := ioutil.ReadFile(filename)
        if err != nil {
            return fmt.Errorf("Erro ao ler o ficheiro JSON: %v", err)
        }
        if len(file) > 0 {
            if err := json.Unmarshal(file, &existingData); err != nil {
                log.Printf("JSON inválido detectado. Corrigindo...\n")
                existingData = []MetricsData{}
            }
        }
    }

    // Adicionar os novos dados
    existingData = append(existingData, data)

    // Serializar e escrever de volta no ficheiro
    file, err := json.MarshalIndent(existingData, "", "  ")
    if err != nil {
        return fmt.Errorf("Erro ao formatar o JSON: %v", err)
    }
    return ioutil.WriteFile(filename, file, 0644)
}

func main() {
    // Criar o diretório para armazenar o ficheiro JSON, caso não exista
    if _, err := os.Stat("prometheus_data"); os.IsNotExist(err) {
        os.Mkdir("prometheus_data", 0755)
    }

    // Configurar o handler para receber métricas
    http.HandleFunc("/receive_metrics", receiveMetricsHandler)

    // Iniciar o servidor para escutar requisições na porta 8001
    fmt.Println("Data Collector rodando na porta 8001...")
    http.ListenAndServe(":8001", nil)
}
