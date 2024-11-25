FROM grafana/k6:latest

# Adiciona o script de teste ao container
COPY load_test.js /scripts/load_test.js

# Define o diretório de trabalho
WORKDIR /scripts

# Comando padrão para rodar o k6 com o script de teste
ENTRYPOINT ["k6", "run", "/scripts/load_test.js"]
