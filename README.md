# htlib.go

A comprehensive Go library for programmatic terminal control using [ht](https://github.com/andyk/ht) (Headless Terminal). This library provides both synchronous and asynchronous APIs for interacting with terminal applications, making it ideal for CLI automation, testing, and AI agent integration.

## Features

- 🚀 **Dual API Design**: Both synchronous (blocking) and asynchronous (channel-based) APIs
- 🔄 **Event Streaming**: Subscribe to terminal events (init, output, resize, snapshot)
- 🎯 **Context Support**: Full `context.Context` support for cancellation and timeouts
- ⌨️  **Rich Key Support**: Comprehensive key mapping including function keys, modifiers, and special keys
- 📸 **Snapshots**: Capture terminal state as text or raw VT100 sequences
- 🧪 **Testing Ready**: Perfect for CLI application testing and automation
- 🤖 **AI-Friendly**: Originally designed to make terminals accessible to LLMs
- 📦 **Zero Dependencies**: Uses only the Go standard library

## Installation

First, install the `ht` binary:

```bash
cargo install --git https://github.com/andyk/ht
```

Then install the Go library:

```bash
go get github.com/io41/htlib.go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/io41/htlib.go"
)

func main() {
    // Create a virtual terminal
    vt := htlib.New(htlib.DefaultConfig())

    // Start the terminal
    ctx := context.Background()
    if err := vt.Start(ctx); err != nil {
        log.Fatal(err)
    }
    defer vt.Close()

    // Wait for initialization
    <-vt.Events()

    // Send a command
    vt.Input(ctx, "echo 'Hello, World!'\n")

    // Get a snapshot
    snapshot, _ := vt.WaitForSnapshot(ctx)
    fmt.Println(snapshot.Text)
}
```

## API Overview

### Creating a Terminal

```go
// Use default configuration (bash, 120x40)
vt := htlib.New(htlib.DefaultConfig())

// Or customize
config := htlib.Config{
    Binary:   "/bin/zsh",
    Args:     []string{"-l"},
    Cols:     80,
    Rows:     24,
    HtBinary: "ht",
}
vt := htlib.New(config)
```

### Synchronous API

```go
ctx := context.Background()

// Send raw input
vt.Input(ctx, "ls -la\n")

// Send named keys
vt.SendKeys(ctx, "nano", htlib.KeyEnter)
vt.SendKeys(ctx, "Hello", htlib.KeySpace, "World")
vt.SendKeys(ctx, htlib.Ctrl('x'), "y", htlib.KeyEnter)

// Resize terminal
vt.Resize(ctx, 100, 30)

// Get snapshot (blocking)
snapshot, err := vt.WaitForSnapshot(ctx)
if err == nil {
    fmt.Println(snapshot.Text) // Rendered text
    fmt.Println(snapshot.Seq)  // Raw VT100 sequences
}
```

### Asynchronous API (Event Streaming)

```go
// Main event stream (includes all events)
for event := range vt.Events() {
    switch e := event.(type) {
    case htlib.InitEvent:
        fmt.Printf("PID: %d, Size: %dx%d\n", e.PID, e.Cols, e.Rows)
    case htlib.OutputEvent:
        fmt.Print(e.Seq)
    case htlib.ResizeEvent:
        fmt.Printf("Resized to %dx%d\n", e.Cols, e.Rows)
    case htlib.SnapshotEvent:
        fmt.Println(e.Text)
    }
}

// Create a subscriber (multiple subscribers supported)
subscriber := vt.Subscribe()
defer vt.Unsubscribe(subscriber)

go func() {
    for event := range subscriber {
        // Handle events independently
    }
}()
```

### Key Helpers

```go
// Control keys
vt.SendKeys(ctx, htlib.Ctrl('c'))  // Ctrl+C
vt.SendKeys(ctx, htlib.Ctrl('d'))  // Ctrl+D

// Alt keys
vt.SendKeys(ctx, htlib.Alt('x'))   // Alt+X

// Modifier combinations
vt.SendKeys(ctx, htlib.CtrlShift("Left"))  // Ctrl+Shift+Left
vt.SendKeys(ctx, htlib.AltShift("Right"))  // Alt+Shift+Right
vt.SendKeys(ctx, htlib.CtrlAlt("Up"))      // Ctrl+Alt+Up

// Named keys
vt.SendKeys(ctx,
    htlib.KeyEnter,
    htlib.KeyEscape,
    htlib.KeyTab,
    htlib.KeyBackspace,
    htlib.KeyF1,
    htlib.KeyF12,
    htlib.KeyPageUp,
    htlib.KeyPageDown,
    htlib.KeyHome,
    htlib.KeyEnd,
)
```

## Event Types

### InitEvent
Emitted once at startup with initial terminal state.

```go
type InitEvent struct {
    Cols int
    Rows int
    PID  int
    Seq  string    // Raw VT100 output
    Text string    // Rendered text view
    Time time.Time
}
```

### OutputEvent
Emitted when the terminal produces output.

```go
type OutputEvent struct {
    Seq  string    // Raw VT100 output
    Time time.Time
}
```

### ResizeEvent
Emitted when the terminal is resized.

```go
type ResizeEvent struct {
    Cols int
    Rows int
    Time time.Time
}
```

### SnapshotEvent
Emitted in response to `TakeSnapshot()`.

```go
type SnapshotEvent struct {
    Cols int
    Rows int
    Seq  string    // Raw VT100 output
    Text string    // Rendered text view
    Time time.Time
}
```

## Examples

The `examples/` directory contains complete working examples:

- **[basic.go](examples/basic.go)**: Simple terminal interaction and snapshots
- **[interactive.go](examples/interactive.go)**: Interactive terminal session with user input
- **[vim_automation.go](examples/vim_automation.go)**: Automate nano/vim editor
- **[event_streaming.go](examples/event_streaming.go)**: Multiple event subscribers and monitoring
- **[cli_testing.go](examples/cli_testing.go)**: CLI testing framework example

Run examples:

```bash
cd examples
go run basic.go
go run interactive.go
go run vim_automation.go
go run event_streaming.go
go run cli_testing.go
```

## Testing

The library includes comprehensive tests:

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -cover

# Run specific test
go test -run TestWaitForSnapshot
```

## Use Cases

### CLI Application Testing

```go
func TestMyCLI(t *testing.T) {
    vt := htlib.New(htlib.DefaultConfig())
    ctx := context.Background()
    vt.Start(ctx)
    defer vt.Close()

    <-vt.Events() // Wait for init

    vt.Input(ctx, "my-cli-tool --version\n")
    time.Sleep(100 * time.Millisecond)

    snapshot, _ := vt.WaitForSnapshot(ctx)
    if !strings.Contains(snapshot.Text, "v1.0.0") {
        t.Error("Version not found")
    }
}
```

### Interactive Application Automation

```go
// Automate an interactive installer
vt.Input(ctx, "make install\n")
time.Sleep(1 * time.Second)

// Answer prompts
vt.SendKeys(ctx, "y", htlib.KeyEnter)  // Confirm installation
time.Sleep(1 * time.Second)

vt.Input(ctx, "/usr/local\n")          // Installation path
time.Sleep(1 * time.Second)

vt.SendKeys(ctx, htlib.KeyEnter)       // Confirm
```

### Terminal Recording

```go
// Record terminal session
var recording []htlib.Event

subscriber := vt.Subscribe()
go func() {
    for event := range subscriber {
        recording = append(recording, event)
    }
}()

// ... perform actions ...

// Later: replay or analyze recording
for _, event := range recording {
    fmt.Printf("%s: %T\n", event.Time(), event)
}
```

## Context and Cancellation

All blocking operations support context cancellation:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

snapshot, err := vt.WaitForSnapshot(ctx)
if err == context.DeadlineExceeded {
    log.Println("Snapshot timed out")
}

// With manual cancellation
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel after 2 seconds
}()

err := vt.Input(ctx, "long-running-command\n")
```

## Error Handling

```go
// Predefined errors
var (
    ErrNotStarted     // Terminal not started yet
    ErrAlreadyStarted // Terminal already running
    ErrClosed         // Terminal closed
    ErrTimeout        // Operation timed out
    ErrInvalidEvent   // Invalid event received
    ErrProcessExited  // ht process exited
)

// Check errors
if err := vt.Start(ctx); err != nil {
    if errors.Is(err, htlib.ErrAlreadyStarted) {
        // Handle already started
    }
}

// Get accumulated errors
if err := vt.Err(); err != nil {
    log.Printf("Terminal error: %v", err)
}
```

## Configuration Options

```go
type Config struct {
    Binary   string   // Binary to run (default: /bin/bash)
    Args     []string // Arguments to pass to binary
    Size     string   // Terminal size "COLSxROWS" (default: 120x40)
    Cols     int      // Explicit columns (overrides Size)
    Rows     int      // Explicit rows (overrides Size)
    HtBinary string   // Path to ht binary (default: "ht")
    Env      []string // Additional environment variables
}
```

## Architecture

htlib.go wraps the `ht` binary, communicating via JSON over STDIN/STDOUT:

```
┌─────────────┐
│   Your App  │
└──────┬──────┘
       │ (Go API)
┌──────▼──────┐
│  htlib.go   │
└──────┬──────┘
       │ (JSON/STDIO)
┌──────▼──────┐
│  ht binary  │
└──────┬──────┘
       │ (PTY)
┌──────▼──────┐
│ bash/app    │
└─────────────┘
```

## Performance Considerations

- Event channels are buffered (100 events by default)
- Multiple subscribers share events efficiently
- Snapshots are on-demand to minimize overhead
- Output events stream in real-time

## Comparison with Python/TypeScript Versions

| Feature | htlib.go | htlib.py | htlib.ts |
|---------|----------|----------|----------|
| Sync API | ✅ | ✅ | ✅ |
| Async API | ✅ (channels) | ✅ (asyncio) | ✅ (promises) |
| Context support | ✅ | ❌ | ❌ |
| Type safety | ✅ (static) | ❌ (runtime) | ✅ (static) |
| Multiple subscribers | ✅ | ❌ | ❌ |
| Dependencies | 0 | asyncio | Node.js |

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `go test` passes
5. Submit a pull request

## License

Apache 2.0 - See LICENSE file for details.

## Related Projects

- [ht](https://github.com/andyk/ht) - The underlying headless terminal (Rust)
- [htlib.py](https://github.com/andyk/headlong/blob/main/packages/env/htlib.py) - Python wrapper
- [htlib.ts](https://github.com/andyk/headlong/blob/main/packages/env/htlib.ts) - TypeScript wrapper

## Credits

- Created for the [ht](https://github.com/andyk/ht) project by [@andyk](https://github.com/andyk)
- Go implementation with full sync/async support and context integration

## Support

- Report issues: [GitHub Issues](https://github.com/andyk/htlib.go/issues)
- Discussions: [GitHub Discussions](https://github.com/andyk/htlib.go/discussions)
- ht documentation: [ht README](https://github.com/andyk/ht#readme)
