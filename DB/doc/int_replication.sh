#createuser -U test -P -c 5 --replication reptest
#pg_basebackup -h postgres-container -p 5432 -U reptest -D /var/lib/postgresql/data -Fp -Xs -R -P


# login to postgres
psql --username=test test

#create a table
CREATE TABLE customers (firstname text, customer_id serial, date_created timestamp);

#show the table
\dt


docker exec -it postgres-2 bash

# login to postgres
psql --username=test test

#show the tables
\dt