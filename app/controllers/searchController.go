package controllers

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"

	"github.com/gocolly/colly/v2"
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
	}

	myValidate = new(lib.MyValidate).InitValidate()
)

type SearchBaseTemplate struct {
}

func (t *SearchBaseTemplate) GetLayoutFile() string {
	return "base.tmpl"
}

type ConfirmForm struct {
	Url   string `form:"url" validate:"required" jaFieldName:"ユーアルエル"`
	Query string `form:"query" validate:"required" jaFieldName:"クエリー"`
	// Url   string `form:"url" binding:"required"`
	// Query string `form:"query" binding:"required"`
	// Url           string `form:"url" validate:"required"`
	// Query         string `form:"query" validate:"required"`
	TargetElement string `form:"target_element"`
	Text          string `form:"target_element_text"`
}

func SearchIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "search/index", gin.H{
		"form": ConfirmForm{},
	})
}

// func ListOfError(err error) (m map[string]string) {
// 	errs := err.(validator.ValidationErrors)

// 	for _, e := range errs {
// 		m[e.Field()] = e.Translate()
// 	}
// }

func SearchConfirm(c *gin.Context) {
	var form ConfirmForm

	c.ShouldBind(&form)

	if err := myValidate.Validate.Struct(form); err != nil {
		c.HTML(http.StatusOK, "search/index", gin.H{
			"form": form,
			"errs": myValidate.GetErrorMessages(err),
		})
		return
	}

	if !isValidURL(form.Url) {
		// formError.addMessage("URLが誤っています。")
		// validator.StructErrors
		// c.HTML(http.StatusOK, "search/index", gin.H{"form": form, "errs": validator.StructErrors})
		return
	}

	co := colly.NewCollector()
	co.UserAgent = c.GetHeader("User-Agent")

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

		var bs []byte
		pDom.Find("*").Each(func(i int, s *goquery.Selection) {
			if len(s.Text()) > 1 {
				// HACK:
				ts := strings.TrimSpace(s.Text())
				if len(ts) > 1 {
					bs = append(bs, ts...)
					bs = append(bs, []byte("\n")...)
				}
			}
		})
		form.Text = string(bs)
	})

	// extract status code
	co.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		// p.StatusCode = r.StatusCode
		// p.Body = r.Body
	})
	co.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		// p.StatusCode = r.StatusCode
	})

	co.Visit(form.Url)

	c.HTML(http.StatusOK, "search/confirm", gin.H{"form": form})
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
