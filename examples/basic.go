package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/andyk/htlib.go"
)

func main() {
	// Create a new virtual terminal with default configuration
	vt := htlib.New(htlib.DefaultConfig())

	// Start the terminal
	ctx := context.Background()
	if err := vt.Start(ctx); err != nil {
		log.Fatalf("Failed to start terminal: %v", err)
	}
	defer vt.Close()

	fmt.Println("Terminal started successfully!")

	// Wait for the init event
	event := <-vt.Events()
	if initEvent, ok := event.(htlib.InitEvent); ok {
		fmt.Printf("Terminal initialized: %dx%d (PID: %d)\n",
			initEvent.Cols, initEvent.Rows, initEvent.PID)
	}

	// Send a simple command
	fmt.Println("\nSending command: echo 'Hello from htlib.go!'")
	if err := vt.Input(ctx, "echo 'Hello from htlib.go!'\n"); err != nil {
		log.Fatalf("Failed to send input: %v", err)
	}

	// Read output for 2 seconds
	fmt.Println("\nTerminal output:")
	timeout := time.After(2 * time.Second)
	for {
		select {
		case event := <-vt.Events():
			if outputEvent, ok := event.(htlib.OutputEvent); ok {
				fmt.Print(outputEvent.Seq)
			}
		case <-timeout:
			fmt.Println("\n\nTimeout reached, taking snapshot...")
			goto snapshot
		}
	}

snapshot:
	// Get a snapshot of the terminal
	snapshot, err := vt.WaitForSnapshot(context.Background())
	if err != nil {
		log.Fatalf("Failed to get snapshot: %v", err)
	}

	fmt.Println("Terminal snapshot (text view):")
	fmt.Println("---")
	fmt.Println(snapshot.Text)
	fmt.Println("---")
}
