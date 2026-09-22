# BUG-011: Self-Closing and Nested Tags Break HTML Text Extraction

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `crawler`
- **Affects:** [`crawler/crawler.go:60-70`](file:///Users/myroslavrozum/src/go-langdetector/crawler/crawler.go#L60-L70)

## Description
In `crawler/crawler.go`, the tokenizer tracks skipped elements using a single boolean flag `inSkipTag`:
```go
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
```
Two distinct issues exist:
1. **Self-Closing Tags:** If an HTML tag like `<script src="..." />` or `<style ... />` is parsed as `html.SelfClosingTagToken`, `inSkipTag` is set to `true`. Self-closing tags never have an `html.EndTagToken`. As a result, `inSkipTag` remains `true` for the entire remainder of the HTML document, causing all subsequent text in the page body to be dropped.
2. **Nested Tags:** If `<head>` contains a `<script>...</script>`, when `</script>` is encountered, `inSkipTag` is set to `false`. Any subsequent content inside `<head>` (e.g. `<title>Text</title>`) is parsed as body text.

## Impact
HTML pages with self-closing script or style tags return empty strings, preventing language profile generation for those articles.

## Suggested Fix
Do not treat self-closing tags as opening an unclosed block, and use a depth counter or stack for nested tags:
```go
case html.StartTagToken:
	token := tokenizer.Token()
	if token.Data == "script" || token.Data == "style" || token.Data == "head" {
		skipDepth++
	}
case html.EndTagToken:
	token := tokenizer.Token()
	if token.Data == "script" || token.Data == "style" || token.Data == "head" {
		if skipDepth > 0 {
			skipDepth--
		}
	}
```
