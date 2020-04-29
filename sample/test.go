package main

import (
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/gin"
)

func main() {
	templates := multitemplate.New()
	templates.AddFromFiles("index",
		"Base.html",
		"Navbar.html",
		"Index.html")

	templates.AddFromFiles("contact",
		"Base.html",
		"Navbar.html",
		"Contact.html")

	router := gin.New()
	router.HTMLRender = templates

	router.GET("", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{
			"title": "Home",
			"stuff": "Interesting home stuff",
		})
	})

	router.GET("/contact", func(c *gin.Context) {
		c.HTML(200, "contact", gin.H{
			"title": "Contact",
			"stuff": "Interesting contact stuff",
		})
	})
	router.Run(":8080")
}
