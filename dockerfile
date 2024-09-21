# Usar a imagem oficial do Nginx como base
FROM nginx:alpine

# Definir o diretório de trabalho dentro do container
WORKDIR /usr/share/nginx/html

# Copiar os arquivos da interface web (HTML, CSS, JS) para o diretório do Nginx
COPY ./ /usr/share/nginx/html

# Expor a porta 80 para acessar a aplicação
EXPOSE 80

# Comando para iniciar o Nginx
CMD ["nginx", "-g", "daemon off;"]
