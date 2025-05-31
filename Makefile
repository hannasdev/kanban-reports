# Makefile for Kanban Reports

# Build configuration
BINARY_NAME=kanban-reports
BUILD_DIR=bin
MAIN_PATH=./cmd/kanban-reports

# Go configuration
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --abbrev=0 2>/dev/null || echo 'dev') -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

.PHONY: all build clean test coverage help install run example

# Default target
all: clean test build

# Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests complete"

# Run tests with coverage
coverage:
	@echo "📊 Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies installed"

# Install the binary to GOPATH/bin
install: build
	@echo "📥 Installing $(BINARY_NAME) to GOPATH/bin..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "✅ Installed: $(GOPATH)/bin/$(BINARY_NAME)"

# Run the application in interactive mode
run: build
	@echo "🚀 Running $(BINARY_NAME) in interactive mode..."
	./$(BUILD_DIR)/$(BINARY_NAME) --interactive

# Show example usage
example: build
	@echo "📖 Example usage:"
	./$(BUILD_DIR)/$(BINARY_NAME) --examples

# Show help
help:
	@echo "🔄 Kanban Reports - Make Targets"
	@echo "================================="
	@echo ""
	@echo "Build targets:"
	@echo "  build      Build the application"
	@echo "  clean      Clean build artifacts"
	@echo "  install    Install binary to GOPATH/bin"
	@echo ""
	@echo "Testing targets:"
	@echo "  test       Run all tests"
	@echo "  coverage   Run tests with coverage report"
	@echo ""
	@echo "Development targets:"
	@echo "  deps       Install dependencies"
	@echo "  run        Run application in interactive mode"
	@echo "  example    Show usage examples"
	@echo ""
	@echo "Usage:"
	@echo "  make build              # Build the application"
	@echo "  make test               # Run tests"
	@echo "  make run                # Run interactively"
	@echo "  ./bin/kanban-reports --help  # Show application help"