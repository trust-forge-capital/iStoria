# iStoria Makefile
# Cross-platform hardware monitoring CLI

.PHONY: all build test clean install run help

# Variables
BINARY_NAME=istoria
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

# Output directory
DIST=dist

# Default target
all: test build

# Build for current platform
build:
	@mkdir -p $(DIST)
	go build ${LDFLAGS} -o $(DIST)/${BINARY_NAME} .

# Build for all platforms
build-all:
	@mkdir -p $(DIST)
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o $(DIST)/istoria-darwin-arm64 .
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o $(DIST)/istoria-darwin-amd64 .
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $(DIST)/istoria-linux-amd64 .
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o $(DIST)/istoria-linux-arm64 .
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o $(DIST)/istoria-windows-amd64.exe .
	@echo "Building for Windows (arm64)..."
	GOOS=windows GOARCH=arm64 go build ${LDFLAGS} -o $(DIST)/istoria-windows-arm64.exe .
	@echo "Build complete!"

# Run tests
test:
	go test -v -cover ./...

# Run tests with race detector
test-race:
	go test -v -race -cover ./...

# Clean build artifacts
clean:
	rm -rf $(DIST)
	rm -f coverage.out

# Install to GOPATH/bin
install:
	go install ${LDFLAGS}

# Run the application
run:
	go run main.go

# Run with specific command
run-cpu:
	go run main.go cpu

run-mem:
	go run main.go mem

run-disk:
	go run main.go disk

run-net:
	go run main.go net

run-sensor:
	go run main.go sensor

run-power:
	go run main.go power

run-stat:
	go run main.go stat

# Live mode
run-live:
	go run main.go cpu --live

# Show help
help:
	@echo "iStoria Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build for current platform (output to dist/)"
	@echo "  make build-all   - Build for all platforms (output to dist/)"
	@echo "  make test        - Run tests"
	@echo "  make test-race   - Run tests with race detector"
	@echo "  make clean       - Clean dist/ directory"
	@echo "  make install     - Install to GOPATH/bin"
	@echo "  make run         - Run the application"
	@echo "  make run-cpu    - Run cpu command"
	@echo "  make run-mem    - Run mem command"
	@echo "  make run-disk   - Run disk command"
	@echo "  make run-net    - Run net command"
	@echo "  make run-sensor - Run sensor command"
	@echo "  make run-power  - Run power command"
	@echo "  make run-stat   - Run stat command"
	@echo "  make run-live   - Run in live mode"
	@echo "  make help       - Show this help"

# Development targets
dev: test build
	@echo "Dev build complete!"

# Release target
release: clean test build-all
	@echo "Release build complete!"

# Show built artifacts
list:
	@echo "Built artifacts:"
	@ls -la $(DIST)/ 2>/dev/null || echo "(none)"
