# Etapa 1: Construir a aplicação
FROM golang:1.20-buster AS builder  

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar todos os ficheiros de código fonte para o container
COPY . .

# Baixar as dependências e compilar a aplicação
RUN go mod tidy && go build -o app

# Etapa 2: Criar a imagem final
FROM debian:latest  

# Instalar dependências necessárias
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Definir o diretório de trabalho no container final
WORKDIR /app

# Copiar o binário da aplicação para o container final
COPY --from=builder /app/app .

# Expor a porta que a aplicação está utilizando
EXPOSE 8080

# Comando padrão para rodar a aplicação web
CMD ["./app"]
