package htlib

import "testing"

func TestKeyConstants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{KeyEnter, "Enter"},
		{KeySpace, "Space"},
		{KeyEscape, "Escape"},
		{KeyTab, "Tab"},
		{KeyBackspace, "Backspace"},
		{KeyLeft, "Left"},
		{KeyRight, "Right"},
		{KeyUp, "Up"},
		{KeyDown, "Down"},
		{KeyHome, "Home"},
		{KeyEnd, "End"},
		{KeyPageUp, "PageUp"},
		{KeyPageDown, "PageDown"},
		{KeyF1, "F1"},
		{KeyF12, "F12"},
	}

	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.constant)
		}
	}
}

func TestCtrl(t *testing.T) {
	tests := []struct {
		input    rune
		expected string
	}{
		{'c', "C-c"},
		{'a', "C-a"},
		{'z', "C-z"},
		{'x', "C-x"},
	}

	for _, tt := range tests {
		result := Ctrl(tt.input)
		if result != tt.expected {
			t.Errorf("Ctrl(%c): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestAlt(t *testing.T) {
	tests := []struct {
		input    rune
		expected string
	}{
		{'x', "A-x"},
		{'a', "A-a"},
		{'1', "A-1"},
	}

	for _, tt := range tests {
		result := Alt(tt.input)
		if result != tt.expected {
			t.Errorf("Alt(%c): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestShift(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Left", "S-Left"},
		{"Right", "S-Right"},
		{"F1", "S-F1"},
	}

	for _, tt := range tests {
		result := Shift(tt.input)
		if result != tt.expected {
			t.Errorf("Shift(%s): expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}

func TestCtrlShift(t *testing.T) {
	result := CtrlShift("Left")
	expected := "C-S-Left"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAltShift(t *testing.T) {
	result := AltShift("Right")
	expected := "A-S-Right"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestCtrlAlt(t *testing.T) {
	result := CtrlAlt("Up")
	expected := "C-A-Up"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
