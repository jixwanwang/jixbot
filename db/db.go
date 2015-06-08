package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// TODO: write stuff that resembles a database interface here.
// For when you finally use a database instead of text files, you lazy poop.

func New(host, port, name, user string) (*sql.DB, error) {
	pgConnect := fmt.Sprintf("dbname=%s user=%s host=%s port=%s sslmode=disable",
		name, user, host, port)
	db, err := sql.Open("postgres", pgConnect)

	if err != nil {
		log.Printf("couldn't connect to db: %s", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("couldn't ping db: %s", err.Error())
		return nil, err
	}

	return db, nil
}
