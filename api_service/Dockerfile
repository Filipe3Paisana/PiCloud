# Etapa 1: Construção do binário
FROM rust:latest AS builder

WORKDIR /app

# Copiar o arquivo Cargo.toml
COPY Cargo.toml ./

# Criar um diretório vazio para o build
RUN mkdir src
RUN echo "fn main() {}" > src/main.rs

# Baixar as dependências e gerar o Cargo.lock
RUN cargo fetch

# Copiar o restante do código
COPY . .

# Compilar o binário final
RUN cargo build --release

# Etapa 2: Criar a imagem final
FROM rust:latest

WORKDIR /app

# Instalar dependências necessárias
RUN apt-get update && apt-get install -y libpq-dev && rm -rf /var/lib/apt/lists/*

# Copiar o binário compilado
COPY --from=builder /app/target/release/api_service /usr/local/bin/api_service

EXPOSE 8080

CMD ["api_service"]
