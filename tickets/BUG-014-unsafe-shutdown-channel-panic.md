# BUG-014: Panic on Channel Send During Shutdown and Lack of Graceful Termination

- **Status:** Open
- **Severity:** Medium
- **Component:** `main`
- **Affects:** [`main.go:18-36`](file:///Users/myroslavrozum/src/go-langdetector/main.go#L18-L36)

## Description
In `main.go`:
```go
func main() {
	c := make(chan os.Signal, 1)
	logger := make(chan string)
	defer close(logger)
	defer close(c)

	signal.Notify(c, os.Interrupt)
	...
	go trainer.Train(store, logger)
	go webapp.Run(store, logger, langDetectorVersion)

	s := <-c
	log.Println("Got signal:", s)
}
```
1. `signal.Notify(c, os.Interrupt)` only listens for SIGINT (`os.Interrupt`), ignoring SIGTERM (`syscall.SIGTERM`), which is the default termination signal used by Docker, Kubernetes, and systemd.
2. When the signal is caught, `main()` exits immediately and executes deferred statements: `close(logger)` and `store.Close()`.
3. The background goroutines (`trainer.Train` and `webapp.Run`) are not canceled or waited upon with `context.Context` or `sync.WaitGroup`.

## Impact
If `trainer.Train` attempts to send to `logger` while `main` is closing the channel, Go panics with:
`panic: send on closed channel`
Furthermore, closing the Badger database while background writes or web requests are in flight can cause Badger transaction failures or partial write errors.

## Suggested Fix
1. Listen for both `os.Interrupt` and `syscall.SIGTERM`.
2. Use a root `context.WithCancel(context.Background())` to signal shutdown to the trainer and web server.
3. Wait for background goroutines to finish before closing resources.
