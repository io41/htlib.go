package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/io41/htlib.go"
)

func main() {
	// Create a virtual terminal
	vt := htlib.New(htlib.DefaultConfig())

	ctx := context.Background()
	if err := vt.Start(ctx); err != nil {
		log.Fatalf("Failed to start terminal: %v", err)
	}
	defer vt.Close()

	fmt.Println("Event streaming example - monitoring all terminal events")
	fmt.Println()

	// Create multiple subscribers to demonstrate async event streaming
	subscriber1 := vt.Subscribe()
	defer vt.Unsubscribe(subscriber1)

	subscriber2 := vt.Subscribe()
	defer vt.Unsubscribe(subscriber2)

	// Subscriber 1: Log all events with timestamps
	go func() {
		fmt.Println("[Subscriber 1] Started - logging all events")
		for event := range subscriber1 {
			switch e := event.(type) {
			case htlib.InitEvent:
				fmt.Printf("[Sub1 %s] INIT: PID=%d, Size=%dx%d\n",
					e.Time.Format("15:04:05.000"), e.PID, e.Cols, e.Rows)
			case htlib.OutputEvent:
				fmt.Printf("[Sub1 %s] OUTPUT: %q\n",
					e.Time.Format("15:04:05.000"), truncate(e.Seq, 50))
			case htlib.ResizeEvent:
				fmt.Printf("[Sub1 %s] RESIZE: %dx%d\n",
					e.Time.Format("15:04:05.000"), e.Cols, e.Rows)
			case htlib.SnapshotEvent:
				fmt.Printf("[Sub1 %s] SNAPSHOT: %dx%d\n",
					e.Time.Format("15:04:05.000"), e.Cols, e.Rows)
			}
		}
	}()

	// Subscriber 2: Count events by type
	eventCounts := make(map[htlib.EventType]int)
	go func() {
		fmt.Println("[Subscriber 2] Started - counting events by type")
		for event := range subscriber2 {
			eventCounts[event.Type()]++
		}
		fmt.Println("\n[Subscriber 2] Final event counts:")
		for eventType, count := range eventCounts {
			fmt.Printf("  %s: %d\n", eventType, count)
		}
	}()

	// Main event stream consumer
	go func() {
		for event := range vt.Events() {
			if outputEvent, ok := event.(htlib.OutputEvent); ok {
				// Print actual output (without prefix)
				fmt.Print(outputEvent.Seq)
			}
		}
	}()

	// Wait for init
	time.Sleep(100 * time.Millisecond)

	// Perform various operations to generate events
	fmt.Println("\n--- Performing operations to generate events ---\n")

	// 1. Send some commands
	fmt.Println("1. Sending command: ls -la")
	if err := vt.Input(ctx, "ls -la\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// 2. Take a snapshot
	fmt.Println("\n2. Taking snapshot...")
	if _, err := vt.WaitForSnapshot(ctx); err != nil {
		log.Printf("Error: %v", err)
	}

	// 3. Resize terminal
	fmt.Println("\n3. Resizing terminal to 80x24...")
	if err := vt.Resize(ctx, 80, 24); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// 4. Send another command
	fmt.Println("\n4. Sending command: echo 'Testing event streaming'")
	if err := vt.Input(ctx, "echo 'Testing event streaming'\n"); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// 5. Take another snapshot
	fmt.Println("\n5. Taking another snapshot...")
	if _, err := vt.WaitForSnapshot(ctx); err != nil {
		log.Printf("Error: %v", err)
	}

	// 6. Resize back
	fmt.Println("\n6. Resizing back to 120x40...")
	if err := vt.Resize(ctx, 120, 40); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n--- Operations complete ---")
	fmt.Println("\nWaiting for subscribers to finish...")
	time.Sleep(500 * time.Millisecond)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
