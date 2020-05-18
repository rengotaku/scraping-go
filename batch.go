package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	cron "github.com/robfig/cron/v3"
	models "github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
)

var (
	maxHistoryNum = 3
)

func main() {
	// cronExec()
	checkReserves()
}

func cronExec() {
	c := cron.New()
	c.AddFunc("@daily", checkReserves)

	fmt.Println(fmt.Sprintf("start cron - %s", time.Now().String()))

	c.Run()
}

func checkReserves() {
	executedAt := time.Now()
	fmt.Println(fmt.Sprintf("start checkReserves - %s", executedAt.String()))

	db, err := models.Connection()
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()

	reserves := []models.Reserve{}
	db.Find(&reserves, &models.Reserve{Model: gorm.Model{DeletedAt: nil}})

	fmt.Println(fmt.Sprintf("target records - %d", len(reserves)))

	// should split array such as using offset
	for _, reserve := range reserves {
		var jobHistories []models.JobHistory
		db.Order("id desc").Limit(maxHistoryNum).Find(&jobHistories, &models.JobHistory{ReserveID: reserve.ID})

		// When all of latest three record are error, maybe reserve data is old.
		errCnt := 0
		for _, jobHistory := range jobHistories {
			if jobHistory.StatusCode != 200 {
				errCnt++
			}
			if !jobHistory.IsNotice {
				errCnt++
			}
		}

		if errCnt >= maxHistoryNum && len(jobHistories) == errCnt {
			reserve.CancelReserve(db)
			continue
		}

		history := diffHtml(db, reserve, jobHistories[0])
		if history == nil {
			continue
		}

		reserve.LastExecutedAt = executedAt

		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&reserve).Error; err != nil {
				return err
			}

			if err := tx.Create(&history).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			fmt.Println(fmt.Sprintf("予期しないエラーが発生しました: %s", err))
			return
		}

	}
}

func diffHtml(db *gorm.DB, reserve models.Reserve, jobHistory models.JobHistory) *models.JobHistory {
	wa := lib.WebAnalyser{
		UserAgent: reserve.UserAgent,
		Url:       reserve.Url,
		Query:     reserve.HtmlSelector,
	}

	res := wa.Search()
	if res.StatusCode != 200 {
		return nil
	}

	formatedTarEle := ""

	if res.TargetElement == "" {
		reserve.CancelReserve(db)
		return nil
	}

	history := models.JobHistory{}
	history.StatusCode = res.StatusCode

	for _, line := range strings.Split(res.TargetElement, "\n") {
		formatedTarEle += strings.TrimSpace(line)
	}

	if formatedTarEle == jobHistory.Html {
		return nil
	}

	history.Html = formatedTarEle

	// HACK: go into reserve model
	switch notifer := reserve.Notifier; notifer {
	case 1:
		if lib.SendToSlack(reserve.NotifierValue, fmt.Sprintf("Scraping Notifer - %s に変更がありました。", reserve.Url)) {
			history.IsNotice = true
		} else {
			history.IsNotice = false
		}
	case 2:
		history.IsNotice = false
	default:
		history.IsNotice = false
	}

	history.Reserve = reserve

	return &history
}
