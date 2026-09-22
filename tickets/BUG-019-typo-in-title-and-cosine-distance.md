# BUG-019: Typo in Application Title and Cosine Distance Function

- **Status:** Resolved
- **Severity:** Low
- **Component:** `webapp`
- **Affects:** [`webapp/webapp.go:28`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webapp.go#L28), [`webapp/calculateColsineDistances.go:7`](file:///Users/myroslavrozum/src/go-langdetector/webapp/calculateColsineDistances.go#L7)

## Description
Several typos exist in symbols, filenames, and user-facing text:
1. In `webapp/webapp.go:28`:
   ```go
   model.Title = "Language Detectur"
   ```
   "Detectur" instead of "Detector".
2. In `webapp/calculateColsineDistances.go`:
   - Filename: `calculateColsineDistances.go` ("Colsine")
   - Function name: `calculateColsineDistances` ("Colsine")
   - Variable name: `norm2comapare` ("comapare")
   - Header comment: `algorythm` ("algorythm")

## Impact
Degrades UI polish and code readability.

## Suggested Fix
1. Fix `model.Title = "Language Detector"`.
2. Rename `calculateColsineDistances` to `calculateCosineDistance` and fix internal typos.
