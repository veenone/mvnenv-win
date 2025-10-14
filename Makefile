.PHONY: build build-plugins build-all build-shim clean test dist help

# Build variables
BINARY_NAME=mvnenv.exe
SHIM_NAME=shim.exe
BUILD_DIR=bin
DIST_DIR=dist

# Read version from VERSION file
VERSION_FILE=VERSION
VERSION=$(shell type $(VERSION_FILE) 2>nul || cat $(VERSION_FILE))
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

# Plugin build tags
PLUGIN_TAGS=mirror

# Default target
all: build

# Build the binary without plugins
build:
	@echo Building $(BINARY_NAME) standard...
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/mvnenv/main.go
	@echo Build complete: $(BUILD_DIR)/$(BINARY_NAME)

# Build the binary with all plugins enabled
build-plugins:
	@echo Building $(BINARY_NAME) with plugins [$(PLUGIN_TAGS)]...
	@mkdir -p $(BUILD_DIR)
	go build -tags "$(PLUGIN_TAGS)" $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/mvnenv/main.go
	@echo Build complete: $(BUILD_DIR)/$(BINARY_NAME)

# Build shim executable
build-shim:
	@echo Building $(SHIM_NAME)...
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(SHIM_NAME) cmd/shim/main.go
	@echo Build complete: $(BUILD_DIR)/$(SHIM_NAME)

# Build everything (main + plugins + shim)
build-all: build-plugins build-shim
	@echo All builds complete

# Create production distribution package
dist: clean build-plugins build-shim
	@echo Creating production distribution...
	@mkdir -p $(DIST_DIR)/mvnenv-$(VERSION)/bin
	@mkdir -p $(DIST_DIR)/mvnenv-$(VERSION)/config
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(DIST_DIR)/mvnenv-$(VERSION)/bin/
	@cp $(BUILD_DIR)/$(SHIM_NAME) $(DIST_DIR)/mvnenv-$(VERSION)/bin/
	@cp VERSION $(DIST_DIR)/mvnenv-$(VERSION)/
	@cp README.md $(DIST_DIR)/mvnenv-$(VERSION)/
	@cp SETUP.md $(DIST_DIR)/mvnenv-$(VERSION)/
	@cp NEXUS.md $(DIST_DIR)/mvnenv-$(VERSION)/
	@cp PLUGINS.md $(DIST_DIR)/mvnenv-$(VERSION)/
	@cp config.example.yaml $(DIST_DIR)/mvnenv-$(VERSION)/config/
	@echo
	@echo Production distribution created: $(DIST_DIR)/mvnenv-$(VERSION)
	@echo
	@echo Contents:
	@echo   - bin/$(BINARY_NAME) with plugins
	@echo   - bin/$(SHIM_NAME)
	@echo   - VERSION
	@echo   - README.md, SETUP.md, NEXUS.md, PLUGINS.md
	@echo   - config/config.example.yaml
	@echo
	@echo To install:
	@echo   1. Copy the mvnenv-$(VERSION) directory to a permanent location
	@echo   2. Add the bin directory to your PATH
	@echo   3. Run: mvnenv.exe rehash

# Clean build artifacts
clean:
	@echo Cleaning...
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo Clean complete

# Run tests
test:
	go test -v ./...

# Display help
help:
	@echo Available targets:
	@echo   build          - Build the binary without plugins
	@echo   build-plugins  - Build the binary with all plugins enabled
	@echo   build-shim     - Build the shim executable
	@echo   build-all      - Build everything
	@echo   dist           - Create production distribution package
	@echo   clean          - Remove build artifacts
	@echo   test           - Run tests
	@echo   help           - Display this help message
	@echo
	@echo Plugins: $(PLUGIN_TAGS)
	@echo Version: $(VERSION)
