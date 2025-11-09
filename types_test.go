package htlib

import (
	"encoding/json"
	"testing"
)

func TestEventTypeConstants(t *testing.T) {
	if EventTypeInit != "init" {
		t.Errorf("expected 'init', got '%s'", EventTypeInit)
	}
	if EventTypeOutput != "output" {
		t.Errorf("expected 'output', got '%s'", EventTypeOutput)
	}
	if EventTypeResize != "resize" {
		t.Errorf("expected 'resize', got '%s'", EventTypeResize)
	}
	if EventTypeSnapshot != "snapshot" {
		t.Errorf("expected 'snapshot', got '%s'", EventTypeSnapshot)
	}
	if EventTypeMouse != "mouse" {
		t.Errorf("expected 'mouse', got '%s'", EventTypeMouse)
	}
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		event     Event
		eventType EventType
	}{
		{"init", InitEvent{}, EventTypeInit},
		{"output", OutputEvent{}, EventTypeOutput},
		{"resize", ResizeEvent{}, EventTypeResize},
		{"snapshot", SnapshotEvent{}, EventTypeSnapshot},
		{"mouse", MouseEvent{}, EventTypeMouse},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.event.Type() != tt.eventType {
				t.Errorf("expected type %s, got %s", tt.eventType, tt.event.Type())
			}
		})
	}
}

func TestCommandMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		cmd      command
		expected string
	}{
		{
			name:     "input command",
			cmd:      command{Type: "input", Payload: "test"},
			expected: `{"type":"input","payload":"test"}`,
		},
		{
			name:     "sendKeys command",
			cmd:      command{Type: "sendKeys", Keys: []string{"Enter", "C-c"}},
			expected: `{"type":"sendKeys","keys":["Enter","C-c"]}`,
		},
		{
			name:     "resize command",
			cmd:      command{Type: "resize", Cols: 80, Rows: 24},
			expected: `{"type":"resize","cols":80,"rows":24}`,
		},
		{
			name:     "takeSnapshot command",
			cmd:      command{Type: "takeSnapshot"},
			expected: `{"type":"takeSnapshot"}`,
		},
		{
			name: "mouse click command",
			cmd: command{
				Type:   "mouse",
				Event:  "click",
				Button: "left",
				Row:    10,
				Col:    20,
			},
			expected: `{"type":"mouse","event":"click","button":"left","row":10,"col":20}`,
		},
		{
			name: "mouse click with modifiers command",
			cmd: command{
				Type:   "mouse",
				Event:  "click",
				Button: "left",
				Row:    5,
				Col:    15,
				Shift:  true,
				Ctrl:   true,
			},
			expected: `{"type":"mouse","event":"click","button":"left","row":5,"col":15,"shift":true,"ctrl":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.cmd)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			// Check that required fields are present
			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}

			if result["type"] != tt.cmd.Type {
				t.Errorf("expected type %s, got %v", tt.cmd.Type, result["type"])
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Binary != "/bin/bash" {
		t.Errorf("expected binary /bin/bash, got %s", cfg.Binary)
	}

	if cfg.Size != "120x40" {
		t.Errorf("expected size 120x40, got %s", cfg.Size)
	}

	if cfg.HtBinary != "ht" {
		t.Errorf("expected ht binary 'ht', got %s", cfg.HtBinary)
	}

	if len(cfg.Args) != 0 {
		t.Errorf("expected empty args, got %v", cfg.Args)
	}
}

func TestParseEvent(t *testing.T) {
	vt := New(DefaultConfig())

	tests := []struct {
		name      string
		json      string
		eventType EventType
		checkFunc func(*testing.T, Event)
	}{
		{
			name:      "init event",
			json:      `{"type":"init","data":{"cols":120,"rows":40,"pid":12345,"seq":"test","text":"hello"}}`,
			eventType: EventTypeInit,
			checkFunc: func(t *testing.T, e Event) {
				init := e.(InitEvent)
				if init.Cols != 120 {
					t.Errorf("expected cols 120, got %d", init.Cols)
				}
				if init.Rows != 40 {
					t.Errorf("expected rows 40, got %d", init.Rows)
				}
				if init.PID != 12345 {
					t.Errorf("expected pid 12345, got %d", init.PID)
				}
			},
		},
		{
			name:      "output event",
			json:      `{"type":"output","data":{"seq":"test output"}}`,
			eventType: EventTypeOutput,
			checkFunc: func(t *testing.T, e Event) {
				output := e.(OutputEvent)
				if output.Seq != "test output" {
					t.Errorf("expected seq 'test output', got %s", output.Seq)
				}
			},
		},
		{
			name:      "resize event",
			json:      `{"type":"resize","data":{"cols":80,"rows":24}}`,
			eventType: EventTypeResize,
			checkFunc: func(t *testing.T, e Event) {
				resize := e.(ResizeEvent)
				if resize.Cols != 80 {
					t.Errorf("expected cols 80, got %d", resize.Cols)
				}
				if resize.Rows != 24 {
					t.Errorf("expected rows 24, got %d", resize.Rows)
				}
			},
		},
		{
			name:      "snapshot event",
			json:      `{"type":"snapshot","data":{"cols":120,"rows":40,"seq":"test","text":"snapshot text"}}`,
			eventType: EventTypeSnapshot,
			checkFunc: func(t *testing.T, e Event) {
				snapshot := e.(SnapshotEvent)
				if snapshot.Cols != 120 {
					t.Errorf("expected cols 120, got %d", snapshot.Cols)
				}
				if snapshot.Text != "snapshot text" {
					t.Errorf("expected text 'snapshot text', got %s", snapshot.Text)
				}
			},
		},
		{
			name:      "mouse event",
			json:      `{"type":"mouse","data":{"event":"click","button":"left","row":10,"col":20,"shift":false,"ctrl":true,"alt":false}}`,
			eventType: EventTypeMouse,
			checkFunc: func(t *testing.T, e Event) {
				mouse := e.(MouseEvent)
				if mouse.Event != "click" {
					t.Errorf("expected event 'click', got %s", mouse.Event)
				}
				if mouse.Button != "left" {
					t.Errorf("expected button 'left', got %s", mouse.Button)
				}
				if mouse.Row != 10 {
					t.Errorf("expected row 10, got %d", mouse.Row)
				}
				if mouse.Col != 20 {
					t.Errorf("expected col 20, got %d", mouse.Col)
				}
				if mouse.Ctrl != true {
					t.Errorf("expected ctrl true, got %v", mouse.Ctrl)
				}
				if mouse.Shift != false {
					t.Errorf("expected shift false, got %v", mouse.Shift)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := vt.parseEvent(tt.json)
			if err != nil {
				t.Fatalf("failed to parse event: %v", err)
			}

			if event.Type() != tt.eventType {
				t.Errorf("expected type %s, got %s", tt.eventType, event.Type())
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, event)
			}
		})
	}
}
