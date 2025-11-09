package htlib

import (
	"encoding/json"
	"time"
)

// Config represents the configuration for a VirtualTerminal.
type Config struct {
	// Binary is the path to the binary to run inside the terminal (default: /bin/bash)
	Binary string
	// Args are arguments to pass to the binary
	Args []string
	// Size is the terminal size in "COLSxROWS" format (default: 120x40)
	Size string
	// Cols is the number of columns (overrides Size if set)
	Cols int
	// Rows is the number of rows (overrides Size if set)
	Rows int
	// HtBinary is the path to the ht binary (default: "ht")
	HtBinary string
	// Env is additional environment variables to pass to the process
	Env []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Binary:   "/bin/bash",
		Args:     []string{},
		Size:     "120x40",
		Cols:     0,
		Rows:     0,
		HtBinary: "ht",
		Env:      []string{},
	}
}

// EventType represents the type of event received from ht.
type EventType string

const (
	// EventTypeInit is emitted once at startup
	EventTypeInit EventType = "init"
	// EventTypeOutput is emitted when terminal produces output
	EventTypeOutput EventType = "output"
	// EventTypeResize is emitted when terminal is resized
	EventTypeResize EventType = "resize"
	// EventTypeSnapshot is emitted in response to takeSnapshot command
	EventTypeSnapshot EventType = "snapshot"
	// EventTypeMouse is emitted when mouse events occur
	EventTypeMouse EventType = "mouse"
)

// Event represents an event received from the ht process.
type Event interface {
	Type() EventType
}

// InitEvent is emitted once at startup and contains the initial terminal state.
type InitEvent struct {
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
	PID  int    `json:"pid"`
	Seq  string `json:"seq"`  // Raw VT100 output
	Text string `json:"text"` // Rendered text view
	Time time.Time
}

func (e InitEvent) Type() EventType { return EventTypeInit }

// OutputEvent is emitted when the terminal produces output.
type OutputEvent struct {
	Seq  string `json:"seq"` // Raw VT100 output
	Time time.Time
}

func (e OutputEvent) Type() EventType { return EventTypeOutput }

// ResizeEvent is emitted when the terminal is resized.
type ResizeEvent struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
	Time time.Time
}

func (e ResizeEvent) Type() EventType { return EventTypeResize }

// SnapshotEvent is emitted in response to a takeSnapshot command.
type SnapshotEvent struct {
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
	Seq  string `json:"seq"`  // Raw VT100 output
	Text string `json:"text"` // Rendered text view
	Time time.Time
}

func (e SnapshotEvent) Type() EventType { return EventTypeSnapshot }

// MouseEvent is emitted when mouse events occur in the terminal.
// Note: The application running in the terminal must enable mouse tracking
// for these events to be emitted.
type MouseEvent struct {
	Event  string `json:"event"`  // "click", "press", "release", "drag"
	Button string `json:"button"` // "left", "right", "middle", "wheel_up", "wheel_down"
	Row    int    `json:"row"`    // 1-based row coordinate
	Col    int    `json:"col"`    // 1-based column coordinate
	Shift  bool   `json:"shift"`  // Shift modifier key pressed
	Ctrl   bool   `json:"ctrl"`   // Control modifier key pressed
	Alt    bool   `json:"alt"`    // Alt modifier key pressed
	Time   time.Time
}

func (e MouseEvent) Type() EventType { return EventTypeMouse }

// MouseModifiers represents modifier keys for mouse events.
type MouseModifiers struct {
	Shift bool
	Ctrl  bool
	Alt   bool
}

// rawEvent is used for unmarshaling JSON events from ht.
type rawEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// command represents a command to send to ht via STDIN.
type command struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
	Keys    []string    `json:"keys,omitempty"`
	Cols    int         `json:"cols,omitempty"`
	Rows    int         `json:"rows,omitempty"`
	Event   string      `json:"event,omitempty"`
	Button  string      `json:"button,omitempty"`
	Row     int         `json:"row,omitempty"`
	Col     int         `json:"col,omitempty"`
	Shift   bool        `json:"shift,omitempty"`
	Ctrl    bool        `json:"ctrl,omitempty"`
	Alt     bool        `json:"alt,omitempty"`
}
