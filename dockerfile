FROM nginx:alpine

RUN rm -rf /usr/share/nginx/html/*

COPY ./WEB/html /usr/share/nginx/html
COPY ./WEB/css /usr/share/nginx/html/css
COPY ./WEB/js /usr/share/nginx/html/js

EXPOSE 80
