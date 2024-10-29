#!/bin/bash

# Definir variáveis para o usuário e senha de replicação
REPL_USER="reptest"


# Criar o usuário de replicação com a senha definida
createuser -U test -P -c 5 --replication "$REPL_USER"  


# Iniciar o pg_basebackup para replicação
#pg_basebackup -h postgres-1-container -p 5432 -U reptest -D /data/ -Fp -Xs -R
