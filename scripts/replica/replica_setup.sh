#!/bin/bash

# Aguarde o master estar disponível
until pg_isready -h postgres-container -U replication_user; do
    echo "Aguardando o master estar disponível..."
    sleep 2
done

# Realiza o backup base
pg_basebackup -h postgres-container -D /var/lib/postgresql/data -U replication_user -P --wal-method=stream

# Cria o arquivo recovery.conf
cp /docker-entrypoint-initdb.d/recovery.conf /var/lib/postgresql/data/recovery.conf
