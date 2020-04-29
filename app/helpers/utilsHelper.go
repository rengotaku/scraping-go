package helpers

import (
	"fmt"
	"html/template"
	"strings"
	"time"
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
