# PICloud

O **PICloud** é um projeto que utiliza Docker, Nginx e K6 para simular e testar a carga de servidores web distribuídos. O Nginx atua como um balanceador de carga, enquanto o K6 é usado para realizar testes de desempenho, fornecendo métricas de tempo de resposta e taxas de sucesso das requisições.

## Execução

Para executar o projeto, siga os seguintes passos:

1. Clone o repositório.
2. Execute os containers usando Docker Compose.
3. Os testes de carga serão executados automaticamente com o K6.
4. Para encerrar os containers, utilize `docker-compose down`.

## Configuração do Nginx

O Nginx é configurado para balancear o tráfego entre múltiplos servidores, garantindo a distribuição eficiente das requisições.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).
