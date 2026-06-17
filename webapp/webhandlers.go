package webapp

import (
	"go-langdetector/algos"
	"go-langdetector/constants"
	"go-langdetector/crawler"
	"go-langdetector/trainer"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func serveIndexPage(c *gin.Context) {
	type Webdata struct {
		Title              string
		SupportedLanguages map[string]string
	}

	webdata := Webdata{
		Title:              "language Detectur",
		SupportedLanguages: make(map[string]string),
	}

	for shortName, v := range constants.UrlDictionary {
		fullName := v[0]
		webdata.SupportedLanguages[shortName] = fullName
	}

	c.HTML(http.StatusOK, "index.html.tpl", webdata)
}

func Detect(trigrammes map[string]map[string]float64) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		} else {
			data = contentToCheck
		}

		trigrammes2investigate := trainer.ExtractTrigrammesFromText(data)
		log.Printf("Got %d trigrammes for: \"%s\"", len(trigrammes2investigate), contentToCheck)

		distances := make(map[string]float64)
		var minD float64 = math.MaxFloat64
		var minLangFull string
		var minLang string
		for lang, v := range constants.UrlDictionary {
			d := algos.CalculateColsineDistances(trigrammes[lang], trigrammes2investigate)
			distances[lang] = d
			log.Printf("Calculated distance to %s is %f", v[0], d)
			if d < minD {
				minD = d
				minLangFull = v[0]
				minLang = lang
			}
		}
		log.Printf("Minimum distance (%f) is to %s", minD, minLangFull)

		c.JSON(http.StatusOK, gin.H{
			"status":      "posted",
			"minD":        minD,
			"minLangFull": minLangFull,
			"minLang":     minLang,
		})
	}
}
