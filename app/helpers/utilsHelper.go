package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/user/scraping-go/lib"
)

var (
	myValidate = new(lib.MyValidate).InitValidate()
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

//pipeline function
func Filter(str string) string {
	return strings.Replace(str, "r", "R", -1)
}

func HtmlSafe(html string) template.HTML {
	return template.HTML(html)
}

func MaskSlackWebHook(str string) string {
	baseUrl := "https://hooks.slack.com/services/"
	return baseUrl + "****/****" + str[len(str)-4:len(str)-1]
}

// func GetMessage(messages map[string]string, key string) string {
// 	return messages[key]
// }

// func GetMessage(messages map[string]string) string {
// 	fmt.Println(messages)
// 	return "few"
// }

// func HasMessage(err error, key string) bool {
// 	return myValidate.HasErrorMessage(err, key)
// }
