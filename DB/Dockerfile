FROM postgres:latest

ENV POSTGRES_USER=${POSTGRES_USER}
ENV POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
ENV POSTGRES_DB=${POSTGRES_DB}
ENV PGDATA=/data

COPY ./postgres-1/config/pg_hba.conf /etc/postgresql/pg_hba.conf
COPY ./postgres-1/config/pg_ident.conf /etc/postgresql/pg_ident.conf
COPY ./postgres-1/config/postgresql.conf /etc/postgresql/postgresql.conf

COPY ./scripts/init/10_create_tables.sql /docker-entrypoint-initdb.d/
COPY ./scripts/init/init_replication.sh /docker-entrypoint-initdb.d/

EXPOSE 5432
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]
