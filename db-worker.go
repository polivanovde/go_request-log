package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	// 	"sync"
	// 	"time"

	_ "github.com/lib/pq"
)

const (
	PGUSER     = "postgres"
	PGPASSWORD = "uaQYs4E34q9k"
	PGHOST     = "localhost"
	PGPORT     = "5432"
	PGDATABASE = "postgres"
)

func initDb() (*sql.DB, error) {
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

	return conn, nil
}

func insert(conn *sql.DB, model ModelTransferJson) {
	transactionId := prepareTransaction(conn, model, "payment")

	_, err = conn.Query(`
        INSERT INTO users (uid, balance) VALUES ($1, $2)
        ON CONFLICT (uid) DO UPDATE SET balance=(select balance from users where uid = $1) + $2`, model.UID, model.Value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	} else {
		updateTransaction(conn, transactionId)
	}
}

func transfer(conn *sql.DB, model ModelTransferJson) {
	transactionId := prepareTransaction(conn, model, "transfer")

	var updatedId int
	err = conn.QueryRow(`
        UPDATE users SET balance = balance-$1 WHERE uid = $2 AND balance-$1 > 0 RETURNING id;
	`, model.Value, model.SenderID).Scan(&updatedId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	} else {
		_, err = conn.Query(`
			UPDATE users SET balance = balance+$1 WHERE uid = $2;
		`, model.Value, model.UID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		} else {
			updateTransaction(conn, transactionId)
		}
	}
}

func prepareTransaction(conn *sql.DB, model ModelTransferJson, transactionType string) int {
	//create transaction with false status
	var lastInsertId int
	err = conn.QueryRow(`
	INSERT INTO transactions (sender_id, recipient_id, transaction_sum, transaction_type, success) VALUES ($1, $2, $3, $4, false) RETURNING id;
	`, model.SenderID, model.UID, model.Value, transactionType).Scan(&lastInsertId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	}

	return lastInsertId
}

func updateTransaction(conn *sql.DB, lastInsertId int) {
	//update transaction on status true
	_, err = conn.Query(`
        UPDATE transactions SET success = true WHERE id = $1;
        `, lastInsertId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	}
}

func selectLastTransaction(conn *sql.DB) string {
	var createdAt string
	err = conn.QueryRow(`
	SELECT created_at FROM transactions ORDER BY created_at DESC LIMIT 1;
	`).Scan(&createdAt)
	if err != nil {
		fmt.Print("No transactions history\n")
	}

	return createdAt
}
