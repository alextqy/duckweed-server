package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnDB() (bool, string, *sql.DB) {
	//打开数据库，如果不存在，则创建
	db, err := sql.Open("sqlite3", "../Dao.db")
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", db
}
