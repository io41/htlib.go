package htlib

import "errors"

var (
	// ErrNotStarted is returned when attempting to use a VirtualTerminal that hasn't been started.
	ErrNotStarted = errors.New("virtual terminal not started")

	// ErrAlreadyStarted is returned when attempting to start a VirtualTerminal that's already running.
	ErrAlreadyStarted = errors.New("virtual terminal already started")

	// ErrClosed is returned when attempting to use a closed VirtualTerminal.
	ErrClosed = errors.New("virtual terminal closed")

	// ErrTimeout is returned when an operation times out.
	ErrTimeout = errors.New("operation timed out")

	// ErrInvalidEvent is returned when an invalid event is received.
	ErrInvalidEvent = errors.New("invalid event received")

	// ErrProcessExited is returned when the ht process exits unexpectedly.
	ErrProcessExited = errors.New("ht process exited")
)
