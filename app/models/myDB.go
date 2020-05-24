package models

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var (
	dbFileName = ""
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dbFileName = dir + "/db/" + gin.Mode() + ".sqlite3"
}

func InitMigration() error {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&Reserve{})
	db.AutoMigrate(&JobHistory{})

	return nil
}

func Connection() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, err
	}

	if gin.Mode() != gin.ReleaseMode {
		db.LogMode(true)
	}

	return db, nil
}
