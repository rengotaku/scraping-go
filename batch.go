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

func main() {
	// cronExec()
	checkReserves()
}

func cronExec() {
	c := cron.New()
	c.AddFunc("@hourly", checkReserves)
	// c.AddFunc("@daily", checkReserves)

	fmt.Println(fmt.Sprintf("start cron - %s", time.Now().String()))

	c.Run()
}

func checkReserves() {
	fmt.Println(fmt.Sprintf("start checkReserves - %s", time.Now().String()))

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
		db.Order("id desc").Limit(3).Find(&jobHistories, &models.JobHistory{ReserveID: reserve.ID})

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

		if errCnt >= 3 && len(jobHistories) == errCnt {
			db.Delete(&reserve)

			switch notifer := reserve.Notifier; notifer {
			case 1:
				if !lib.SendToSlack(reserve.NotifierValue, fmt.Sprintf("Scraping Notifer - %s はスクレイピングできません。再度、登録をし直して下さい。", reserve.Url)) {
					// nothing
				}
			case 2:
				// nothing
			default:
				// output err to logger
			}

			continue
		}

		history := models.JobHistory{
			Reserve: reserve,
		}

		wa := lib.WebAnalyser{
			UserAgent: reserve.UserAgent,
			Url:       reserve.Url,
			Query:     reserve.HtmlSelector,
		}

		res := wa.Search()
		history.StatusCode = res.StatusCode

		if history.StatusCode == 200 {
			formatedTarEle := ""

			if res.TargetElement != "" {
				for _, line := range strings.Split(res.TargetElement, "\n") {
					formatedTarEle += strings.TrimSpace(line)
				}

				if formatedTarEle != reserve.PreHtml {
					history.Html = formatedTarEle
					reserve.PreHtml = formatedTarEle

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
				}
			}
		}

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
