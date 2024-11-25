# Etapa 1: Construir a aplicação Go
FROM golang:1.20-alpine AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar todos os ficheiros da pasta NODE para o container
COPY . .

# Baixar as dependências e compilar a aplicação
RUN go mod tidy && go build -o node_app

# Etapa 2: Imagem para executar o Node com node-exporter
FROM alpine:latest

# Instalar ferramentas necessárias, incluindo bash, curl, wget, e jq
RUN apk add --no-cache bash curl wget jq

# Definir o diretório de trabalho no container final
WORKDIR /app

# Copiar o binário compilado da aplicação Node da etapa anterior
COPY --from=builder /app/node_app .

# Baixar o node_exporter e torná-lo executável
RUN wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz \
    && tar -xzf node_exporter-1.3.1.linux-amd64.tar.gz \
    && mv node_exporter-1.3.1.linux-amd64/node_exporter /app/node_exporter \
    && rm -rf node_exporter-1.3.1.linux-amd64* \
    && chmod +x /app/node_exporter

# Copiar o script de supervisão
COPY node-exporter/supervise_node.sh /app/supervise_node.sh
RUN chmod +x /app/supervise_node.sh

# Expor as portas do node e do node-exporter
EXPOSE 8082 9100

# Comando padrão para rodar o script de supervisão, que inicia o node e o node-exporter
ENTRYPOINT ["/app/supervise_node.sh"]
