# BUG-018: Negative Cosine Distance From Floating-Point Roundoff

- **Status:** Resolved
- **Severity:** Low
- **Component:** `webapp`
- **Affects:** [`webapp/calculateColsineDistances.go:25-26`](file:///Users/myroslavrozum/src/go-langdetector/webapp/calculateColsineDistances.go#L25-L26)

## Description
In `webapp/calculateColsineDistances.go`:
```go
cosineSimilarity := dotProduct / (math.Sqrt(normModel) * math.Sqrt(norm2comapare))
return 1.0 - cosineSimilarity
```
Due to standard IEEE-754 floating-point inaccuracies, when comparing two identical or nearly identical vectors, `cosineSimilarity` can evaluate to a value slightly greater than 1.0 (e.g. `1.0000000000000002`). As a result, `1.0 - cosineSimilarity` evaluates to a negative number (e.g. `-2.220446049250313e-16`).

## Impact
Cosine distance is mathematically bounded in $[0.0, 2.0]$. Negative distances violate distance metric axioms and can cause subtle bugs if other logic expects non-negative distances.

## Suggested Fix
Clamp `cosineSimilarity` to $[-1.0, 1.0]$:
```go
if cosineSimilarity > 1.0 {
	cosineSimilarity = 1.0
} else if cosineSimilarity < -1.0 {
	cosineSimilarity = -1.0
}
return 1.0 - cosineSimilarity
```
