# BUG-008: Non-Standard Language Code for Ukrainian (`ua` vs ISO 639-1 `uk`)

- **Status:** Resolved
- **Severity:** Low
- **Component:** `constants`
- **Affects:** [`constants/constants.go:6`](file:///Users/myroslavrozum/src/go-langdetector/constants/constants.go#L6)

## Description
In `constants/constants.go`, the Ukrainian entry uses the key `ua`:
```go
`ua`: {`Ukrainian`, `https://uk.wikipedia.org/wiki/...`},
```
`ua` is the ISO 3166-1 alpha-2 country code for Ukraine. The ISO 639-1 two-letter code for the Ukrainian language is `uk` (as reflected in `uk.wikipedia.org`). All other entries in `UrlDictionary` use ISO 639-1 language codes (`en`, `fr`, `de`, `es`, `pl`, `pt`, `zh`, `ar`, `ru`, `eo`).

## Impact
Creates inconsistencies for any API consumers or downstream tools expecting standard ISO 639-1 language codes.

## Suggested Fix
Rename the key from `ua` to `uk` in `constants.UrlDictionary` and any database entries.
