package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/io41/htlib.go"
)

// This example demonstrates mouse interaction with a TUI application.
// It shows how to send mouse clicks, drags, scrolls, and listen for mouse events.
//
// Note: The application running in the terminal must enable mouse tracking
// for mouse events to have any effect. Many modern TUI applications like
// vim, tmux, htop, etc. support mouse tracking.
func main() {
	// Create a terminal running a shell
	config := htlib.DefaultConfig()
	vt := htlib.New(config)

	ctx := context.Background()
	if err := vt.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer vt.Close()

	// Wait for initialization
	<-vt.Events()

	fmt.Println("Mouse Demo - Demonstrating mouse support in htlib.go")
	fmt.Println("========================================================")
	fmt.Println()

	// Example 1: Basic Mouse Click
	fmt.Println("Example 1: Sending left click at position (5, 10)")
	if err := vt.MouseClick(ctx, "left", 5, 10); err != nil {
		log.Printf("Error sending mouse click: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Example 2: Right Click
	fmt.Println("Example 2: Sending right click at position (10, 20)")
	if err := vt.MouseClick(ctx, "right", 10, 20); err != nil {
		log.Printf("Error sending right click: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Example 3: Mouse Drag (Press, Drag, Release)
	fmt.Println("Example 3: Performing drag operation")
	fmt.Println("  - Press at (5, 10)")
	if err := vt.MousePress(ctx, "left", 5, 10); err != nil {
		log.Printf("Error sending mouse press: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("  - Drag to (5, 20)")
	if err := vt.MouseDrag(ctx, "left", 5, 20); err != nil {
		log.Printf("Error sending mouse drag: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("  - Release at (5, 20)")
	if err := vt.MouseRelease(ctx, "left", 5, 20); err != nil {
		log.Printf("Error sending mouse release: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Example 4: Mouse Scroll
	fmt.Println("Example 4: Scrolling with mouse wheel")
	fmt.Println("  - Scroll down")
	if err := vt.MouseScroll(ctx, "wheel_down", 10, 10); err != nil {
		log.Printf("Error sending scroll: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	fmt.Println("  - Scroll up")
	if err := vt.MouseScroll(ctx, "wheel_up", 10, 10); err != nil {
		log.Printf("Error sending scroll: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Example 5: Mouse Click with Modifiers
	fmt.Println("Example 5: Click with modifier keys (Ctrl+Click)")
	if err := vt.MouseClickWithModifiers(ctx, "left", 8, 15, htlib.MouseModifiers{
		Ctrl:  true,
		Shift: false,
		Alt:   false,
	}); err != nil {
		log.Printf("Error sending mouse click with modifiers: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	fmt.Println()
	fmt.Println("Example 6: Listening for mouse events")
	fmt.Println("Subscribing to terminal events (will monitor for 3 seconds)...")

	// Subscribe to events to capture any mouse events from the terminal
	subscriber := vt.Subscribe()
	defer vt.Unsubscribe(subscriber)

	// Start a goroutine to monitor events
	done := make(chan bool)
	go func() {
		timeout := time.After(3 * time.Second)
		for {
			select {
			case event := <-subscriber:
				if mouseEvent, ok := event.(htlib.MouseEvent); ok {
					fmt.Printf("  Received mouse event: %s %s at (%d, %d)",
						mouseEvent.Event, mouseEvent.Button,
						mouseEvent.Row, mouseEvent.Col)
					if mouseEvent.Shift {
						fmt.Print(" [Shift]")
					}
					if mouseEvent.Ctrl {
						fmt.Print(" [Ctrl]")
					}
					if mouseEvent.Alt {
						fmt.Print(" [Alt]")
					}
					fmt.Println()
				}
			case <-timeout:
				done <- true
				return
			}
		}
	}()

	// Generate some mouse events for demonstration
	time.Sleep(500 * time.Millisecond)
	vt.MouseClick(ctx, "left", 3, 5)
	time.Sleep(500 * time.Millisecond)
	vt.MouseClickWithModifiers(ctx, "right", 7, 12, htlib.MouseModifiers{Shift: true})

	// Wait for monitoring to complete
	<-done

	fmt.Println()
	fmt.Println("Mouse demo completed!")
	fmt.Println()
	fmt.Println("Note: To see actual mouse interaction, run a TUI application")
	fmt.Println("that supports mouse tracking, such as:")
	fmt.Println("  - vim (with :set mouse=a)")
	fmt.Println("  - tmux")
	fmt.Println("  - htop")
	fmt.Println("  - nano")
}
