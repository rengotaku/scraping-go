package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/user/scraping-go/lib"
)

type Reserve struct {
	gorm.Model

	// FIXME: primary key
	UUID           string `gorm:"type:varchar(100);"`
	Url            string `gorm:"type:text;"`
	HtmlSelector   string `gorm:"type:text;"`
	UserAgent      string
	Notifier       int
	NotifierValue  string `gorm:"type:text;"`
	Interval       int
	LastExecutedAt time.Time
	JobHistories   []JobHistory `gorm:"association_autoupdate:false;"`
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

func (m *Reserve) CancelReserve(db *gorm.DB) {
	db.Delete(m)

	switch notifer := m.Notifier; notifer {
	case 1:
		if !lib.SendToSlack(m.NotifierValue, fmt.Sprintf("Scraping Notifer - %s はスクレイピングできません。再度、登録をし直して下さい。", m.Url)) {
			// nothing
		}
	case 2:
		// nothing
	default:
		// output err to logger
	}
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
