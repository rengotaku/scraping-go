package models

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var (
	dbFileName = "reserves.sqlite3"
)

func InitMigration() error {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&Reserve{})

	return nil
}

func Connection() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}