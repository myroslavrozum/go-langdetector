package webapp

import (
	"fmt"
	"go-langdetector/constants"
	"slices"
)

type Model struct {
	Title                      string
	SupportedLanguages         map[string]string
	SupportedLanguagesRendered string
	logger                     chan string
	Trigrammes                 map[string]map[string]float64
	Version                    string
}

func (m *Model) renderSupportedLanguages(detected ...string) {
	var dl string
	if len(detected) > 0 {
		dl = detected[0]
	}
	shortNames := make([]string, 0, len(constants.UrlDictionary))
	for k := range constants.UrlDictionary {
		shortNames = append(shortNames, k)
	}
	slices.Sort(shortNames)

	m.SupportedLanguages = make(map[string]string)

	m.SupportedLanguagesRendered = ``
	for _, shortName := range shortNames {
		fullName := constants.UrlDictionary[shortName][0]
		trainingUrl := constants.UrlDictionary[shortName][1]
		m.SupportedLanguages[shortName] = fullName
		class := `text-[#3b5998] hover:underline text-sm`
		if dl != "" && dl == shortName {
			class += " detected"
		}
		m.SupportedLanguagesRendered += fmt.Sprintf(`<li id="langmarker_" + %s class="%s"><a href="%s">%s</a></li>`, shortName, class, trainingUrl, fullName)
	}
}
