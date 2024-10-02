# Etapa 1: Construir a aplicação
FROM golang:1.20-alpine AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos de configuração e código
COPY go.mod ./
COPY app.go ./

# Baixar as dependências
RUN go mod tidy

# Compilar os aplicativos Go
RUN go build -o app app.go

# Etapa 2: Executar a aplicação
FROM alpine:latest

# Definir o diretório de trabalho no container final
WORKDIR /app

# Copiar apenas o binário da aplicação web para o container final
COPY --from=builder /app/app .

# Comando padrão para rodar a aplicação web
CMD ["./app"]

