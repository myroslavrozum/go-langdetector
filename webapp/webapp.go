package webapp

import (
	"go-langdetector/constants"
	"go-langdetector/crawler"
	"go-langdetector/db"
	"go-langdetector/trainer"
	"log"
	"maps"
	"math"
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

func calculateDistances(model map[string]float64, trigrammes2investigate map[string]float64) float64 {
	// var distance float64
	distances := make(map[string]float64)

	for k, observedFrequency := range trigrammes2investigate {
		modelFrequency := model[k]
		d := math.Abs(observedFrequency - modelFrequency)
		// distance += d
		distances[k] = d
	}
	var d float64
	for v := range maps.Values(distances) {
		d += v
	}
	return d
}

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

	router.POST(`/detect`, func(c *gin.Context) {
		contentToCheck := c.PostForm("content")

		data := ""

		if strings.HasPrefix(contentToCheck, `http://`) || strings.HasPrefix(contentToCheck, `https://`) {
			var err error
			data, err = crawler.GetTextFromURL(contentToCheck)
			if err != nil {
				log.Printf("Error getting data from URL %s: %v", contentToCheck, err)
				log.Println("Setting `data` to: ", contentToCheck)
				data = contentToCheck
			}
		}

		trigrammes2investigate := trainer.ExtractTrigrammesFromText(data)
		log.Printf("Got %d trigrammes for: \"%s\"", len(trigrammes2investigate), contentToCheck)

		distances := make(map[string]float64)
		var minD float64 = math.MaxFloat64
		var minLang string
		for lang, v := range constants.UrlDictionary {
			d := calculateDistances(trigrammes[lang], trigrammes2investigate)
			distances[lang] = d
			log.Printf("Calculated distance to %s is %f", v[0], d)
			if d < minD {
				minD = d
				minLang = v[0]
			}
		}
		log.Printf("Minimum distance (%f) is to %s", minD, minLang)

		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"minD":    minD,
			"minLang": minLang,
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080
}
