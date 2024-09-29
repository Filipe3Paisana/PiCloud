# Dockerfile para o servidor Rocket 2

# Etapa 1: Usar uma imagem base do Rust para compilar
FROM rust:1.70 as builder

# Instalar dependências do sistema
RUN apt-get update && apt-get install -y libpq-dev

# Definir diretório de trabalho
WORKDIR /app

# Copiar os arquivos do projeto
COPY ./Cargo.toml ./Cargo.lock ./
COPY ./src ./src

# Construir o projeto
RUN cargo build --release

# Etapa 2: Criar a imagem final para execução
FROM debian:buster-slim

# Copiar o executável do servidor
COPY --from=builder /app/target/release/rocket_server2 /usr/local/bin/rocket_server2

# Expor a porta que o Rocket vai rodar
EXPOSE 8000

# Executar o servidor
CMD ["rocket_server2"]
