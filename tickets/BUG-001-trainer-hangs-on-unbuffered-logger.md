# BUG-001: Background Trainer Hangs Indefinitely on Unbuffered Logger Channel

- **Status:** Fixed
- **Severity:** Critical
- **Component:** `trainer` / `main`
- **Affects:** [`trainer/trainer.go:82`](file:///Users/myroslavrozum/src/go-langdetector/trainer/trainer.go#L82), [`main.go:19`](file:///Users/myroslavrozum/src/go-langdetector/main.go#L19)

## Description
The `logger` channel passed from `main()` into `trainer.Train` is an unbuffered channel (`make(chan string)`). During training, `trainer.Train` logs statistics to this channel via:
```go
logger <- m
```
If no SSE client is connected to `/logStream` (which is the default state whenever a user is not actively viewing the web page), there is no receiver on `logger`. As a result, the send blocks indefinitely.

## Impact
Training permanently freezes on the very first language processed. The trainer stops fetching subsequent languages, never writes updated trigrams to BadgerDB, and never runs subsequent training cycles.

## Steps to Reproduce
1. Start the application with `go run .`.
2. Do not open the web browser to `/` or `/logStream`.
3. Observe that `trainer.Train` halts after crawling the first language and never logs or updates further languages.

## Suggested Fix
Make sending to the log channel non-blocking so training continues regardless of whether an active web viewer is listening:
```go
select {
case logger <- m:
default:
}
```
Alternatively, implement a fan-out hub for broadcast subscribers.

## Resolution
Fixed in [`trainer/trainer.go:82-85`](file:///Users/myroslavrozum/src/go-langdetector/trainer/trainer.go#L82-L85) by replacing the blocking channel send `logger <- m` with a non-blocking `select` statement (`select { case logger <- m: default: }`).
