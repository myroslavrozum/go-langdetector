package webapp

import (
	"go-langdetector/constants"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func serveIndexPage(c *gin.Context) {
	type Webdata struct {
		Title              string
		SupportedLanguages string
	}

	webdata := Webdata{
		Title: "language Detectur",
	}

	var langs []string
	for _, v := range constants.UrlDictionary {
		langs = append(langs, v[0])
	}
	webdata.SupportedLanguages = strings.Join(langs, " | ")

	c.HTML(http.StatusOK, "index.html.tpl", webdata)
}

func Run() {
	router := gin.Default()
	router.LoadHTMLGlob("./webapp/templates/*")
	router.Static("/assets", "./webapp/assets")
	router.GET("/", serveIndexPage)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
