.PHONY: build run run-once clean test help

# Build the application
build:
	@echo "Building CloudFlare Backuper..."
	go build -o cloudflare-backuper

# Run the application as a service
run: build
	@echo "Starting CloudFlare Backuper service..."
	./cloudflare-backuper

# Run a single backup and exit
run-once: build
	@echo "Running one-time backup..."
	./cloudflare-backuper -once

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f cloudflare-backuper
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run static analysis
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from: https://golangci-lint.run/usage/install/"; \
	fi

# Display help
help:
	@echo "CloudFlare Backuper - Makefile Commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make build      - Build the application binary"
	@echo "  make run        - Build and run as a service"
	@echo "  make run-once   - Build and run a single backup"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Install/update dependencies"
	@echo "  make fmt        - Format source code"
	@echo "  make lint       - Run static analysis (requires golangci-lint)"
	@echo "  make help       - Display this help message"
	@echo ""
	@echo "Before running, make sure to:"
	@echo "  1. Copy config.example.yml to config.yml"
	@echo "  2. Edit config.yml with your CloudFlare and Discord credentials"
	@echo ""
	@echo "Example:"
	@echo "  cp config.example.yml config.yml"
	@echo "  nano config.yml"
	@echo "  make run"
