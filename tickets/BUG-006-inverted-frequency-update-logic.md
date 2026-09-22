# BUG-006: Inverted Frequency Update Logic Discards Trained History

- **Status:** Open
- **Severity:** High
- **Component:** `trainer`
- **Affects:** [`trainer/trainer.go:86-94`](file:///Users/myroslavrozum/src/go-langdetector/trainer/trainer.go#L86-L94)

## Description
In `trainer/trainer.go`, when merging new trigrams into `storedTrigrammes`, the update condition is inverted:
```go
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
```

## Impact
- If a trigram **does not exist** yet in stored memory (`!exists`), it is a brand-new trigram; but the code sets its initial frequency to `(0 + newFreq) / 2.0`, arbitrarily dividing it by 2.
- If a trigram **already exists** in stored memory (`exists`), the code sets `updatedFreq = newFreq`, completely discarding all prior training history and replacing it with the single new observation.

## Steps to Reproduce
1. Store a trigram with frequency `0.10`.
2. Run an update batch where that trigram has frequency `0.02`.
3. Inspect `storedTrigrammes`: the frequency becomes `0.02` instead of the average `0.06`.

## Suggested Fix
Invert the condition so existing trigrams are smoothed/averaged with prior history, while new trigrams receive their observed frequency:
```go
if exists {
	updatedFreq = (originalFreq + newFreq) / 2.0
} else {
	updatedFreq = newFreq
}
```
