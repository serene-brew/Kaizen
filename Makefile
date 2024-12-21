# Variables
APP_NAME := kaizen
GO_FILES := $(wildcard *.go)
BUILD_DIR := build

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	@echo "Building $(APP_NAME) from $(GO_FILES)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(GO_FILES)
	@echo "Build complete! Executable is in $(BUILD_DIR)/$(APP_NAME)"

# Run the app
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup complete!"

# Help target
.PHONY: help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo "Usage:"
	@echo "  make         Build the app (default target)"
	@echo "  make build   Build the app"
	@echo "  make run     Build and run the app"
	@echo "  make clean   Remove build artifacts"
	@echo "  make help    Show this help message"

