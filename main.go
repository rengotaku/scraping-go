package main

import (
	"html/template"
	"path/filepath"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	controllers "github.com/user/scraping-go/app/controllers"
	helpers "github.com/user/scraping-go/app/helpers"
	models "github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
)

func main() {
	router := gin.Default()

	// HACK: laod from config
	store := sessions.NewCookieStore([]byte("7c392fb14fe25f428f3194f59b5b01e1c6adf8702e41755abb774812de3238dc"))
	router.Use(sessions.Sessions("scraping", store))
	// HACK: should split another file.
	router.Use(lib.CsrfMiddleware(lib.CsrfOptions{
		Secret: "49da0b13f4aa987332efec012e370bf7",
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	router.HTMLRender = createRender()

	routing(router)

	err := models.InitMigration()
	if err != nil {
		panic(err)
		return
	}

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
		addFromFilesFuncs(r, t.Name, funcMap, append([]string{t.GetFullLayoute(), t.GetFullCss()}, t.GetFullViews()...)...)
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
