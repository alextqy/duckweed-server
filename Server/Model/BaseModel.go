package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

func ConnDB() (bool, string, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("../Dao.db"), &gorm.Config{})
	if err != nil {
		return false, err.Error(), nil
	} else {
		return true, "", db
	}
}
