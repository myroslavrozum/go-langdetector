package webapp

import (
	"go-langdetector/constants"
	"go-langdetector/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run(store *db.Store) {
	router := gin.Default()

	trigrammes := make(map[string]map[string]float64)
	for lang := range constants.UrlDictionary {
		var err error
		trigrammes[lang], err = store.RestoreTrigrammes(lang)
		if err != nil {
			log.Panic(err)
		}
	}

	router.LoadHTMLGlob("./webapp/templates/*")

	router.Static("/assets", "./webapp/assets")

	router.Static("/js", "./webapp/js")

	router.GET("/", serveIndexPage)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST(`/detect`, Detect(trigrammes))

	router.Run() // listen and serve on 0.0.0.0:8080
}
