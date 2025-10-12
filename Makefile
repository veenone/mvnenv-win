.PHONY: build clean test help

# Build variables
BINARY_NAME=mvnenv.exe
BUILD_DIR=bin
VERSION=$(shell cat VERSION)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Default target
all: build

# Build the binary
build:
	@echo Building $(BINARY_NAME)...
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/mvnenv/main.go
	@echo Build complete: $(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo Cleaning...
	@if exist "$(BUILD_DIR)" rmdir /s /q "$(BUILD_DIR)"
	@echo Clean complete

# Run tests
test:
	go test -v ./...

# Display help
help:
	@echo Available targets:
	@echo   build  - Build the binary (default)
	@echo   clean  - Remove build artifacts
	@echo   test   - Run tests
	@echo   help   - Display this help message
