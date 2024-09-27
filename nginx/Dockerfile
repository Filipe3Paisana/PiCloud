FROM nginx:latest

# Instalar o bash
RUN apt-get update && apt-get install -y bash

# Copiar o arquivo de configuração do Nginx
COPY nginx/nginx.conf /etc/nginx/nginx.conf

# Copiar o script wait-for-it.sh
COPY api_service/scripts/wait-for-it.sh /usr/local/bin/wait-for-it.sh

# Dar permissão de execução ao script
RUN chmod +x /usr/local/bin/wait-for-it.sh

# Definir o ponto de entrada
ENTRYPOINT []

# Comando padrão para iniciar o Nginx usando o wait-for-it.sh
CMD ["/usr/local/bin/wait-for-it.sh", "api:8080", "--", "nginx", "-g", "daemon off;"]
