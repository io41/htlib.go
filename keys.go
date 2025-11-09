package htlib

// Common key constants for use with SendKeys.
// These match the key names supported by ht.
const (
	// Special keys
	KeyEnter     = "Enter"
	KeySpace     = "Space"
	KeyEscape    = "Escape"
	KeyTab       = "Tab"
	KeyBackspace = "Backspace"

	// Arrow keys
	KeyLeft  = "Left"
	KeyRight = "Right"
	KeyUp    = "Up"
	KeyDown  = "Down"

	// Navigation keys
	KeyHome     = "Home"
	KeyEnd      = "End"
	KeyPageUp   = "PageUp"
	KeyPageDown = "PageDown"

	// Function keys
	KeyF1  = "F1"
	KeyF2  = "F2"
	KeyF3  = "F3"
	KeyF4  = "F4"
	KeyF5  = "F5"
	KeyF6  = "F6"
	KeyF7  = "F7"
	KeyF8  = "F8"
	KeyF9  = "F9"
	KeyF10 = "F10"
	KeyF11 = "F11"
	KeyF12 = "F12"
)

// Ctrl returns a control key sequence for the given character.
// Example: Ctrl('c') returns "C-c"
func Ctrl(c rune) string {
	return "C-" + string(c)
}

// Alt returns an alt key sequence for the given character.
// Example: Alt('x') returns "A-x"
func Alt(c rune) string {
	return "A-" + string(c)
}

// Shift returns a shift key sequence for the given key.
// Example: Shift("Left") returns "S-Left"
func Shift(key string) string {
	return "S-" + key
}

// CtrlShift returns a control+shift key sequence for the given key.
// Example: CtrlShift("Left") returns "C-S-Left"
func CtrlShift(key string) string {
	return "C-S-" + key
}

// AltShift returns an alt+shift key sequence for the given key.
// Example: AltShift("Left") returns "A-S-Left"
func AltShift(key string) string {
	return "A-S-" + key
}

// CtrlAlt returns a control+alt key sequence for the given key.
// Example: CtrlAlt("Left") returns "C-A-Left"
func CtrlAlt(key string) string {
	return "C-A-" + key
}
