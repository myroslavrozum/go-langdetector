# BUG-002: Database Test Teardown Wipes Production Badger Data Directory

- **Status:** Resolved
- **Severity:** Critical
- **Component:** `db`
- **Affects:** [`db/db_test.go:11-42`](file:///Users/myroslavrozum/src/go-langdetector/db/db_test.go#L11-L42)

## Description
In `db/db_test.go`, the test database path is hardcoded to the production database path:
```go
var dirPath = "data/langdetector-badger-db"
```
During test teardown, `teardown()` runs:
```go
func teardown() {
	log.Println("Performing teardown after all tests...")
	store.Close()
	err := os.RemoveAll(dirPath)
	if err != nil {
		log.Printf("Failed to delete directory: %v\n", err)
		return
	}
}
```

## Impact
Running tests (`go test ./...`, `make test`, or `make build`) destroys the real application database directory (`data/langdetector-badger-db`), deleting all trained language profiles that have been accumulated.

## Steps to Reproduce
1. Train some profiles or place data in `data/langdetector-badger-db`.
2. Run `go test ./db`.
3. Check the filesystem; `data/langdetector-badger-db` has been deleted by `os.RemoveAll`.

## Suggested Fix
Use a temporary directory for tests via Go's `t.TempDir()` or `os.MkdirTemp("", "badger-test-*")`:
```go
tempDir, err := os.MkdirTemp("", "badger-test-*")
// use tempDir as dirPath
```
Or run Badger with `badger.DefaultOptions("").WithInMemory(true)` during unit tests.
