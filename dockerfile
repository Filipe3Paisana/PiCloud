# Usar a imagem oficial do Nginx como base
FROM nginx:alpine

# Definir o diretório de trabalho dentro do container
WORKDIR /usr/share/nginx/html

# Remover o conteúdo padrão do Nginx
RUN rm -rf ./*

# Copiar os arquivos do projeto para o diretório correto no container
COPY ./WEB/html/ /usr/share/nginx/html/
COPY ./WEB/css/ /usr/share/nginx/html/css/
COPY ./WEB/js/ /usr/share/nginx/html/js/

# Expor a porta 80 para acessar a aplicação
EXPOSE 80

# Comando para iniciar o Nginx
CMD ["nginx", "-g", "daemon off;"]
