package models

import (
	"github.com/jinzhu/gorm"
)

type JobHistory struct {
	gorm.Model

	ReserveID  uint
	Reserve    Reserve `gorm:"association_autoupdate:false;association_autocreate:false"`
	Html       string  `gorm:"type:text;"`
	StatusCode int
	IsNotice   bool
}
