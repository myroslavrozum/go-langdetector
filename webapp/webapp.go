package webapp

import (
	"embed"
	"go-langdetector/constants"
	"go-langdetector/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func readModelFromContext(value any) Model {
	if model, ok := value.(Model); ok {
		return model
	}
	return Model{}
}

//go:embed css/*
var assetsCss embed.FS

//go:generate go tool templ generate
//go:generate tailwindcss -i css/styles_template.css -o css/styles_compiled.css
func Run(store *db.Store, logger chan string, version string) {
	var model Model

	model.Title = "Language Detectur"
	model.SupportedLanguages = make(map[string]string)
	model.Version = string(version)
	model.logger = logger

	model.renderSupportedLanguages()
	model.Trigrammes = make(map[string]map[string]float64)
	for lang := range constants.UrlDictionary {
		var err error
		model.Trigrammes[lang], err = store.RestoreTrigrammes(lang)
		if err != nil {
			log.Println(err)
		}
	}

	router := gin.Default()
	router.StaticFS("/assets", http.FS(assetsCss))

	router.GET("/", serveIndexPage(model))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST(`/detect`, Detect(model))
	router.GET("/logStream", logStream(model))

	router.Run() // listen and serve on 0.0.0.0:8080
}
