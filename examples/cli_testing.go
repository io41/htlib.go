package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andyk/htlib.go"
)

// CLITester provides utilities for testing CLI applications
type CLITester struct {
	vt *htlib.VirtualTerminal
}

// NewCLITester creates a new CLI testing helper
func NewCLITester() *CLITester {
	return &CLITester{
		vt: htlib.New(htlib.DefaultConfig()),
	}
}

// Start initializes the terminal
func (t *CLITester) Start(ctx context.Context) error {
	return t.vt.Start(ctx)
}

// Close cleans up the terminal
func (t *CLITester) Close() error {
	return t.vt.Close()
}

// RunCommand runs a command and returns the output
func (t *CLITester) RunCommand(ctx context.Context, cmd string) (string, error) {
	// Send command
	if err := t.vt.Input(ctx, cmd+"\n"); err != nil {
		return "", err
	}

	// Wait for output
	time.Sleep(300 * time.Millisecond)

	// Get snapshot
	snapshot, err := t.vt.WaitForSnapshot(ctx)
	if err != nil {
		return "", err
	}

	return snapshot.Text, nil
}

// ExpectOutput checks if the output contains the expected text
func (t *CLITester) ExpectOutput(output, expected string) bool {
	return strings.Contains(output, expected)
}

// ExpectPrompt waits for a prompt and verifies it appears
func (t *CLITester) ExpectPrompt(ctx context.Context, prompt string) (bool, error) {
	snapshot, err := t.vt.WaitForSnapshot(ctx)
	if err != nil {
		return false, err
	}

	return strings.Contains(snapshot.Text, prompt), nil
}

func main() {
	fmt.Println("=== CLI Testing Framework Demo ===\n")

	// Create tester
	tester := NewCLITester()
	ctx := context.Background()

	if err := tester.Start(ctx); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
	defer tester.Close()

	// Wait for init
	time.Sleep(200 * time.Millisecond)

	// Test 1: Check basic command execution
	fmt.Println("Test 1: Running 'echo hello'")
	output, err := tester.RunCommand(ctx, "echo hello")
	if err != nil {
		log.Fatalf("Test 1 failed: %v", err)
	}
	if tester.ExpectOutput(output, "hello") {
		fmt.Println("✓ Test 1 passed: Found 'hello' in output")
	} else {
		fmt.Println("✗ Test 1 failed: Expected 'hello' in output")
	}

	// Test 2: Check command with arguments
	fmt.Println("\nTest 2: Running 'ls -la'")
	output, err = tester.RunCommand(ctx, "ls -la")
	if err != nil {
		log.Fatalf("Test 2 failed: %v", err)
	}
	if tester.ExpectOutput(output, "total") {
		fmt.Println("✓ Test 2 passed: ls -la executed successfully")
	} else {
		fmt.Println("✗ Test 2 failed: ls -la output not as expected")
	}

	// Test 3: Check environment variable
	fmt.Println("\nTest 3: Checking $HOME environment variable")
	output, err = tester.RunCommand(ctx, "echo $HOME")
	if err != nil {
		log.Fatalf("Test 3 failed: %v", err)
	}
	if tester.ExpectOutput(output, "/") {
		fmt.Println("✓ Test 3 passed: $HOME is set")
	} else {
		fmt.Println("✗ Test 3 failed: $HOME not found in output")
	}

	// Test 4: Check command that should fail
	fmt.Println("\nTest 4: Running non-existent command")
	output, err = tester.RunCommand(ctx, "nonexistentcommand123")
	if err != nil {
		log.Fatalf("Test 4 failed: %v", err)
	}
	if tester.ExpectOutput(output, "not found") || tester.ExpectOutput(output, "command not found") {
		fmt.Println("✓ Test 4 passed: Error message detected for invalid command")
	} else {
		fmt.Println("✗ Test 4 failed: Expected error message not found")
	}

	// Test 5: Multi-line command
	fmt.Println("\nTest 5: Testing multi-line output")
	output, err = tester.RunCommand(ctx, "for i in 1 2 3; do echo \"Line $i\"; done")
	if err != nil {
		log.Fatalf("Test 5 failed: %v", err)
	}
	foundAll := tester.ExpectOutput(output, "Line 1") &&
		tester.ExpectOutput(output, "Line 2") &&
		tester.ExpectOutput(output, "Line 3")
	if foundAll {
		fmt.Println("✓ Test 5 passed: All lines found in output")
	} else {
		fmt.Println("✗ Test 5 failed: Not all lines found in output")
	}

	// Test 6: Working directory
	fmt.Println("\nTest 6: Checking working directory")
	output, err = tester.RunCommand(ctx, "pwd")
	if err != nil {
		log.Fatalf("Test 6 failed: %v", err)
	}
	if tester.ExpectOutput(output, "/") {
		fmt.Println("✓ Test 6 passed: Working directory detected")
	} else {
		fmt.Println("✗ Test 6 failed: Could not detect working directory")
	}

	fmt.Println("\n=== All tests completed ===")
}
