package htlib

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

// VirtualTerminal represents a headless terminal session managed by ht.
type VirtualTerminal struct {
	config Config
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	// Event handling
	events      chan Event
	subscribers []chan Event
	mu          sync.RWMutex
	started     bool
	closed      bool

	// Background goroutine management
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Error handling
	err error
}

// New creates a new VirtualTerminal with the given configuration.
func New(config Config) *VirtualTerminal {
	if config.Binary == "" {
		config.Binary = "/bin/bash"
	}
	if config.HtBinary == "" {
		config.HtBinary = "ht"
	}
	if config.Size == "" && config.Cols == 0 && config.Rows == 0 {
		config.Size = "120x40"
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &VirtualTerminal{
		config:      config,
		events:      make(chan Event, 100),
		subscribers: make([]chan Event, 0),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start launches the ht subprocess and begins processing events.
func (vt *VirtualTerminal) Start(ctx context.Context) error {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	if vt.started {
		return ErrAlreadyStarted
	}
	if vt.closed {
		return ErrClosed
	}

	// Build command arguments
	args := vt.buildArgs()

	// Create command
	vt.cmd = exec.CommandContext(vt.ctx, vt.config.HtBinary, args...)
	if len(vt.config.Env) > 0 {
		vt.cmd.Env = append(vt.cmd.Env, vt.config.Env...)
	}

	// Setup pipes
	var err error
	vt.stdin, err = vt.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	vt.stdout, err = vt.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	vt.stderr, err = vt.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := vt.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ht process: %w", err)
	}

	vt.started = true

	// Start background goroutines
	vt.wg.Add(2)
	go vt.readEvents()
	go vt.waitForExit()

	return nil
}

// buildArgs constructs the command line arguments for ht.
func (vt *VirtualTerminal) buildArgs() []string {
	args := []string{}

	// Add size flag
	var size string
	if vt.config.Cols > 0 && vt.config.Rows > 0 {
		size = fmt.Sprintf("%dx%d", vt.config.Cols, vt.config.Rows)
	} else {
		size = vt.config.Size
	}
	args = append(args, "--size", size)

	// Add subscription to all events
	args = append(args, "--subscribe", "init,output,resize,snapshot")

	// Add binary and its arguments
	args = append(args, vt.config.Binary)
	args = append(args, vt.config.Args...)

	return args
}

// readEvents reads events from stdout and dispatches them.
func (vt *VirtualTerminal) readEvents() {
	defer vt.wg.Done()
	defer close(vt.events)

	scanner := bufio.NewScanner(vt.stdout)
	for scanner.Scan() {
		line := scanner.Text()
		event, err := vt.parseEvent(line)
		if err != nil {
			// Log error but continue
			continue
		}

		// Send to main events channel
		select {
		case vt.events <- event:
		case <-vt.ctx.Done():
			return
		}

		// Send to subscribers
		vt.mu.RLock()
		for _, sub := range vt.subscribers {
			select {
			case sub <- event:
			default:
				// Skip if subscriber is not ready
			}
		}
		vt.mu.RUnlock()
	}

	if err := scanner.Err(); err != nil {
		vt.mu.Lock()
		vt.err = fmt.Errorf("error reading stdout: %w", err)
		vt.mu.Unlock()
	}
}

// waitForExit waits for the ht process to exit.
func (vt *VirtualTerminal) waitForExit() {
	defer vt.wg.Done()

	err := vt.cmd.Wait()
	vt.mu.Lock()
	if err != nil && vt.err == nil {
		vt.err = fmt.Errorf("ht process exited: %w", err)
	}
	vt.mu.Unlock()

	// Cancel context to stop all operations
	vt.cancel()
}

// parseEvent parses a JSON event line from ht.
func (vt *VirtualTerminal) parseEvent(line string) (Event, error) {
	var raw rawEvent
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse event: %w", err)
	}

	now := time.Now()

	switch raw.Type {
	case "init":
		var data struct {
			Cols int    `json:"cols"`
			Rows int    `json:"rows"`
			PID  int    `json:"pid"`
			Seq  string `json:"seq"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, err
		}
		return InitEvent{
			Cols: data.Cols,
			Rows: data.Rows,
			PID:  data.PID,
			Seq:  data.Seq,
			Text: data.Text,
			Time: now,
		}, nil

	case "output":
		var data struct {
			Seq string `json:"seq"`
		}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, err
		}
		return OutputEvent{
			Seq:  data.Seq,
			Time: now,
		}, nil

	case "resize":
		var data struct {
			Cols int `json:"cols"`
			Rows int `json:"rows"`
		}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, err
		}
		return ResizeEvent{
			Cols: data.Cols,
			Rows: data.Rows,
			Time: now,
		}, nil

	case "snapshot":
		var data struct {
			Cols int    `json:"cols"`
			Rows int    `json:"rows"`
			Seq  string `json:"seq"`
			Text string `json:"text"`
		}
		if err := json.Unmarshal(raw.Data, &data); err != nil {
			return nil, err
		}
		return SnapshotEvent{
			Cols: data.Cols,
			Rows: data.Rows,
			Seq:  data.Seq,
			Text: data.Text,
			Time: now,
		}, nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidEvent, raw.Type)
	}
}

// sendCommand sends a JSON command to ht via stdin.
func (vt *VirtualTerminal) sendCommand(cmd command) error {
	vt.mu.RLock()
	defer vt.mu.RUnlock()

	if !vt.started {
		return ErrNotStarted
	}
	if vt.closed {
		return ErrClosed
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	data = append(data, '\n')
	if _, err := vt.stdin.Write(data); err != nil {
		return fmt.Errorf("failed to write command: %w", err)
	}

	return nil
}

// Input sends raw input to the terminal.
func (vt *VirtualTerminal) Input(ctx context.Context, text string) error {
	cmd := command{
		Type:    "input",
		Payload: text,
	}
	return vt.sendCommand(cmd)
}

// SendKeys sends named keys to the terminal.
// Examples: "Enter", "C-c", "Left", "F1", etc.
func (vt *VirtualTerminal) SendKeys(ctx context.Context, keys ...string) error {
	cmd := command{
		Type: "sendKeys",
		Keys: keys,
	}
	return vt.sendCommand(cmd)
}

// Resize resizes the terminal to the specified dimensions.
func (vt *VirtualTerminal) Resize(ctx context.Context, cols, rows int) error {
	cmd := command{
		Type: "resize",
		Cols: cols,
		Rows: rows,
	}
	return vt.sendCommand(cmd)
}

// TakeSnapshot requests a snapshot of the terminal state.
// Use WaitForSnapshot to receive the snapshot event.
func (vt *VirtualTerminal) TakeSnapshot(ctx context.Context) error {
	cmd := command{
		Type: "takeSnapshot",
	}
	return vt.sendCommand(cmd)
}

// WaitForSnapshot requests a snapshot and waits for the response.
// This is a convenience method that combines TakeSnapshot with event waiting.
func (vt *VirtualTerminal) WaitForSnapshot(ctx context.Context) (*SnapshotEvent, error) {
	// Subscribe to events temporarily
	eventChan := vt.Subscribe()
	defer vt.Unsubscribe(eventChan)

	// Request snapshot
	if err := vt.TakeSnapshot(ctx); err != nil {
		return nil, err
	}

	// Wait for snapshot event
	for {
		select {
		case event := <-eventChan:
			if snapshot, ok := event.(SnapshotEvent); ok {
				return &snapshot, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-vt.ctx.Done():
			return nil, ErrClosed
		}
	}
}

// Events returns a channel that receives all events from the terminal.
// This channel is closed when the terminal is closed.
func (vt *VirtualTerminal) Events() <-chan Event {
	return vt.events
}

// Subscribe creates a new subscriber channel for receiving events.
// The caller is responsible for reading from this channel to avoid blocking.
// Call Unsubscribe when done.
func (vt *VirtualTerminal) Subscribe() chan Event {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	ch := make(chan Event, 100)
	vt.subscribers = append(vt.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber channel.
func (vt *VirtualTerminal) Unsubscribe(ch chan Event) {
	vt.mu.Lock()
	defer vt.mu.Unlock()

	for i, sub := range vt.subscribers {
		if sub == ch {
			// Remove from slice
			vt.subscribers = append(vt.subscribers[:i], vt.subscribers[i+1:]...)
			close(ch)
			return
		}
	}
}

// Close terminates the ht process and cleans up resources.
func (vt *VirtualTerminal) Close() error {
	vt.mu.Lock()
	if vt.closed {
		vt.mu.Unlock()
		return nil
	}
	vt.closed = true
	vt.mu.Unlock()

	// Cancel context to stop background goroutines
	vt.cancel()

	// Close stdin to signal ht to exit
	if vt.stdin != nil {
		vt.stdin.Close()
	}

	// Wait for background goroutines
	vt.wg.Wait()

	// Close all subscriber channels
	vt.mu.Lock()
	for _, sub := range vt.subscribers {
		close(sub)
	}
	vt.subscribers = nil
	vt.mu.Unlock()

	return vt.err
}

// Err returns any error that occurred during operation.
func (vt *VirtualTerminal) Err() error {
	vt.mu.RLock()
	defer vt.mu.RUnlock()
	return vt.err
}

// Size returns the current terminal size.
func (vt *VirtualTerminal) Size() (cols, rows int) {
	if vt.config.Cols > 0 && vt.config.Rows > 0 {
		return vt.config.Cols, vt.config.Rows
	}

	// Parse from size string
	parts := strings.Split(vt.config.Size, "x")
	if len(parts) == 2 {
		cols, _ = strconv.Atoi(parts[0])
		rows, _ = strconv.Atoi(parts[1])
	}
	return
}
