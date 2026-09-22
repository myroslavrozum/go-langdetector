# BUG-022: Makefile Subshell Recipe Lines Ignore Build Environment Variables

- **Status:** Open
- **Severity:** Medium
- **Component:** `build` / `Makefile`
- **Affects:** [`Makefile:8-12`](file:///Users/myroslavrozum/src/go-langdetector/Makefile#L8-L12), [`Makefile:26-29`](file:///Users/myroslavrozum/src/go-langdetector/Makefile#L26-L29)

## Description
In `Makefile`:
```makefile
build: makeversion generate test
	@echo "Building the executable...."
	CGO_ENABLED=0
	GOOS=darwin
	GOARCH=amd64

	go build -o ./bin/go-langdetector .
```
In GNU Make, each line of a target recipe executes in a new subshell. `CGO_ENABLED=0`, `GOOS=darwin`, and `GOARCH=amd64` each execute and terminate in separate isolated shells without affecting the subsequent `go build` command.

## Impact
The binary is built using the host system's default environment and compiler settings rather than the intended target architecture and CGO flags.

## Steps to Reproduce
Run `make build` on a Linux or ARM64 Mac system and observe `go env` or run `file ./bin/go-langdetector`. The binary will match host architecture instead of `darwin/amd64` with `CGO_ENABLED=0`.

## Suggested Fix
Prefix the environment variables on the same command line or export them globally:
```makefile
build: makeversion generate test
	@echo "Building the executable...."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/go-langdetector .
```
Apply the same fix to `build-static`.
