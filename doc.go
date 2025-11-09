// Package htlib provides a comprehensive Go library for programmatic terminal control
// using ht (Headless Terminal).
//
// htlib wraps the ht binary to provide both synchronous and asynchronous APIs for
// interacting with terminal applications. It's ideal for CLI automation, testing,
// and AI agent integration.
//
// # Quick Start
//
// Create and start a virtual terminal:
//
//	vt := htlib.New(htlib.DefaultConfig())
//	ctx := context.Background()
//	if err := vt.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer vt.Close()
//
// # Synchronous API
//
// Send input and get snapshots:
//
//	vt.Input(ctx, "echo 'Hello, World!'\n")
//	snapshot, err := vt.WaitForSnapshot(ctx)
//	fmt.Println(snapshot.Text)
//
// Send named keys:
//
//	vt.SendKeys(ctx, "ls", htlib.KeySpace, "-la", htlib.KeyEnter)
//	vt.SendKeys(ctx, htlib.Ctrl('c'))
//
// # Asynchronous API
//
// Stream events in real-time:
//
//	for event := range vt.Events() {
//	    switch e := event.(type) {
//	    case htlib.InitEvent:
//	        fmt.Printf("Terminal started: PID %d\n", e.PID)
//	    case htlib.OutputEvent:
//	        fmt.Print(e.Seq)
//	    case htlib.SnapshotEvent:
//	        fmt.Println(e.Text)
//	    }
//	}
//
// Create independent subscribers:
//
//	subscriber := vt.Subscribe()
//	defer vt.Unsubscribe(subscriber)
//
//	go func() {
//	    for event := range subscriber {
//	        // Process events independently
//	    }
//	}()
//
// # Context Support
//
// All blocking operations support context cancellation:
//
//	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
//	defer cancel()
//
//	snapshot, err := vt.WaitForSnapshot(ctx)
//	if err == context.DeadlineExceeded {
//	    log.Println("Operation timed out")
//	}
//
// # Key Helpers
//
// Use provided helpers for key combinations:
//
//	htlib.Ctrl('c')              // Control+C
//	htlib.Alt('x')               // Alt+X
//	htlib.Shift(htlib.KeyLeft)   // Shift+Left
//	htlib.CtrlShift(htlib.KeyUp) // Ctrl+Shift+Up
//
// # Event Types
//
// Four event types are supported:
//
//   - InitEvent: Emitted once at startup with initial terminal state
//   - OutputEvent: Emitted when the terminal produces output
//   - ResizeEvent: Emitted when the terminal is resized
//   - SnapshotEvent: Emitted in response to TakeSnapshot command
//
// # Testing
//
// Perfect for CLI application testing:
//
//	func TestCLI(t *testing.T) {
//	    vt := htlib.New(htlib.DefaultConfig())
//	    vt.Start(context.Background())
//	    defer vt.Close()
//
//	    <-vt.Events() // Wait for init
//
//	    vt.Input(context.Background(), "my-cli --version\n")
//	    time.Sleep(100 * time.Millisecond)
//
//	    snapshot, _ := vt.WaitForSnapshot(context.Background())
//	    if !strings.Contains(snapshot.Text, "v1.0.0") {
//	        t.Error("Version not found")
//	    }
//	}
//
// For more examples, see the examples/ directory in the repository.
package htlib
