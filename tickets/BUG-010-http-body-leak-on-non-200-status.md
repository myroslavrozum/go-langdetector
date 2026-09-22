# BUG-010: HTTP Response Body Leak on Non-200 Status Codes in Crawler

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `crawler`
- **Affects:** [`crawler/crawler.go:35-43`](file:///Users/myroslavrozum/src/go-langdetector/crawler/crawler.go#L35-L43)

## Description
In `crawler/crawler.go`:
```go
resp, err := client.Do(req)
if err != nil {
	return "", err
}
if resp.StatusCode != 200 {
	return "", errors.New(resp.Status)
}

defer resp.Body.Close()
```
When `resp.StatusCode != 200`, the function returns immediately with an error. The `defer resp.Body.Close()` statement is positioned *after* this check, so it is never executed when the status code is not 200.

## Impact
Whenever a crawled URL returns a non-200 response (e.g. 429 Too Many Requests, 404 Not Found, 500 Internal Server Error), `resp.Body` remains unread and unclosed. In Go's `net/http`, this leaks file descriptors and prevents connection reuse in the HTTP client pool.

## Steps to Reproduce
1. Call `GetTextFromURL("https://httpbin.org/status/404")`.
2. Inspect open file descriptors; the socket connection remains unclosed.

## Suggested Fix
Move `defer resp.Body.Close()` immediately after the `err != nil` check:
```go
resp, err := client.Do(req)
if err != nil {
	return "", err
}
defer resp.Body.Close()

if resp.StatusCode != 200 {
	return "", errors.New(resp.Status)
}
```
