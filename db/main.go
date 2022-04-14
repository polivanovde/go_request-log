package main

import (
	"database/sql"
	"log"
	"fmt"

	"github.com/pressly/goose"
	_ "github.com/lib/pq"
)

const (
    PGUSER = "postgres"
    PGPASSWORD = "uaQYs4E34q9k"
    PGHOST = "localhost"
    PGPORT = "5432"
    PGDATABASE = "postgres"
)

func main() {
	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
            PGUSER,
            PGPASSWORD,
            PGHOST,
            PGPORT,
            PGDATABASE,
        ))
    if err != nil {
        log.Fatal(err)
    }

	if err := goose.Up(conn, "."); err != nil {
		log.Fatalf("goose %v", err)
	}
}