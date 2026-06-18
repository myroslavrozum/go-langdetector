package trainer

import (
	"fmt"
	"go-langdetector/constants"
	"go-langdetector/crawler"
	"go-langdetector/db"
	"log"
	"maps"
	"math"
	"slices"
	"time"
)

const trainInterval = 5

func sum[T int | float64](arr []T) (s T) {
	for _, v := range arr {
		s += v
	}
	return s
}

func ExtractTrigrammesFromText(text string) map[string]float64 {
	newTrigrammes := make(map[string]int)
	frequencies := make(map[string]float64)

	runes := []rune(text)
	for i := 0; i <= len(runes)-3; i++ {
		trigramme := string(runes[i : i+3])
		newTrigrammes[trigramme]++
	}

	numberOfTrigrammes := len(newTrigrammes)
	for trigramme, count := range newTrigrammes {
		frequencies[trigramme] = float64(count) / float64(numberOfTrigrammes)
	}

	return frequencies
}

func Train(store *db.Store, logger chan string) {
	c := make(chan string)
	defer close(c)

	for {
		updatedTrigrammes := make(map[string]map[string]float64)
		for lang, data := range constants.UrlDictionary {
			url := data[1]
			log.Printf("Fetching content for language: %s from URL: %s", lang, url)
			txtContent, err := crawler.GetTextFromURL(url)
			if err != nil {
				log.Printf("Failed to get text for %s: %v", lang, err)
				continue
			}
			trigrammes := ExtractTrigrammesFromText(txtContent)

			log.Printf("0. Got %d trigrammes for %s language...\n", len(trigrammes), data[0])
			storedTrigrammes, err := store.RestoreTrigrammes(lang)
			log.Println("1. Number of stored trigrammes: ", len(storedTrigrammes))
			if err == nil && storedTrigrammes != nil {
				numberOfTrigrammes := len(storedTrigrammes)
				freqs := slices.Collect(maps.Values(storedTrigrammes))
				avgFreq := sum(freqs) / float64(numberOfTrigrammes)
				minFreq := slices.Min(freqs)

				absolutes := make([]float64, 0, len(freqs))
				for _, x := range freqs {
					absolutes = append(absolutes, math.Abs(avgFreq-x))
				}

				dispersion := sum(absolutes) / float64(numberOfTrigrammes)
				lowFreqValues := make([]float64, 0)
				for _, x := range freqs {
					if x < (avgFreq - dispersion) {
						lowFreqValues = append(lowFreqValues, x)
					}
				}
				m := fmt.Sprintf("%s: Total values %d,  min_freq: %f, avg frequency: %f, dispersion %f, LF items %d\n",
					lang, numberOfTrigrammes, minFreq, avgFreq, dispersion, len(lowFreqValues))
				log.Println(m)
				logger <- m

				log.Printf("2. Calculating updated values:")
				for trigramme, newFreq := range trigrammes {
					originalFreq, exists := storedTrigrammes[trigramme]
					var updatedFreq float64
					if !exists {
						originalFreq = 0
						updatedFreq = (originalFreq + newFreq) / 2.0
					} else {
						updatedFreq = newFreq
					}
					storedTrigrammes[trigramme] = updatedFreq
				}
			} else {
				storedTrigrammes = trigrammes
			}
			updatedTrigrammes[lang] = storedTrigrammes
		}
		store.DumpTrigrammes(updatedTrigrammes)
		time.Sleep(trainInterval * time.Minute)
	}
}
