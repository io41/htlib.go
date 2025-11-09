package htlib

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	vt := New(DefaultConfig())
	if vt == nil {
		t.Fatal("expected non-nil VirtualTerminal")
	}

	if vt.config.Binary != "/bin/bash" {
		t.Errorf("expected binary /bin/bash, got %s", vt.config.Binary)
	}

	if vt.config.HtBinary != "ht" {
		t.Errorf("expected ht binary 'ht', got %s", vt.config.HtBinary)
	}
}

func TestStartAndClose(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	// Wait a moment for initialization
	time.Sleep(100 * time.Millisecond)

	// Close can return an error if the process was killed
	// This is expected behavior
	_ = vt.Close()
}

func TestDoubleStart(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Second start should fail
	if err := vt.Start(ctx); err != ErrAlreadyStarted {
		t.Errorf("expected ErrAlreadyStarted, got %v", err)
	}
}

func TestInputBeforeStart(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	err := vt.Input(ctx, "test")
	if err != ErrNotStarted {
		t.Errorf("expected ErrNotStarted, got %v", err)
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		cols int
		rows int
	}{
		{
			name: "default size",
			cfg:  Config{Size: "120x40"},
			cols: 120,
			rows: 40,
		},
		{
			name: "explicit cols/rows",
			cfg:  Config{Cols: 80, Rows: 24},
			cols: 80,
			rows: 24,
		},
		{
			name: "cols/rows override size",
			cfg:  Config{Size: "120x40", Cols: 100, Rows: 30},
			cols: 100,
			rows: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vt := New(tt.cfg)
			cols, rows := vt.Size()
			if cols != tt.cols || rows != tt.rows {
				t.Errorf("expected size %dx%d, got %dx%d", tt.cols, tt.rows, cols, rows)
			}
		})
	}
}

func TestInitEvent(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Wait for init event
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case event := <-vt.Events():
		initEvent, ok := event.(InitEvent)
		if !ok {
			t.Fatalf("expected InitEvent, got %T", event)
		}
		if initEvent.PID == 0 {
			t.Error("expected non-zero PID")
		}
		if initEvent.Cols == 0 || initEvent.Rows == 0 {
			t.Error("expected non-zero dimensions")
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for init event")
	}
}

func TestSendInput(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Wait for init
	<-vt.Events()

	// Send input
	if err := vt.Input(ctx, "echo hello\n"); err != nil {
		t.Fatalf("failed to send input: %v", err)
	}

	// Should receive output events
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	receivedOutput := false
	for {
		select {
		case event := <-vt.Events():
			if _, ok := event.(OutputEvent); ok {
				receivedOutput = true
				return
			}
		case <-ctx.Done():
			if !receivedOutput {
				t.Fatal("timeout waiting for output event")
			}
			return
		}
	}
}

func TestSubscribe(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Create subscriber
	sub := vt.Subscribe()
	defer vt.Unsubscribe(sub)

	// Wait for init event
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case event := <-sub:
		if _, ok := event.(InitEvent); !ok {
			t.Errorf("expected InitEvent, got %T", event)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for event on subscriber")
	}
}

func TestWaitForSnapshot(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Wait for init
	<-vt.Events()

	// Get snapshot
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	snapshot, err := vt.WaitForSnapshot(ctx)
	if err != nil {
		t.Fatalf("failed to get snapshot: %v", err)
	}

	if snapshot == nil {
		t.Fatal("expected non-nil snapshot")
	}

	if snapshot.Cols == 0 || snapshot.Rows == 0 {
		t.Error("expected non-zero dimensions")
	}

	if snapshot.Text == "" {
		t.Error("expected non-empty text")
	}
}

func TestResize(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Wait for init
	<-vt.Events()

	// Send resize
	if err := vt.Resize(ctx, 80, 24); err != nil {
		t.Fatalf("failed to resize: %v", err)
	}

	// Should receive resize event
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	for {
		select {
		case event := <-vt.Events():
			if resizeEvent, ok := event.(ResizeEvent); ok {
				if resizeEvent.Cols != 80 || resizeEvent.Rows != 24 {
					t.Errorf("expected size 80x24, got %dx%d", resizeEvent.Cols, resizeEvent.Rows)
				}
				return
			}
		case <-ctx.Done():
			t.Fatal("timeout waiting for resize event")
		}
	}
}

func TestSendKeys(t *testing.T) {
	vt := New(DefaultConfig())
	ctx := context.Background()

	if err := vt.Start(ctx); err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	defer vt.Close()

	// Wait for init
	<-vt.Events()

	// Send keys
	if err := vt.SendKeys(ctx, "echo", "Space", "test", "Enter"); err != nil {
		t.Fatalf("failed to send keys: %v", err)
	}

	// Should receive output
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	receivedOutput := false
	for {
		select {
		case event := <-vt.Events():
			if _, ok := event.(OutputEvent); ok {
				receivedOutput = true
				return
			}
		case <-ctx.Done():
			if !receivedOutput {
				t.Fatal("timeout waiting for output")
			}
			return
		}
	}
}
