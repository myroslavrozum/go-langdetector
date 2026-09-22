# BUG-021: Dockerfile Stage 3 Fails on Non-Existent Asset and Template Paths

- **Status:** Open
- **Severity:** High
- **Component:** `docker`
- **Affects:** [`Dockerfile:32-33`](file:///Users/myroslavrozum/src/go-langdetector/Dockerfile#L32-L33)

## Description
In `Dockerfile`:
```dockerfile
#Stage 3: Create the minimal production runtime image
FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/webapp/assets webapp/assets
COPY --from=builder /app/webapp/templates webapp/templates
```
Neither `/app/webapp/assets` nor `/app/webapp/templates` exists in the codebase:
- CSS files are located in `webapp/css` and embedded into the Go binary via `//go:embed css/*`.
- Templates are compiled directly into Go code by Templ (`webapp/index_templ.go`).

## Impact
Running `docker build .` fails during Stage 3 with an error:
`COPY failed: stat /var/lib/docker/.../app/webapp/assets: file does not exist`.

## Steps to Reproduce
Run:
```bash
docker build -t go-langdetector .
```

## Suggested Fix
Remove lines 32 and 33 from `Dockerfile`. Ensure styles are generated before the binary build if running standalone.
