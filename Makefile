.PHONY: help build test test-coverage clean run web install lint fmt vet check deps

# Variables
BINARY_NAME=akeneo-migrator
BUILD_DIR=./bin
CMD_DIR=./cmd/app
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)          - Build the application binary"
	@echo "  $(COLOR_GREEN)test$(COLOR_RESET)           - Run all tests"
	@echo "  $(COLOR_GREEN)test-coverage$(COLOR_RESET) - Run tests with coverage report"
	@echo "  $(COLOR_GREEN)clean$(COLOR_RESET)          - Remove build artifacts"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)            - Build and run the application"
	@echo "  $(COLOR_GREEN)web$(COLOR_RESET)            - Build and start the web UI"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)        - Install dependencies"
	@echo "  $(COLOR_GREEN)lint$(COLOR_RESET)           - Run linter (golangci-lint)"
	@echo "  $(COLOR_GREEN)fmt$(COLOR_RESET)            - Format code"
	@echo "  $(COLOR_GREEN)vet$(COLOR_RESET)            - Run go vet"
	@echo "  $(COLOR_GREEN)check$(COLOR_RESET)          - Run all checks (fmt, vet, test)"
	@echo "  $(COLOR_GREEN)deps$(COLOR_RESET)           - Download dependencies"

## build: Build the application binary
build:
	@echo "$(COLOR_BLUE)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(COLOR_GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	@$(GO) test $(GOFLAGS) ./...
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	@$(GO) test -cover -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

## clean: Remove build artifacts
clean:
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## run: Build and run the application
run: build
	@echo "$(COLOR_BLUE)Running $(BINARY_NAME)...$(COLOR_RESET)"
	@$(BUILD_DIR)/$(BINARY_NAME)

## web: Build and start the web UI
web: build
	@echo "$(COLOR_BLUE)Starting Web UI...$(COLOR_RESET)"
	@$(BUILD_DIR)/$(BINARY_NAME) web

## install: Install dependencies
install:
	@echo "$(COLOR_BLUE)Installing dependencies...$(COLOR_RESET)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

## deps: Download dependencies
deps:
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	@$(GO) mod download
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "$(COLOR_BLUE)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not installed. Install it from https://golangci-lint.run/$(COLOR_RESET)"; \
	fi

## fmt: Format code
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	@$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	@$(GO) vet ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## check: Run all checks (fmt, vet, test)
check: fmt vet test
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

# Default target
.DEFAULT_GOAL := help
