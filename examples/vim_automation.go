package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andyk/htlib.go"
)

func main() {
	// Create a virtual terminal
	config := htlib.DefaultConfig()
	config.Cols = 80
	config.Rows = 24
	vt := htlib.New(config)

	ctx := context.Background()
	if err := vt.Start(ctx); err != nil {
		log.Fatalf("Failed to start terminal: %v", err)
	}
	defer vt.Close()

	// Wait for init
	<-vt.Events()

	fmt.Println("Automating vim to create and edit a file...")

	// Start vim with a new file
	fmt.Println("1. Opening nano (easier to automate than vim)...")
	if err := vt.Input(ctx, "nano test.txt\n"); err != nil {
		log.Fatalf("Failed to start nano: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Enter insert mode and type some text
	fmt.Println("2. Writing content...")
	lines := []string{
		"Hello from htlib.go!",
		"",
		"This file was created automatically using the htlib.go library.",
		"It demonstrates how to automate terminal-based editors.",
		"",
		"Features:",
		"- Programmatic terminal control",
		"- Send keys and input",
		"- Capture terminal state",
	}

	for _, line := range lines {
		if err := vt.Input(ctx, line+"\n"); err != nil {
			log.Fatalf("Failed to write line: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Get a snapshot to see what we've typed
	fmt.Println("3. Taking snapshot of editor...")
	snapshot, err := vt.WaitForSnapshot(ctx)
	if err != nil {
		log.Fatalf("Failed to get snapshot: %v", err)
	}

	fmt.Println("\n=== Editor Content ===")
	fmt.Println(snapshot.Text)
	fmt.Println("=== End ===\n")

	// Save and exit (Ctrl+X, Y, Enter in nano)
	fmt.Println("4. Saving and exiting...")
	if err := vt.SendKeys(ctx, htlib.Ctrl('x')); err != nil {
		log.Fatalf("Failed to send Ctrl+X: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Confirm save
	if err := vt.SendKeys(ctx, "y"); err != nil {
		log.Fatalf("Failed to confirm save: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Confirm filename
	if err := vt.SendKeys(ctx, htlib.KeyEnter); err != nil {
		log.Fatalf("Failed to confirm filename: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// Verify the file was created
	fmt.Println("5. Verifying file was created...")
	if err := vt.Input(ctx, "cat test.txt\n"); err != nil {
		log.Fatalf("Failed to cat file: %v", err)
	}
	time.Sleep(500 * time.Millisecond)

	// Get final snapshot
	finalSnapshot, err := vt.WaitForSnapshot(ctx)
	if err != nil {
		log.Fatalf("Failed to get final snapshot: %v", err)
	}

	fmt.Println("\n=== Final Terminal State ===")
	fmt.Println(finalSnapshot.Text)
	fmt.Println("=== End ===")

	// Clean up
	fmt.Println("\n6. Cleaning up...")
	if err := vt.Input(ctx, "rm test.txt\n"); err != nil {
		log.Fatalf("Failed to remove file: %v", err)
	}

	fmt.Println("\nAutomation complete!")
}
