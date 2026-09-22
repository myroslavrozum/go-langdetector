# BUG-012: Point-to-Point Unbuffered Channel Fragments SSE Stream Across Clients

- **Status:** Resolved
- **Severity:** Medium
- **Component:** `webapp`
- **Affects:** [`webapp/webhandlers.go:79-100`](file:///Users/myroslavrozum/src/go-langdetector/webapp/webhandlers.go#L79-L100)

## Description
In `webapp/webhandlers.go`:
```go
func logStream(model Model) gin.HandlerFunc {
	return func(c *gin.Context) {
		...
		c.Stream(func(w io.Writer) bool {
			select {
			case <-c.Request.Context().Done():
				return false
			case l := <-model.logger:
				c.SSEvent("message", strings.TrimSpace(l))
				return true
			}
		})
	}
}
```
`model.logger` is a single Go channel. In Go, reading from a channel is point-to-point: each value is delivered to exactly one receiver.

## Impact
If two browser tabs or users open the application at the same time:
1. Messages from the trainer are randomly partitioned between the two SSE streams. One tab receives message 1, the other receives message 2.
2. Neither client sees the complete log stream.

## Steps to Reproduce
1. Open the web UI in two separate browser windows or curl connections to `/logStream`.
2. Observe that log lines alternate or distribute between the two clients.

## Suggested Fix
Implement a simple broadcast broker/hub where each SSE connection registers its own channel, and messages from the trainer are fanned out to all active subscribers.
