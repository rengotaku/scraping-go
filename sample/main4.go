package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

type confirmForm struct {
	Url   string `form:"url" binding:"required"`
	Query string `form:"query" binding:"required"`
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

//pipeline function
func filter(str string) string {
	return strings.Replace(str, "r", "R", -1)
}

func main() {
	router := gin.Default()
	router.Delims("{[{", "}]}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
		"filter":       filter,
	})

	router.HTMLRender = createRender()

	router.GET("/search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search", gin.H{})
	})

	router.POST("/search/confirm", func(c *gin.Context) {
		var form confirmForm
		if c.ShouldBind(&form) == nil {
			c.HTML(http.StatusOK, "confirm", gin.H{"form": form})
		}
	})

	router.Run(":8080")
}

func createRender() multitemplate.Render {
	r := multitemplate.New()
	r.AddFromFiles("search", "./app/views/layout/base.tmpl", "./app/views/search.tmpl")
	r.AddFromFiles("confirm", "./app/views/layout/base.tmpl", "./app/views/confirm.tmpl")

	return r
}
