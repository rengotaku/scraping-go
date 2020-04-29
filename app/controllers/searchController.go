package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gocolly/colly/v2"
	"github.com/yosssi/gohtml"
)

var (
	layout                 = "./app/views/layout/base.tmpl"
	SearchIndexTemplates   = []string{layout, "./app/views/search.tmpl"}
	SearchConfirmTemplates = []string{layout, "./app/views/confirm.tmpl"}
)

type ConfirmForm struct {
	// FormError

	Url   string `form:"url" binding:"required" validate:"required"`
	Query string `form:"query" binding:"required" validate:"required"`
	// Url   string `form:"url" binding:"required"`
	// Query string `form:"query" binding:"required"`
	// Url           string `form:"url" validate:"required"`
	// Query         string `form:"query" validate:"required"`
	TargetElement string `form:"target_element"`
	Text          string `form:"target_element_text"`
}

func SearchIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "search/index", gin.H{})
}

// func ListOfError(err error) (m map[string]string) {
// 	errs := err.(validator.ValidationErrors)

// 	for _, e := range errs {
// 		m[e.Field()] = e.Translate()
// 	}
// }

func TransFunc(ut ut.Translator, fe validator.FieldError) string {
	fld, _ := ut.T(fe.Field())
	t, err := ut.T(fe.Tag(), fld)
	if err != nil {
		return fe.(error).Error()
	}
	return t
}

func SearchConfirm(c *gin.Context) {
	var form ConfirmForm

	japanese := ja_JP.New()
	uni := ut.New(japanese, japanese)

	trans, _ := uni.GetTranslator("ja_JP")
	_ = trans.Add("ConfirmForm.Url", "フォームユーアルエル", false)
	_ = trans.Add("Url", "ユーアルエル", false)
	_ = trans.Add("Query", "クエリ", false)

	validate := validator.New()

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0}は必須項目です", false)
	}, TransFunc)

	// _ = trans.Add("Url", "ユーアルエル", false)
	// _ = ja_translations.RegisterDefaultTranslations(validate, trans)

	//
	// japanese := ja_JP.New()
	// uni := ut.New(japanese, japanese)
	// trans, _ := uni.GetTranslator("ja_JP")

	// validate := validator.New()
	// _ = ja_translations.RegisterDefaultTranslations(validate, trans)
	// _ = trans.Add("Url", "ユーアルエル", false)
	//

	// if err := c.BindQuery(&form); err != nil {
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }

	if err := c.ShouldBind(&form); err != nil {
		fmt.Println(form)
		fmt.Println(err)

		errs := err.(validator.ValidationErrors)
		fmt.Println(errs[0].Translate(trans))
		// fmt.Println(errs[0].Translate(ja))

		c.HTML(http.StatusOK, "search/index", gin.H{
			"form": form,
			"errs": err,
		})

		// c.String(http.StatusBadRequest, "bad request")
		return
	}

	// if err := validate.Struct(form); err != nil {
	// 	fmt.Println(form)

	// 	// if err := validate.Struct(form); err != nil {
	// 	// if err := c.ShouldBind(&form); err != nil {
	// 	errs := err.(validator.ValidationErrors)

	// 	fmt.Println(errs)

	// 	// fmt.Println(errs[0].Translate(trans))
	// 	// fmt.Println(errs[0].Field())

	// 	c.HTML(http.StatusOK, "search/index", gin.H{
	// 		"form": form,
	// 		"errs": errs,
	// 	})
	// 	return
	// }

	if !isValidURL(form.Url) {
		formError := new(FormError)
		formError.addMessage("URLが誤っています。")
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
