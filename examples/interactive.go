package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

	fmt.Println("Interactive terminal session started!")
	fmt.Println("Commands you can try:")
	fmt.Println("  - Type any bash command")
	fmt.Println("  - Use 'snapshot' to see current terminal view")
	fmt.Println("  - Use 'exit' or Ctrl+C to quit")
	fmt.Println()

	// Wait for init event
	<-vt.Events()

	// Start goroutine to print terminal output
	go func() {
		for event := range vt.Events() {
			if outputEvent, ok := event.(htlib.OutputEvent); ok {
				fmt.Print(outputEvent.Seq)
			}
		}
	}()

	// Read user input and send to terminal
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		// Handle special commands
		switch strings.TrimSpace(input) {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		case "snapshot":
			snapshot, err := vt.WaitForSnapshot(context.Background())
			if err != nil {
				log.Printf("Failed to get snapshot: %v", err)
				continue
			}
			fmt.Println("\n=== SNAPSHOT ===")
			fmt.Println(snapshot.Text)
			fmt.Println("=== END SNAPSHOT ===\n")
			continue

		case "clear":
			if err := vt.SendKeys(ctx, htlib.Ctrl('l')); err != nil {
				log.Printf("Failed to send keys: %v", err)
			}
			continue
		}

		// Send input to terminal
		if err := vt.Input(ctx, input+"\n"); err != nil {
			log.Printf("Failed to send input: %v", err)
		}

		// Give terminal time to process
		time.Sleep(100 * time.Millisecond)
	}
}
