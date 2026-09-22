# BUG-003: Data Race and Concurrent Map Writes in `/detect` Handler

- **Status:** Resolved
- **Severity:** Critical
- **Component:** `webapp`
- **Affects:** [`webapp/webhandlers.go:28-68`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webhandlers.go#L28-L68), [`webapp/model.go:18-42`](file:///Users/myroslavrozum/src/go-langdetector/webapp/model.go#L18-L42)

## Description
In `webapp/webapp.go`, `router.POST("/detect", Detect(model))` registers the detector handler. In `Detect(m Model)`:
```go
func Detect(m Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		...
		m.renderSupportedLanguages(minLang)
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(m.SupportedLanguagesRendered))
	}
}
```
`m` is captured in the closure. Calling `m.renderSupportedLanguages(minLang)` invokes a method with a pointer receiver:
```go
func (m *Model) renderSupportedLanguages(detected ...string) {
	...
	m.SupportedLanguages = make(map[string]string)
	m.SupportedLanguagesRendered = ``
	for _, shortName := range shortNames {
		...
		m.SupportedLanguages[shortName] = fullName
		m.SupportedLanguagesRendered += fmt.Sprintf(...)
	}
}
```

## Impact
Each incoming HTTP request runs in its own goroutine in Gin. When multiple users or requests access `/detect` concurrently, they mutate `m.SupportedLanguages` and `m.SupportedLanguagesRendered` simultaneously without synchronization. This triggers:
1. `fatal error: concurrent map writes`, instantly crashing the process.
2. Data races on `m.SupportedLanguagesRendered`.
3. Request contamination where one request's rendered result overrides another's.

## Steps to Reproduce
1. Start the web server.
2. Send concurrent POST requests to `/detect` using `hey`, `ab`, or `curl`.
3. Run with `go test -race` or inspect with race detector enabled.

## Suggested Fix
Refactor `renderSupportedLanguages` to be a pure function that does not mutate a shared receiver:
```go
func renderSupportedLanguages(detected string) string
```
Return the rendered HTML string directly to the handler instead of storing it on a shared `Model`.
