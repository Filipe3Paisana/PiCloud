package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "time"
)

// Estrutura para armazenar a resposta da API do Prometheus
type PrometheusResponse struct {
    Status string `json:"status"`
    Data   struct {
        ResultType string `json:"resultType"`
        Result     []struct {
            Metric map[string]string `json:"metric"`
            Value  []interface{}      `json:"value"`
        } `json:"result"`
    } `json:"data"`
}

// Estrutura para armazenar todas as métricas em um único arquivo
type MetricsData struct {
    Timestamp int64                          `json:"timestamp"`
    Metrics   map[string]*PrometheusResponse `json:"metrics"`
}

func collectMetrics(query string) (*PrometheusResponse, error) {
    prometheusURL := os.Getenv("PROMETHEUS_URL")
    // Escapa a query para evitar problemas de URL encoding
    encodedQuery := url.QueryEscape(query)
    fullURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, encodedQuery)

    resp, err := http.Get(fullURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Log para capturar a resposta completa
    body, err := ioutil.ReadAll(resp.Body)
    fmt.Printf("Resposta completa do Prometheus para query '%s': %s\n", query, string(body))

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Erro na resposta do Prometheus: %s", string(body))
    }

    var result PrometheusResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }
    return &result, nil
}

func appendAllMetricsData(data MetricsData) error {
    filename := "prometheus_data/all_metrics.json"
    var existingData []MetricsData

    // Carregar dados existentes, se houver
    if _, err := os.Stat(filename); err == nil {
        file, err := ioutil.ReadFile(filename)
        if err != nil {
            return err
        }
        if err := json.Unmarshal(file, &existingData); err != nil {
            return err
        }
    }

    // Adicionar nova entrada
    existingData = append(existingData, data)

    file, err := json.MarshalIndent(existingData, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, file, 0644)
}

func main() {
    // Diretório para armazenar o arquivo JSON
    if _, err := os.Stat("prometheus_data"); os.IsNotExist(err) {
        os.Mkdir("prometheus_data", 0755)
    }

    queries := map[string]string{
        "uptime": 	"up{job=\"nodes\"}",
        "cpu": 		"100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[1m])) * 100)",
        "memory": 	"100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes))",
        "disk": 	"100 * (node_filesystem_size_bytes{fstype!=\"tmpfs\"} - node_filesystem_free_bytes{fstype!=\"tmpfs\"}) / node_filesystem_size_bytes{fstype!=\"tmpfs\"}",
    }

    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fmt.Println("Iniciando coleta de métricas...")

            allMetrics := MetricsData{
                Timestamp: time.Now().Unix(),
                Metrics:   make(map[string]*PrometheusResponse),
            }

            for metricName, query := range queries {
                fmt.Printf("Coletando dados de %s...\n", metricName)

                data, err := collectMetrics(query)
                if err != nil {
                    fmt.Printf("Erro ao coletar %s: %v\n", metricName, err)
                    continue
                }

                allMetrics.Metrics[metricName] = data
            }

            if err := appendAllMetricsData(allMetrics); err != nil {
                fmt.Printf("Erro ao salvar todas as métricas: %v\n", err)
            } else {
                fmt.Println("Todas as métricas salvas.")
            }
        }
    }
}
