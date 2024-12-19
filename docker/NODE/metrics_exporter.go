package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "bytes"
    "time"
    "strconv"
)

type PrometheusResponse struct {
    Status string `json:"status"`
    Data   struct {
        ResultType string `json:"resultType"`
        Result     []struct {
            Metric map[string]string `json:"metric"`
            Value  []interface{}     `json:"value"`
        } `json:"result"`
    } `json:"data"`
}

type Metric struct {
    Instance string  `json:"instance"`
    Value    float64 `json:"value"`
}

var instanceLabel string

func init() {
    instanceLabel = os.Getenv("INSTANCE_LABEL")
    if instanceLabel == "" {
        log.Fatal("A variável de ambiente INSTANCE_LABEL não está definida.")
    }
}

func fetchMetric(query string) ([]Metric, error) {
    prometheusURL := os.Getenv("PROMETHEUS_URL")
    if prometheusURL == "" {
        prometheusURL = "http://prometheus:9090"
    }
    encodedQuery := url.QueryEscape(query)
    fullURL := fmt.Sprintf("%s/api/v1/query?query=%s", prometheusURL, encodedQuery)
    log.Printf("Consultando o Prometheus na URL: %s", fullURL)
    resp, err := http.Get(fullURL)
    if err != nil {
        log.Printf("Erro ao realizar a consulta %s: %v", query, err)
        return nil, err
    }
    defer resp.Body.Close()

    // Ler a resposta
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Erro ao ler a resposta da consulta %s: %v", query, err)
        return nil, err
    }
    log.Printf("Resposta do Prometheus para a consulta '%s': %s", query, string(bodyBytes))

    // Decodificar a resposta
    var result PrometheusResponse
    if err := json.Unmarshal(bodyBytes, &result); err != nil {
        log.Printf("Erro ao decodificar resposta da consulta %s: %v", query, err)
        return nil, err
    }

    // Verificar se a consulta retornou resultados
    if len(result.Data.Result) == 0 {
        log.Printf("A consulta %s não retornou resultados", query)
        return nil, fmt.Errorf("A consulta %s não retornou resultados", query)
    }

    metrics := make([]Metric, 0)
    for _, r := range result.Data.Result {
        instance := r.Metric["instance"]
        valueStr, ok := r.Value[1].(string)
        if !ok {
            log.Printf("Formato inesperado no campo 'value': %v", r.Value)
            continue
        }
        parsedValue, err := strconv.ParseFloat(valueStr, 64)
        if err != nil {
            log.Printf("Erro ao converter valor %s para float: %v", valueStr, err)
            continue
        }
        metrics = append(metrics, Metric{Instance: instance, Value: parsedValue})
    }

    return metrics, nil
}

func sendMetricsToCollector(metrics map[string][]Metric) error {
    collectorURL := os.Getenv("DATA_COLLECTOR_URL")
    if collectorURL == "" {
        collectorURL = "http://data_collector:8001/receive_metrics"
    }

    payload := map[string]interface{}{
        "node_id": os.Getenv("NODE_ID"),
        "metrics": metrics,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // Log do payload
    log.Printf("Enviando métricas para o collector: %s", string(jsonData))

    resp, err := http.Post(collectorURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := ioutil.ReadAll(resp.Body)
        return fmt.Errorf("Recebeu status code %d do data_collector: %s", resp.StatusCode, string(bodyBytes))
    }
    return nil
}

func gatherMetrics() map[string][]Metric {
    metrics := make(map[string][]Metric)

    cpuQuery := fmt.Sprintf("100 - (avg by (instance) (rate(node_cpu_seconds_total{mode=\"idle\", instance=\"%s\"}[1m])) * 100)", instanceLabel)
    memQuery := fmt.Sprintf("100 * (1 - (node_memory_MemAvailable_bytes{instance=\"%s\"} / node_memory_MemTotal_bytes{instance=\"%s\"}))", instanceLabel, instanceLabel)
    diskQuery := fmt.Sprintf("100 * ((node_filesystem_size_bytes{fstype!=\"tmpfs\", instance=\"%s\"} - node_filesystem_free_bytes{fstype!=\"tmpfs\", instance=\"%s\"}) / node_filesystem_size_bytes{fstype!=\"tmpfs\", instance=\"%s\"})", instanceLabel, instanceLabel, instanceLabel)
    responseTimeQuery := fmt.Sprintf("probe_duration_seconds{job='blackbox', instance=\"%s\"}", instanceLabel)
    
    cpuMetrics, err := fetchMetric(cpuQuery)
    if err != nil {
        log.Printf("Erro ao coletar métricas de CPU: %v", err)
    } else {
        metrics["CPU"] = cpuMetrics
    }

    memMetrics, err := fetchMetric(memQuery)
    if err != nil {
        log.Printf("Erro ao coletar métricas de Memória: %v", err)
    } else {
        metrics["Memory"] = memMetrics
    }

    diskMetrics, err := fetchMetric(diskQuery)
    if err != nil {
        log.Printf("Erro ao coletar métricas de Disco: %v", err)
    } else {
        metrics["Disk"] = diskMetrics
    }

    responseTimeMetrics, err := fetchMetric(responseTimeQuery)
    if err != nil {
        log.Printf("Erro ao coletar métricas de Tempo de Resposta: %v", err)
    } else {
        metrics["ResponseTime"] = responseTimeMetrics
    }

    return metrics
}

func startMetricsExporter() {
    for {
        metrics := gatherMetrics()

        // Log das métricas coletadas
        log.Printf("Métricas coletadas no nó %s: %+v", os.Getenv("NODE_ID"), metrics)

        if err := sendMetricsToCollector(metrics); err != nil {
            log.Println("Erro ao enviar métricas:", err)
        }
        time.Sleep(5 * time.Minute)
    }
}
