package main

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	// "github.com/go-playground/locales/ja_JP"
	// ut "github.com/go-playground/universal-translator"
	// "github.com/go-playground/validator/v10"
	// ja_translations "github.com/go-playground/validator/v10/translations/ja"
	// "gopkg.in/bluesuncorp/validator.v5"

	controllers "github.com/user/scraping-go/app/controllers"
	helpers "github.com/user/scraping-go/app/helpers"
)

func main() {
	router := gin.Default()
	router.Delims("{[{", "}]}")

	binding.Validator = new(defaultValidator)

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("bookabledate", bookableDate)
	// }

	router.HTMLRender = createRender()

	routing(router)

	router.Run(":8080")
}

// var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
// 	date, ok := fl.Field().Interface().(time.Time)
// 	if ok {
// 		today := time.Now()
// 		if today.After(date) {
// 			return false
// 		}
// 	}
// 	return true
// }

func routing(r *gin.Engine) {
	r.GET("/search", controllers.SearchIndex)
	r.POST("/search/confirm", controllers.SearchConfirm)
}

func createRender() multitemplate.Render {
	r := multitemplate.New()

	funcMap := template.FuncMap{
		// "formatAsDate": helpers.formatAsDate,
		// "filter":       helpers.filter,
		"htmlSafe": helpers.HtmlSafe,
	}

	AddFromFilesFuncs(r, "search/index", funcMap, controllers.SearchIndexTemplates...)
	AddFromFilesFuncs(r, "search/confirm", funcMap, controllers.SearchConfirmTemplates...)

	return r
}

// AddFromFilesFuncs supply add template from file callback func
// https://sourcegraph.com/github.com/gin-contrib/multitemplate/-/blob/multitemplate.go#L76
func AddFromFilesFuncs(r multitemplate.Render, name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}
