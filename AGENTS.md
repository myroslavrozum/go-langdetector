# Agent Guidance

## Project Shape

- This is a Go language detector. It builds trigram frequency profiles from Wikipedia text and compares input with cosine distance.
- `main.go` opens the Badger store at `data/langdetector-badger-db`, starts training, and starts the web server.
- `crawler/` fetches and sanitizes article text; `trainer/` extracts trigrams and periodically persists language profiles; `db/` owns Badger access.
- `webapp/` owns Gin routes, the Templ UI, embedded CSS assets, detection, and the trainer log stream. `constants/` owns language metadata and source URLs.
- Keep package boundaries intact: domain data access belongs in `db`, profile generation in `trainer`, HTTP/UI behavior in `webapp`, and fetching/parsing in `crawler`.

## Commands

- Run the full test suite with `make test` or `go test ./...`.
- Regenerate generated sources with `make generate` or `go generate ./...`. This requires the `templ` and Tailwind CLIs.
- Run locally with `make run` or `go run .`; the server listens on port `8080` and uses the Badger data directory.
- Build with `make build`. It writes `.version`, generates sources, runs tests, and writes the binary to `bin/`.
- The module currently requires Go `1.26.4`; keep this aligned with `go.mod` and the Dockerfile.

## Change Conventions

- Use the module path `go-langdetector` for internal imports and keep Go formatting idiomatic (`gofmt`).
- Add or update focused tests beside the package being changed. Existing tests cover crawler, database, and trainer behavior; the webapp test file is currently sparse.
- Preserve Unicode handling: language text is processed as runes, not bytes, and trigram keys contain three runes.
- Generated `webapp/*_templ.go` and `webapp/css/*_compiled.css` files are ignored. Change `webapp/index.templ` or `webapp/css/styles_template.css`, then regenerate rather than editing generated output.
- Runtime profiles live under `data/`, and `.version` is build/runtime metadata. Do not treat either as source code.

## Things To Check

- Training runs in a goroutine and can block while sending logs if no `/logStream` consumer is connected; avoid introducing additional unbounded blocking in that path.
- The webapp loads profiles once at startup, while training updates Badger periodically. Changes involving model refresh or shared request state need concurrency tests.
- Network-dependent crawler/trainer behavior should use injected or local test servers where practical; avoid relying on Wikipedia in unit tests.
- Read [README.md](README.md) for setup and frontend tooling details and [TODO.md](TODO.md) for known project direction. The Dockerfile currently references paths that do not match the checked-in `webapp/css` layout; verify container changes against the current tree.
