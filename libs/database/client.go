package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var dbc *sql.DB

func Connect(file string) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err.Error())
	}
	dbc = db
}

func GetInstance() *sql.DB {
	return dbc
}
