package webapp

type Model struct {
	Title              string
	SupportedLanguages map[string]string
	logger             chan string
	Trigrammes         map[string]map[string]float64
	Version            string
}
