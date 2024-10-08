# Etapa 1: Construir a aplicação Go
FROM golang:1.20-alpine AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar todos os arquivos da pasta NODES para o container
COPY . .

# Baixar as dependências e compilar a aplicação
RUN go mod tidy && go build -o node_app

# Etapa 2: Imagem para executar o Node
FROM alpine:latest

# Definir o diretório de trabalho no container final
WORKDIR /app

# Copiar apenas o binário compilado da etapa anterior
COPY --from=builder /app/node_app .

# Definir a porta que será exposta
EXPOSE 8082

# Comando padrão para rodar o node
CMD ["./node_app"]
