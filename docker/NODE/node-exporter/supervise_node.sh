#!/bin/bash

# Define o endpoint do data_collector
DATA_COLLECTOR_URL="http://data_collector:8001/receive_metrics"

# Inicia o node_app em segundo plano
./node_app &

# Inicia o node-exporter em segundo plano
/app/node_exporter &

# Função para coletar e enviar as métricas
send_metrics() {
    # Coletar métricas do node-exporter local
    metrics=$(curl -s http://localhost:9100/metrics | jq -Rsa .)

    # Enviar as métricas para o data_collector
    curl -X POST -H "Content-Type: application/json" -d '{
        "node_id": "'"$NODE_ID"'",
        "metrics": '"$metrics"',
        "timestamp": '$(date +%s)'
    }' $DATA_COLLECTOR_URL
}

# Loop infinito para enviar as métricas periodicamente
while true; do
    send_metrics
    sleep 60  # Enviar a cada 60 segundos
done
