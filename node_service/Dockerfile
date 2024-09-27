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

# Copiar o binário compilado
COPY --from=builder /app/target/release/node_service /usr/local/bin/node_service

EXPOSE 8081

CMD ["node_service"]
