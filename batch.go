package main

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	cron "github.com/robfig/cron/v3"
	models "github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
)

func main() {
	cronExec()
	// checkReserves()
}

func cronExec() {
	c := cron.New()
	c.AddFunc("@daily", checkReserves)

	c.Run()
}

func checkReserves() {
	db, err := models.Connection()
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()

	reserves := []models.Reserve{}
	db.Find(&reserves, &models.Reserve{Model: gorm.Model{DeletedAt: nil}})

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

		if errCnt > 0 && len(jobHistories) == errCnt {
			db.Delete(&reserve)

			switch notifer := reserve.Notifier; notifer {
			case 1:
				if !lib.SendToSlack(reserve.NotifierValue, fmt.Sprintf("Scraping Notifer - %s 指定のサイトがスクレイピングできません。再度、登録をし直して下さい。", reserve.Url)) {
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

		db.Create(&history)
		db.Save(&reserve)
	}
}
