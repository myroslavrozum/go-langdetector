# BUG-020: Crawler Test Hardcodes Port 8080 and Races With Socket Bind

- **Status:** Resolved
- **Severity:** High
- **Component:** `crawler`
- **Affects:** [`crawler/crawler_test.go:18-28`](file:///Users/myroslavrozum/src/go-langdetector/crawler/crawler_test.go#L18-L28)

## Description
In `crawler/crawler_test.go`:
```go
func startServer() {
	http.HandleFunc("/ping", hello)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func setup() {
	go startServer()
}
```
1. It hardcodes port `8080`—the exact same port used by the main web server. If the application or any local service is already running on port 8080, `log.Fatal` kills the entire test process.
2. `http.HandleFunc("/ping", hello)` mutates the global `http.DefaultServeMux`.
3. `go startServer()` starts an asynchronous goroutine without waiting for the socket listener to be ready before tests execute, introducing a startup race condition.

## Impact
Running `go test ./...` fails catastrophically if port 8080 is occupied, or intermittently due to the startup race.

## Steps to Reproduce
1. Start the main application (`go run .`) in one terminal.
2. Run `go test ./crawler` in another terminal.
3. Observe test failure: `ListenAndServe: bind: address already in use`.

## Suggested Fix
Replace the custom server with standard library `httptest.NewServer`:
```go
func TestGetTextFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(hello))
	defer server.Close()

	text, err := GetTextFromURL(server.URL + "/ping")
	...
}
```
Remove `setup`, `teardown`, and `startServer`.
