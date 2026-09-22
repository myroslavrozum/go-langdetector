# BUG-023: `make clean` Deletes Machine-Wide Go Module Cache

- **Status:** Open
- **Severity:** Medium
- **Component:** `build` / `Makefile`
- **Affects:** [`Makefile:74`](file:///Users/myroslavrozum/src/go-langdetector/Makefile#L74)

## Description
In `Makefile`:
```makefile
clean: clean-bin clean-data clean-generated
	go clean -modcache
```
`go clean -modcache` deletes the system-wide Go module cache located at `$GOPATH/pkg/mod` (typically `~/go/pkg/mod`).

## Impact
Running `make clean` wipes out all downloaded Go modules across all projects on the developer's computer. Subsequent builds for this project and every other Go project on the machine must re-download all dependencies from the internet.

## Steps to Reproduce
Run `make clean`. Check `$GOPATH/pkg/mod`; it is completely empty.

## Suggested Fix
Remove `go clean -modcache` from the `clean` target. Use standard local clean actions (`clean-bin`, `clean-data`, `clean-generated`).
