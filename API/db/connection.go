package  db

import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
    connStr := "host=postgres-container port=5432 user=test password=test dbname=test sslmode=disable"

    var db *sql.DB
    var err error

    for {
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            fmt.Println("Erro a conectar à base de dados, tentando novamente em 5 segundos...")
            time.Sleep(5 * time.Second)
            continue
        }

        if err = db.Ping(); err == nil {
            break
        }
        fmt.Println("Database não está pronta, tentando novamente em 5 segundos...")
        time.Sleep(5 * time.Second)
    }

    return db, nil
}
