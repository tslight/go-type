.PHONY: help build test test-verbose test-cover bench bench-textgen bench-cli clean run

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=cli
BINARY_PATH=./$(BINARY_NAME)
MAIN_PKG=./cmd/cli

help: ## Display this help message
	@echo "Go-Type CLI - Makefile targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_PATH) $(MAIN_PKG)
	@echo "✓ Built: $(BINARY_PATH)"

test: ## Run all tests (quick mode)
	@echo "Running tests..."
	@go test ./...
	@echo "✓ All tests passed"

test-verbose: ## Run all tests with verbose output
	@echo "Running tests (verbose)..."
	@go test -v ./...

test-cover: ## Run all tests with coverage report
	@echo "Running tests with coverage..."
	@go test -cover ./...
	@echo "✓ Coverage report complete"

test-textgen: ## Run textgen package tests only
	@echo "Running textgen tests..."
	@go test -v ./internal/textgen

test-cli: ## Run CLI package tests only
	@echo "Running CLI tests..."
	@go test -v ./cmd/cli

bench: bench-textgen bench-cli ## Run all benchmarks

bench-textgen: ## Run textgen package benchmarks
	@echo "Running textgen benchmarks..."
	@go test -bench=. ./internal/textgen -benchmem
	@echo "✓ Benchmarks complete"

bench-cli: ## Run CLI package benchmarks
	@echo "Running CLI benchmarks..."
	@go test -bench=. ./cmd/cli -benchmem
	@echo "✓ Benchmarks complete"

run: build ## Build and run the CLI
	@echo "Running $(BINARY_NAME)..."
	@$(BINARY_PATH) -words 22

run-words: build ## Build and run with custom word count (usage: make run-words WORDS=50)
	@$(BINARY_PATH) -words $(WORDS)

clean: ## Remove build artifacts and test cache
	@echo "Cleaning..."
	@rm -f $(BINARY_PATH)
	@go clean -testcache
	@echo "✓ Cleaned"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ No issues found"

lint: fmt vet ## Run formatters and linters

check: lint test ## Run linters and tests

all: clean check build ## Clean, lint, test, and build
	@echo "✓ All done!"
