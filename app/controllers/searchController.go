package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
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

	completeMessgeFlag = "complete_messge_flag"
)

type SearchBaseTemplate struct {
}

func (t *SearchBaseTemplate) GetLayoutFile() string {
	return "base.tmpl"
}

func (t *SearchBaseTemplate) GetCssFile() string {
	return "base.tmpl"
}

type SearchForm struct {
	Url   string `form:"url" validate:"required" jaFieldName:"WebのURL"`
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

func SearchIndex(c *gin.Context) {
	if !lib.BeginOtt(c, 60*10) {
		c.String(403, "Forbidden")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "search/index", gin.H{
		"form":     SearchForm{},
		"messages": myValidate.GetErrorMessages(nil),
		"csrf":     lib.GetCsrfToken(c),
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
			"csrf":     lib.GetCsrfToken(c),
		})
		return
	}

	var resForm ConfirmSendForm = ConfirmSendForm{}
	resForm.ConfirmForm = form

	session := sessions.Default(c)
	notifier, ok := session.Get("notifier").(string)
	if ok && len(notifier) > 0 {
		resForm.Notifier = notifier
		resForm.NotifierValue = session.Get("notifier_value").(string)
	}

	c.HTML(http.StatusOK, "search/confirm", gin.H{
		"form":     resForm,
		"messages": myValidate.GetErrorMessages(nil),
		"csrf":     lib.GetCsrfToken(c),
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
			"csrf":     lib.GetCsrfToken(c),
		})
		return
	}

	formatedTarEle := ""
	for _, line := range strings.Split(form.TargetElement, "\n") {
		formatedTarEle += strings.TrimSpace(line)
	}

	// HACK: is it late to prevent from double post?
	uuid := lib.EndOtt(c)

	reserve := models.Reserve{
		Url:            form.Url,
		UUID:           uuid,
		HtmlSelector:   form.Query,
		NotifierValue:  form.NotifierValue,
		UserAgent:      c.GetHeader("User-Agent"), // Should relay this from search web site.
		LastExecutedAt: time.Now(),
	}

	if !reserve.SetNotifier(form.Notifier) {
		fmt.Println(fmt.Sprintf("不正なパラメーター %s が設定されました。", form.Notifier))
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	if !reserve.SetInterval(form.Interval) {
		fmt.Println(fmt.Sprintf("不正なパラメーター %s が設定されました。", form.Interval))
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	histories := []models.JobHistory{{
		StatusCode: 200,
		Html:       formatedTarEle,
		IsNotice:   true,
	}}
	reserve.JobHistories = histories

	db, err := models.Connection()
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}
	defer db.Close()
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&reserve).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusInternalServerError, "予期しないエラーが発生しました。")
		return
	}

	session := sessions.Default(c)
	session.Set("notifier", form.Notifier)
	session.Set("notifier_value", form.NotifierValue)
	session.AddFlash("1", completeMessgeFlag)
	session.Save()

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
		fmt.Println(err)
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

	var jobHistory models.JobHistory
	db.Order("id desc").First(&jobHistory, &models.JobHistory{ReserveID: reserve.ID})

	doc, err := htmlquery.Parse(strings.NewReader(jobHistory.Html))

	form := ConfirmSendForm{}
	form.Url = reserve.Url
	form.Query = reserve.HtmlSelector
	// HACK: formated html is better
	form.TargetElement = jobHistory.Html
	form.TargetElementText = htmlquery.InnerText(doc)
	form.Notifier = reserve.GetNotifierAsString().String
	form.NotifierValue = reserve.NotifierValue
	form.Interval = reserve.GetIntervalAsString().String
	form.ExecutedAt = reserve.LastExecutedAt

	session := sessions.Default(c)
	var message string
	if len(session.Flashes(completeMessgeFlag)) > 0 {
		message = "登録しました。"
		session.Save()
	}

	c.HTML(http.StatusOK, "search/finished", gin.H{
		"form":    form,
		"message": message,
	})
}

func validateConfirm(c *gin.Context, form *ConfirmForm) map[string]string {
	if err := myValidate.Validate.Struct(form); err != nil {
		return myValidate.GetErrorMessages(err)
	}

	if !isValidURL(form.Url) {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.WebのURL", "URLの形式が不正です。")
	}

	wa := lib.WebAnalyser{
		UserAgent: c.GetHeader("User-Agent"),
		Url:       form.Url,
		Query:     form.Query,
	}

	res := wa.Search()
	if res.StatusCode != 200 {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.WebのURL", "指定のサイトが開けません。")
	}
	if res.TargetElement == "" {
		return myValidate.PushErrorMessage(nil, "ConfirmForm.SearchForm.比較する要素", "該当する要素が存在しません。")
	}

	form.TargetElement = res.TargetElement
	form.TargetElementText = res.TargetElementText

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

		slackURL := "https://hooks.slack.com/services/"
		if !strings.HasPrefix(form.NotifierValue, slackURL) || form.NotifierValue == slackURL {
			return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法のデータ", "WebhookのURLの形式が不正です。")
		}

		if !lib.SendToSlack(form.NotifierValue, "Scraping Notifer - 確認用通知") {
			return myValidate.PushErrorMessage(nil, "ConfirmSendForm.FinishedForm.通知方法のデータ", "指定のWebhookに通知ができません。")
		}
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
