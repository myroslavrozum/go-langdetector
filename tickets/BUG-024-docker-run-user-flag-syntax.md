# BUG-024: Makefile Syntax Error in `docker-run` User Flag

- **Status:** Open
- **Severity:** Low
- **Component:** `build` / `Makefile`
- **Affects:** [`Makefile:57`](file:///Users/myroslavrozum/src/go-langdetector/Makefile#L57)

## Description
In `Makefile`:
```makefile
docker-run:
	docker run -d --rm -it \
		-u $(id -u):$(id -g) \
		-v $(CURDIR)/data:/app/data -p 8080:8080 \
		--name go-langdetector \
		go-langdetector:latest
```
In GNU Make syntax, `$(name)` expands a Make variable. There are no Make variables named `id -u` or `id -g`. As a result, `$(id -u)` and `$(id -g)` expand to empty strings, producing:
`-u :`

## Impact
Running `make docker-run` errors out because `-u :` is an invalid user/group specifier for Docker.

## Steps to Reproduce
Run `make docker-run`. Docker fails with invalid user specification.

## Suggested Fix
Use the Make `shell` function to evaluate shell commands:
```makefile
		-u $(shell id -u):$(shell id -g) \
```
