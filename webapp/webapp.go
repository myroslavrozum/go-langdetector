package webapp

import (
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

//go:generate go tool templ generate
//go:generate tailwindcss -i assets/css/styles_template.css -o assets/css/styles_compiled.css
func Run(store *db.Store, logger chan string, version string) {
	var model Model

	model.Title = "Language Detectur"
	model.SupportedLanguages = make(map[string]string)

	for shortName, v := range constants.UrlDictionary {
		fullName := v[0]
		model.SupportedLanguages[shortName] = fullName
	}

	model.Trigrammes = make(map[string]map[string]float64)
	for lang := range constants.UrlDictionary {
		var err error
		model.Trigrammes[lang], err = store.RestoreTrigrammes(lang)
		if err != nil {
			log.Println(err)
		}
	}

	model.Version = string(version)
	model.logger = logger

	router := gin.Default()
	router.Static("/assets", "./webapp/assets")
	router.Static("/js", "./webapp/js")

	router.GET("/", serveIndexPage(model))
	router.GET("/notailwind", serveIndexPage(model))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST(`/detect`, Detect(model))
	router.GET("/getTrainerLogs", wsServe(model))
	router.Run() // listen and serve on 0.0.0.0:8080
}
