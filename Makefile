.PHONY: test test-verbose test-coverage examples clean fmt lint help

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  examples      - Run all examples"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter (requires golangci-lint)"
	@echo "  clean         - Clean build artifacts"

# Run tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run all examples
examples:
	@echo "Running basic example..."
	@cd examples && timeout 3 go run basic.go || true
	@echo ""
	@echo "Running vim_automation example..."
	@cd examples && timeout 5 go run vim_automation.go || true
	@echo ""
	@echo "Running event_streaming example..."
	@cd examples && timeout 5 go run event_streaming.go || true
	@echo ""
	@echo "Running cli_testing example..."
	@cd examples && timeout 5 go run cli_testing.go || true

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/"; exit 1)
	golangci-lint run

# Clean build artifacts
clean:
	rm -f coverage.out coverage.html
	rm -f *.test
	go clean
