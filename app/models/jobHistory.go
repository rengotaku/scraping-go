package models

import (
	"github.com/jinzhu/gorm"
)

type JobHistory struct {
	gorm.Model

	ReserveID  uint
	Reserve    Reserve
	Html       string `gorm:"type:text;"`
	StatusCode int
	IsNotice   bool
}
