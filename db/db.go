package db

import (
	"database/sql"
	"log"
)

var database *sql.DB

func Init(path string) {
	var err error
	database, err = initSQLLite(path)
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
	database.Close()
}
