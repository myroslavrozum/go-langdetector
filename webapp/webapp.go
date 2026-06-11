package webapp

import (
	"bytes"
	"go-langdetector/constants"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func serveIndexPage(c *gin.Context) {
	type Webdata struct {
		Title              string
		SupportedLanguages string
	}

	webdata := Webdata{
		"language Detectur",
		"",
	}

	i := 0
	for _, v := range constants.UrlDictionary {
		webdata.SupportedLanguages += v[0]
		i++
		if i < len(constants.UrlDictionary) {
			webdata.SupportedLanguages += " | "
		}
	}

	tf := template.Must(template.ParseGlob("./webapp/templates/*"))
	log.Println("Defined templates: ", tf.DefinedTemplates())

	var stringBuffer bytes.Buffer
	tf.ExecuteTemplate(&stringBuffer, "index.html.tpl", webdata)

	c.DataFromReader(http.StatusOK,
		int64(stringBuffer.Len()),
		"text/html; charset=utf-8",
		&stringBuffer,
		nil,
	)
}

func Run() {
	router := gin.Default()
	router.Static("/assets", "./webapp/assets")
	router.GET("/", serveIndexPage)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
