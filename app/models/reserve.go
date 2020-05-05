package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Reserve struct {
	gorm.Model

	UUID          string `gorm:"type:varchar(100);"`
	Url           string `gorm:"type:text;"`
	HtmlSelector  string `gorm:"type:text;"`
	Notifier      int
	NotifierValue string `gorm:"type:text;"`
	PreHtml       string `gorm:"type:text;"`
	Interval      int
	ExecutedAt    time.Time
}

func (m *Reserve) AfterCreate(scope *gorm.Scope) (err error) {
	if m.UUID == "" {
		seed := []byte(strconv.Itoa(int(m.ID)))
		sha := sha256.Sum256(seed)
		uuid := hex.EncodeToString(sha[:])
		scope.DB().Model(m).Update("UUID", uuid)
	}

	return
}

func (m *Reserve) GetNotifierAsString() sql.NullString {
	nullString := sql.NullString{
		Valid: false,
	}

	switch m.Notifier {
	case 1:
		nullString.String = "slack"
		nullString.Valid = true
	case 2:
		nullString.String = "email"
		nullString.Valid = true
	}

	return nullString
}

func (m *Reserve) SetNotifier(s string) bool {
	switch s {
	case "slack":
		m.Notifier = 1
		return true
	case "email":
		m.Notifier = 2
		return true
	}

	return false
}

func (m *Reserve) GetIntervalAsString() sql.NullString {
	nullString := sql.NullString{
		Valid: false,
	}

	switch m.Notifier {
	case 1:
		nullString.String = "day"
		nullString.Valid = true
	case 2:
		nullString.String = "harf_hour"
		nullString.Valid = true
	}

	return nullString
}

func (m *Reserve) SetInterval(s string) bool {
	switch s {
	case "day":
		m.Notifier = 1
		return true
	case "harf_hour":
		m.Notifier = 2
		return true
	}

	return false
}
