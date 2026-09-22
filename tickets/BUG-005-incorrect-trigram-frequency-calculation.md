# BUG-005: Incorrect Trigram Relative Frequency Calculation Uses Unique Count

- **Status:** Open
- **Severity:** High
- **Component:** `trainer`
- **Affects:** [`trainer/trainer.go:24-40`](file:///Users/myroslavrozum/src/go-langdetector/trainer/trainer.go#L24-L40)

## Description
In `trainer/trainer.go`, `ExtractTrigrammesFromText` computes the frequency of each trigram:
```go
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
```
`len(newTrigrammes)` is the count of **distinct** (unique) trigrams in the text, not the total count of trigrams extracted.

## Impact
A proper relative frequency (empirical probability distribution) is defined as:
$$P(t) = \frac{\text{count}(t)}{\text{total trigrams}} = \frac{\text{count}(t)}{\sum_{t'} \text{count}(t')}$$
Because the denominator is `len(newTrigrammes)`:
1. Frequent trigrams in natural language texts obtain frequencies $> 1.0$ (for example, in `"aaaa"`, the only trigram `"aaa"` occurs twice, so `count / unique = 2 / 1 = 2.0`).
2. Across all unique trigrams, $\sum \text{frequencies} = \frac{\text{total trigrams}}{\text{distinct trigrams}} \gg 1.0$. In long Wikipedia texts, this sum can be 10, 50, or 100+, completely distorting the profile and cosine distance calculations.

## Steps to Reproduce
Run:
```go
freqs := trainer.ExtractTrigrammesFromText("aaaa")
// freqs["aaa"] == 2.0 (impossible for a frequency)
```

## Suggested Fix
Use the total number of trigrams extracted (`len(runes) - 2` when `len(runes) >= 3` or sum of counts) as the denominator:
```go
totalTrigrammes := len(runes) - 2
if totalTrigrammes <= 0 {
	return frequencies
}
for trigramme, count := range newTrigrammes {
	frequencies[trigramme] = float64(count) / float64(totalTrigrammes)
}
```
