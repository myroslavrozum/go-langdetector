## 06162026: Because this is a fun project to hone my skills that I haven't been using
for a while TODO is to replace the ancient jQuery+bootstrap UI inherited from original
app with GoTTH (for no reason other than fun), then cover with tests and never return to it again leaving all those things which could've done better as they are

## 06202026 TODO
Before diving into WebUI, make logger websocket to broadcast messages.
According to our Gemini-friend the architecture should look like this:

Architecture:
- [ ] DesignHub: Holds room state mapping (map[string]map[*Client]bool) and manages thread-safe registration.
- [X] Client: Acts as the middleman holding a specific connection and a buffered outbound channel.
- [X] Write Pump: A dedicated per-client goroutine that handles sequential network writes to prevent write contention.
- [ ] Embed web assets into the bynary (emed.FS)