# Makefile

# Variables
GO_MAIN := ./cmd # Adjust this to the directory where your main.go is located
BUILD_DIR := ./dist
EXECUTABLE := $(BUILD_DIR)/myapp

# Targets
.PHONY: all build wire test run clean

all: wire test run

build:
	@echo "Building the project..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(EXECUTABLE) $(GO_MAIN)

wire:
	@echo "Running wire for dependency injection..."
	@wire ./...

test:
	@echo "Running tests in verbose mode..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

run: build
	@echo "Running the executable..."
	@$(EXECUTABLE)

clean:
	@echo "Cleaning build files..."
	@go clean
	@rm -rf $(BUILD_DIR)