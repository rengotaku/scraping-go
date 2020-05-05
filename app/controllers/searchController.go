package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/gocolly/colly/v2"
	"github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
	"github.com/yosssi/gohtml"
)

var (
	SearchTemplates = []Template{
		Template{
			BaseTemplate: &SearchBaseTemplate{},
			Name:         "search/index",
			Files:        []string{"search.tmpl"},
		},
		Template{
			BaseTemplate: &SearchBaseTemplate{},
			Name:         "search/confirm",
			Files:        []string{"confirm.tmpl"},
		},
		Template{
			BaseTemplate: &SearchBaseTemplate{},
			Name:         "search/finished",
			Files:        []string{"finished.tmpl"},
		},
	}

	myValidate = lib.MyValidate{}.InitValidate()
)

type SearchBaseTemplate struct {
}

func (t *SearchBaseTemplate) GetLayoutFile() string {
	return "base.tmpl"
}

type SearchForm struct {
	Url   string `form:"url" validate:"required" jaFieldName:"サイトのURL"`
	Query string `form:"query" validate:"required" jaFieldName:"比較する要素"`
}

type ConfirmForm struct {
	SearchForm
	TargetElement     string `form:"target_element"`
	TargetElementText string `form:"target_element_text"`
}

type ConfirmSendForm struct {
	ConfirmForm
	FinishedForm
}

type FinishedForm struct {
	Notifier      string    `form:"notifier" validate:"required" jaFieldName:"通知方法"`
	NotifierValue string    `form:"notifier_value" validate:"required" jaFieldName:"通知方法のデータ"`
	Interval      string    `form:"interval" validate:"required" jaFieldName:"通知間隔"`
	ExecutedAt    time.Time `form:"executed_at" jaFieldName:"最終実行日"`
}

type SlackParams struct {
	Text string `json:"text"`
}

func SearchIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "search/index", gin.H{
		"form":     SearchForm{},
		"messages": myValidate.GetErrorMessages(nil),
	})
}

func SearchConfirm(c *gin.Context) {
	var form ConfirmForm

	c.ShouldBind(&form)
	messages := validateConfirm(c, &form)
	if len(messages) > 0 {
		c.HTML(http.StatusOK, "search/index", gin.H{
			"form":     form,
			"messages": messages,
		})
		return
	}

	var resForm ConfirmSendForm = ConfirmSendForm{}
	resForm.ConfirmForm = form

	c.HTML(http.StatusOK, "search/confirm", gin.H{
		"form":     resForm,
		"messages": myValidate.GetErrorMessages(nil),
	})
}

func SearchConfirmLast(c *gin.Context) {
	var form ConfirmSendForm

	c.ShouldBind(&form)
	messages := validateConfirmLast(c, &form)
	if len(messages) > 0 {
		c.HTML(http.StatusOK, "search/confirm", gin.H{
			"form":     form,
			"messages": messages,
		})
		return
	}

	formatedTarEle := ""
	for _, line := range strings.Split(form.TargetElement, "\n") {
		formatedTarEle += strings.TrimSpace(line)
	}

	reserve := models.Reserve{
		Url:           form.Url,
		HtmlSelector:  form.Query,
		NotifierValue: form.NotifierValue,
		PreHtml:       formatedTarEle,
		ExecutedAt:    time.Now(),
	}

	if !reserve.SetNotifier(form.Notifier) {
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	if !reserve.SetInterval(form.Interval) {
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	db, err := models.Connection()
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()
	db.Create(&reserve)

	// HACK: should use r.HandleContext(c) better
	path := fmt.Sprintf("/search/finish?reserved_key=%s", reserve.UUID)
	c.Redirect(http.StatusMovedPermanently, path)
}

func SearchFinished(c *gin.Context) {
	uuid := c.Query("reserved_key")
	if len(uuid) == 0 {
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	db, err := models.Connection()
	if err != nil {
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}
	defer db.Close()

	reserve := models.Reserve{}
	db.Find(&reserve, &models.Reserve{UUID: uuid, Model: gorm.Model{DeletedAt: nil}})
	if reserve.ID == 0 {
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	form := ConfirmSendForm{}
	form.Url = reserve.Url
	form.Query = reserve.HtmlSelector
	// HACK: formated html is better
	form.TargetElement = reserve.PreHtml
	// TODO: extra text from html
	// form.TargetElementText = reserve.PreHtml
	form.Notifier = reserve.GetNotifierAsString().String
	form.NotifierValue = reserve.NotifierValue
	form.Interval = reserve.GetIntervalAsString().String
	form.ExecutedAt = reserve.ExecutedAt

	c.HTML(http.StatusOK, "search/finished", gin.H{"form": form})
}

func validateConfirm(c *gin.Context, form *ConfirmForm) map[string]string {
	if err := myValidate.Validate.Struct(form); err != nil {
		return myValidate.GetErrorMessages(err)
	}

	if !isValidURL(form.Url) {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.サイトのURL", "URLの形式が不正です。")
	}

	co := colly.NewCollector()
	co.UserAgent = c.GetHeader("User-Agent")

	var statusCode int
	co.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find("script").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		e.DOM.Find("style").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		pDom := e.DOM.Find(form.Query).Parent()

		html, _ := pDom.Html()
		form.TargetElement = gohtml.Format(html)

		txt := pDom.Text()
		lines := []string{}
		for _, v := range strings.Split(txt, "\n") {
			tTxt := strings.TrimSpace(v)
			if len(tTxt) > 0 {
				lines = append(lines, tTxt)
			}
		}

		form.TargetElementText = strings.Join(lines, "\n")
	})

	// extract status code
	co.OnResponse(func(r *colly.Response) {
		// log.Println("response received", r.StatusCode)
		statusCode = r.StatusCode
	})
	co.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		// p.StatusCode = r.StatusCode
	})

	co.Visit(form.Url)

	if statusCode != 200 {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.サイトのURL", "指定のサイトが開けません。")
	}
	if form.TargetElement == "" {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.サイトのURL", "該当する要素が存在しません。")
	}

	return map[string]string{}
}

func validateConfirmLast(c *gin.Context, form *ConfirmSendForm) map[string]string {
	if err := myValidate.Validate.Struct(form); err != nil {
		return myValidate.GetErrorMessages(err)
	}

	switch notifer := form.Notifier; notifer {
	case "email":
		return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法", "通知方法が不正です。")
	case "slack":
		if !isValidURL(form.NotifierValue) {
			return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法のデータ", "WebhookのURLの形式が不正です。")
		}

		if !strings.HasPrefix(form.NotifierValue, "https://hooks.slack.com/services") {
			return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法のデータ", "WebhookのURLの形式が不正です。")
		}

		bjsonStr, _ := json.Marshal(SlackParams{Text: "Web diff - 確認用通知"})

		r, _ := http.NewRequest("POST", form.NotifierValue, bytes.NewBuffer(bjsonStr))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil || resp.StatusCode != 200 {
			return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法のデータ", "指定のWebhookに通知ができません。")
		}
		defer resp.Body.Close()
	default:
		return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法", "通知方法が不正です。")
	}

	return map[string]string{}
}

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
