package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnDB() (bool, string, *sql.DB) {
	db, err := sql.Open("sqlite3", "./Dao.db")
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", db
}
