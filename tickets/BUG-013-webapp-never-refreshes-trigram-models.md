# BUG-013: In-Memory Models Never Refreshed in Webapp After Startup

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `webapp`
- **Affects:** [`webapp/webapp.go:34-41`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webapp.go#L34-L41), [`webapp/webhandlers.go:45`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webhandlers.go#L45)

## Description
In `webapp/webapp.go`, language trigram profiles are loaded from BadgerDB exactly once at application startup:
```go
model.Trigrammes = make(map[string]map[string]float64)
for lang := range constants.UrlDictionary {
	var err error
	model.Trigrammes[lang], err = store.RestoreTrigrammes(lang)
	if err != nil {
		log.Println(err)
	}
}
```
Meanwhile, `trainer.Train` runs in the background and continuously writes updated trigram profiles to the database every 5 minutes. The web application never reloads or synchronizes its in-memory `model.Trigrammes`.

## Impact
1. If the application starts with a fresh/empty database, `model.Trigrammes[lang]` is `nil` for all languages. Even after training runs and saves complete models to the database, `/detect` continues using the empty in-memory models, rendering the detector unusable until the entire process is restarted.
2. Even with an existing database, the webapp never benefits from ongoing background training updates.

## Suggested Fix
Provide a synchronized store reader or thread-safe model cache that `/detect` queries, or reload trigram models periodically / on-demand using a `sync.RWMutex`.
