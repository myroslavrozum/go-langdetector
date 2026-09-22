# BUG-004: Panic on Empty Slice in Trainer Statistics Calculation

- **Status:** Resolved
- **Severity:** High
- **Component:** `trainer`
- **Affects:** [`trainer/trainer.go:61-66`](file:///Users/myroslavrozum/src/go-langdetector/trainer/trainer.go#L61-L66)

## Description
In `trainer/trainer.go`, the statistical summary computation assumes `storedTrigrammes` contains at least one entry:
```go
if err == nil && storedTrigrammes != nil {
	numberOfTrigrammes := len(storedTrigrammes)
	freqs := slices.Collect(maps.Values(storedTrigrammes))
	avgFreq := sum(freqs) / float64(numberOfTrigrammes)
	minFreq := slices.Min(freqs)
```
If `storedTrigrammes` is an empty map (`len(storedTrigrammes) == 0`):
1. `slices.Min(freqs)` on an empty slice panics with:
   `panic: slices.Min: empty list`
2. `sum(freqs) / float64(numberOfTrigrammes)` performs division by zero (`0.0 / 0.0`), producing `NaN`.

## Impact
If an empty trigram map is ever stored in the database for a language (for example, if initial extraction produced no trigrams or an empty record was saved), the background training goroutine will panic and crash the entire program.

## Steps to Reproduce
1. Put an empty map for a language into Badger: `store.DumpTrigrammes(map[string]map[string]float64{"en": {}})`
2. Call `trainer.Train(store, logger)`
3. Observe panic in `slices.Min`.

## Suggested Fix
Guard the statistics calculation with `len(storedTrigrammes) > 0`:
```go
if err == nil && len(storedTrigrammes) > 0 {
	...
} else {
	storedTrigrammes = trigrammes
}
```
