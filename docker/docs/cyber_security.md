# Tópicos para Segurança

### docker scout
 - para ver as vulnerabilidades das imagens e retirar o que não está a ser utilizado, criando uma imagem adaptada para o nosso projeto e mais segura. 
 - container as non-root, ou seja, o container não tem todas as permissões, não sendo possível comprometer o sistema

### Kubernetes
 - implementar encriptação entre pods
 - restringir as permissões dos pods

### fragmentos
 - fazendo um hash dos fragmentos, conseguimos garantir a integridade dos mesmos para reconstruir o ficheiro corretamente. 
 - encriptar informação dos fragmentos, para que não sejam acedidos pelo host do node. 

### CloudFlare
 - para prevenir ataques de ddos ao datacenter principal. 

### Ecriptação da Base de Dados
 - encriptar as passwords dos utilizadores. 
