package main

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00002, Down00002)
}

func Up00002(tx *sql.Tx) error {
	_, err := tx.Exec(`
    create table transactions
    (
        id              serial
            constraint transactions_pk
                primary key,
        sender_id       integer,
        recipient_id    integer not null,
        transaction_type text,
        transaction_sum bigint,
        success         boolean,
    	created_at 		timestamp without time zone default (now() at time zone 'utc')
    );

    alter table transactions
        owner to postgres;

    create unique index transactions_id_uindex
        on transactions (id);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down00002(tx *sql.Tx) error {
	_, err := tx.Exec("drop table users; drop table transactions;")
	if err != nil {
		return err
	}
	return nil
}
