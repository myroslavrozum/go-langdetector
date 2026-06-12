package main

import (
	"errors"
	"go-langdetector/constants"
	"go-langdetector/db"
	"io"
	"log"
	"maps"
	"math"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"golang.org/x/net/html"
)

const TrainInterval = 5

func Sum[T int | float64](arr []T) (s T) {
	for _, v := range arr {
		s += v
	}
	return s
}

func getTextFromURL(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "MroCustomBot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	var textBuilder strings.Builder

	var inSkipTag bool

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if errors.Is(err, io.EOF) {
				return textBuilder.String(), nil
			}
			return "", err

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "script" || token.Data == "style" || token.Data == "head" {
				inSkipTag = true
			}

		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "script" || token.Data == "style" || token.Data == "head" {
				inSkipTag = false
			}

		case html.TextToken:
			if !inSkipTag {
				token := tokenizer.Token()
				trimmed := strings.TrimSpace(token.Data)

				if len(trimmed) > 0 {
					trimmed += " "
					textBuilder.WriteString(trimmed)
				}
			}
		}
	}
}

func extractTrigrammesFromText(text string) map[string]float64 {
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

func train(database *badger.DB) {
	for {
		updatedTrigrammes := make(map[string]map[string]float64)

		for lang, data := range constants.UrlDictionary {
			url := data[1]
			log.Printf("Fetching content for language: %s from URL: %s", lang, url)
			txtContent, err := getTextFromURL(url)
			if err != nil {
				log.Printf("Failed to get text for %s: %v", lang, err)
				continue
			}
			trigrammes := extractTrigrammesFromText(txtContent)

			log.Printf("0. Got %d trigrammes for %s language...\n", len(trigrammes), data[0])
			storedTrigrammes, err := db.RestoreTrigrammes(database, lang)
			log.Println("1. Number of stored trigrammes: ", len(storedTrigrammes))
			if err == nil && storedTrigrammes != nil {
				numberOfTrigrammes := len(storedTrigrammes)
				freqs := slices.Collect(maps.Values(storedTrigrammes))
				avgFreq := Sum(freqs) / float64(numberOfTrigrammes)
				minFreq := slices.Min(freqs)

				absolutes := make([]float64, 0, len(freqs))
				for _, x := range freqs {
					absolutes = append(absolutes, math.Abs(avgFreq-x))
				}

				dispersion := Sum(absolutes) / float64(numberOfTrigrammes)
				lowFreqValues := make([]float64, 0)
				for _, x := range freqs {
					if x < (avgFreq - dispersion) {
						lowFreqValues = append(lowFreqValues, x)
					}
				}
				log.Printf("%s: Total values %d,  min_freq: %f, avg frequency: %f, dispersion %f, LF items %d\n",
					lang, numberOfTrigrammes, minFreq, avgFreq, dispersion, len(lowFreqValues))
				log.Println("2. Calculating updated values:")
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
		db.DumpTrigrammes(database, updatedTrigrammes)
		time.Sleep(constants.TrainInterval * time.Minute)
	}
}
