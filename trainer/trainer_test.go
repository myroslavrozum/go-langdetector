package trainer_test

import (
	"go-langdetector/trainer"
	"reflect"
	"testing"
)

func TestExtractTrigrammesFromText(t *testing.T) {
	var testTrigrammes = map[string]float64{
		"ell": 0.09090909090909091,
		"wor": 0.09090909090909091,
		"rld": 0.09090909090909091,
		"ld!": 0.09090909090909091,
		"o, ": 0.09090909090909091,
		" wo": 0.09090909090909091,
		"Hel": 0.09090909090909091,
		"llo": 0.09090909090909091,
		", w": 0.09090909090909091,
		"lo,": 0.09090909090909091,
		"orl": 0.09090909090909091,
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		text string
		want map[string]float64
	}{
		{`Hello, world`, `Hello, world!`, testTrigrammes},
		{`Empty string`, ``, make(map[string]float64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trainer.ExtractTrigrammesFromText(tt.text)
			// TODO: update the condition below to compare got with tt.want.
			trigrammesAreEqual := reflect.DeepEqual(got, tt.want)
			if !trigrammesAreEqual {
				t.Errorf("ExtractTrigrammesFromText() %v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
