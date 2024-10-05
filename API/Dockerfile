# Etapa 1: Construir a aplicação
FROM golang:1.20-alpine AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar todos os arquivos de código fonte para o container
COPY . .

# Baixar as dependências e compilar a aplicação
RUN go mod tidy && go build -o app

# Etapa 2: Criar a imagem final
FROM alpine:latest

# Definir o diretório de trabalho no container final
WORKDIR /app

# Copiar o binário da aplicação para o container final
COPY --from=builder /app/app .

# Comando padrão para rodar a aplicação web
CMD ["./app"]
