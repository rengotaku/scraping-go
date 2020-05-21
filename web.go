package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	controllers "github.com/user/scraping-go/app/controllers"
	helpers "github.com/user/scraping-go/app/helpers"
	models "github.com/user/scraping-go/app/models"
	lib "github.com/user/scraping-go/lib"
)

const (
	SECRET_KEY_MIN_NUM = 128
)

var (
	secretKey = ""
	domain    = ""
	port      = ""
)

func init() {
	secretKey = os.Getenv("SECRET")
	if len(secretKey) < SECRET_KEY_MIN_NUM {
		if gin.Mode() == gin.ReleaseMode {
			panic("Missing `SECRET` for release, set this as 128 digit")
		}

		rand.Seed(time.Now().UTC().UnixNano())
		secretKey = srand(SECRET_KEY_MIN_NUM)
	}

	fmt.Printf("secretKey: %s\n", secretKey)

	domain = os.Getenv("DOMAIN")
	if domain == "" {
		panic("Missing `DOMAIN`, set domain like as `www.hogehoge.com`")
	}

	fmt.Printf("domain: %s\n", domain)

	port = os.Getenv("PORT")
	if port == "" {
		panic("Missing `PORT`, set port number")
	}
}

// HACK: move to lib
var alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// generates a random string of fixed size
func srand(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = alpha[rand.Intn(len(alpha))]
	}
	return string(buf)
}

func main() {
	router := gin.Default()

	// HACK: laod from config
	store := sessions.NewCookieStore([]byte(secretKey))
	router.Use(sessions.Sessions("scraping", store))
	// HACK: should split another file.
	router.Use(lib.CsrfMiddleware(lib.CsrfOptions{
		Secret: secretKey,
		ErrorFunc: func(c *gin.Context) {
			c.String(403, "Forbidden")
			c.Abort()
		},
	}))
	router.Use(lib.TimeTokenMiddleware(lib.TimeTokenOptions{
		ErrorFunc: func(c *gin.Context) {
			c.String(408, "Request Timeout")
			c.Abort()
		},
	}))
	router.Static("/static", "./public")

	router.HTMLRender = createRender()

	routing(router)

	err := models.InitMigration()
	if err != nil {
		panic(err)
		return
	}

	router.Run(":" + port)
}

func routing(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/search")
	})

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
		"htmlSafe":         helpers.HtmlSafe,
		"maskSlackWebHook": helpers.MaskSlackWebHook,
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
