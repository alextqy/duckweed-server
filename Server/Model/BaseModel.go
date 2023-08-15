package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnDB() (bool, string, *sql.Tx, *sql.DB) {
	db, err := sql.Open("sqlite3", "../Dao.db")
	if err != nil {
		log.Fatal(err.Error())
		return false, err.Error(), nil, nil
	} else {
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err.Error())
			return false, err.Error(), nil, nil
		} else {
			return true, "", tx, db
		}
	}
}
