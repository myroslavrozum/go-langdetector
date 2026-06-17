package crawler

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/html"
)

func fixUtf(r rune) rune {
	if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r)) {
		return -1
	}

	if r == utf8.RuneError {
		return -1
	}
	return r
}

func GetTextFromURL(url string) (string, error) {
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
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
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
				trimmed = strings.ReplaceAll(trimmed, "\n", " ")
				trimmed = strings.ReplaceAll(trimmed, "\t", " ")
				trimmed = strings.ReplaceAll(trimmed, "\r", " ")
				trimmed = strings.Map(fixUtf, trimmed)

				if len(trimmed) > 0 {
					trimmed += " "
					textBuilder.WriteString(trimmed)
				}
			}
		}
	}
}
