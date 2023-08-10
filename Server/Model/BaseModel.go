package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnDB() (bool, string, *sql.DB) {
	db, err := sql.Open("sqlite3", "../Dao.db")
	if err != nil {
		log.Fatal(err.Error())
		return false, err.Error(), nil
	} else {
		return true, "", db
	}
}
