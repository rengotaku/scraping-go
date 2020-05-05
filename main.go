package main

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"

	controllers "github.com/user/scraping-go/app/controllers"
	helpers "github.com/user/scraping-go/app/helpers"
)

func main() {
	router := gin.Default()

	router.HTMLRender = createRender()

	routing(router)

	router.Run(":8080")
}

func routing(r *gin.Engine) {
	r.GET("/search", controllers.SearchIndex)
	r.POST("/search/confirm", controllers.SearchConfirm)
	r.POST("/search/last_check", controllers.SearchConfirmLast)
	r.GET("/search/finish", controllers.SearchFinished)
}

func createRender() multitemplate.Render {
	r := multitemplate.New()

	funcMap := template.FuncMap{
		// "formatAsDate": helpers.formatAsDate,
		// "filter":       helpers.filter,
		"htmlSafe": helpers.HtmlSafe,
		// "getErrorMessage": helpers.GetMessage,
		// "hasMessage":      helpers.HasMessage,
	}

	allT := []controllers.Template{}
	allT = append(allT, controllers.SearchTemplates...)

	for _, t := range allT {
		// Layout must be first
		addFromFilesFuncs(r, t.Name, funcMap, append([]string{t.GetFullLayoutes()}, t.GetFullViews()...)...)
	}

	return r
}

// AddFromFilesFuncs supply add template from file callback func
// https://sourcegraph.com/github.com/gin-contrib/multitemplate/-/blob/multitemplate.go#L76
func addFromFilesFuncs(r multitemplate.Render, name string, funcMap template.FuncMap, files ...string) {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
}
