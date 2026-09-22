# BUG-009: Random Language Falsely Detected on Short or Empty Input

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `webapp`
- **Affects:** [`webapp/webhandlers.go:49-65`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webhandlers.go#L49-L65)

## Description
In `webapp/webhandlers.go`, `Detect` iterates over all languages to find the minimum cosine distance:
```go
var minD float64 = math.MaxFloat64
var minLangFull string
var minLang string
for lang, v := range constants.UrlDictionary {
	d := calculateColsineDistances(learnedTrigrammes[lang], trigrammes2investigate)
	distances[lang] = d
	if d < minD {
		minD = d
		minLangFull = v[0]
		minLang = lang
	}
}
m.renderSupportedLanguages(minLang)
```
When input has fewer than 3 runes (e.g. empty string `""`, single character `"a"`, or two characters `"hi"`), `trigrammes2investigate` contains 0 trigrams. `calculateColsineDistances` returns `1.0` (the maximum possible distance, meaning 0 similarity) for all languages.

Because `1.0 < math.MaxFloat64` evaluates to `true` on the very first iteration, the first language iterated in `constants.UrlDictionary` is assigned to `minLang`. Because map iteration in Go is non-deterministic, a different language is randomly picked each time.

## Impact
Entering short or empty text falsely highlights a random language in the UI as detected, even though there was zero similarity.

## Steps to Reproduce
1. Submit empty text or "hi" to `/detect`.
2. Observe that a random language is highlighted as detected.
3. Submit again; a different language is highlighted.

## Suggested Fix
1. Reject inputs with fewer than 3 runes with a user-friendly message (e.g. "Input text is too short to detect language").
2. Only mark a language as detected if `minD < threshold` (e.g. `minD < 0.95` or when similarity $> 0$).
