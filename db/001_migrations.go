package main

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up00001, Down00001)
}

func Up00001(tx *sql.Tx) error {
	_, err := tx.Exec(`
	create table users
    (
        id      serial
            constraint users_pk
                primary key,
        uid     integer,
        balance bigint
    );

    alter table users
        owner to postgres;

    create unique index users_id_uindex
        on users (id);

    create unique index users_uid_uindex
        on users (uid);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down00001(tx *sql.Tx) error {
	_, err := tx.Exec("drop table users; drop table transactions;")
	if err != nil {
		return err
	}
	return nil
}