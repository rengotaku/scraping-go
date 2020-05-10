package lib

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/yosssi/gohtml"
)

type WebAnalyser struct {
	UserAgent string
	Url       string
	Query     string
}

type AnalysedResponse struct {
	TargetElement     string
	TargetElementText string
	StatusCode        int
}

func (wa *WebAnalyser) Search() (res AnalysedResponse) {
	co := colly.NewCollector()
	co.UserAgent = wa.UserAgent

	co.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find("script").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		e.DOM.Find("style").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		pDom := e.DOM.Find(wa.Query).Parent()

		html, _ := pDom.Html()
		res.TargetElement = gohtml.Format(html)

		txt := pDom.Text()
		lines := []string{}
		for _, v := range strings.Split(txt, "\n") {
			tTxt := strings.TrimSpace(v)
			if len(tTxt) > 0 {
				lines = append(lines, tTxt)
			}
		}

		res.TargetElementText = strings.Join(lines, "\n")
	})

	// extract status code
	co.OnResponse(func(r *colly.Response) {
		// log.Println("response received", r.StatusCode)
		res.StatusCode = r.StatusCode
	})
	co.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		// p.StatusCode = r.StatusCode
	})

	co.Visit(wa.Url)

	return res
}
