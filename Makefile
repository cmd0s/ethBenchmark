# Makefile for ethbench - Ethereum Node Benchmark Tool

VERSION := 0.1.0
BINARY := ethbench
BUILD_DIR := build

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOGET := $(GOCMD) get

# Build flags
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# Target architectures
.PHONY: all build build-arm64 build-all clean test deps tidy help

all: build

# Download dependencies
deps:
	$(GOMOD) download

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Build for current platform
build: deps
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/ethbench
	@echo "Built: $(BUILD_DIR)/$(BINARY)"

# Build for Raspberry Pi 5 (ARM64 Linux)
build-arm64: deps
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/ethbench
	@echo "Built: $(BUILD_DIR)/$(BINARY)-linux-arm64"

# Build for AMD64 Linux
build-amd64: deps
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/ethbench
	@echo "Built: $(BUILD_DIR)/$(BINARY)-linux-amd64"

# Build for macOS ARM64 (Apple Silicon)
build-darwin-arm64: deps
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/ethbench
	@echo "Built: $(BUILD_DIR)/$(BINARY)-darwin-arm64"

# Build for macOS AMD64
build-darwin-amd64: deps
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/ethbench
	@echo "Built: $(BUILD_DIR)/$(BINARY)-darwin-amd64"

# Build for all platforms
build-all: build-arm64 build-amd64 build-darwin-arm64 build-darwin-amd64
	@echo "Built all platforms in $(BUILD_DIR)/"

# Create release archives
release: build-all
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && tar -czvf release/$(BINARY)-$(VERSION)-linux-arm64.tar.gz $(BINARY)-linux-arm64
	@cd $(BUILD_DIR) && tar -czvf release/$(BINARY)-$(VERSION)-linux-amd64.tar.gz $(BINARY)-linux-amd64
	@cd $(BUILD_DIR) && tar -czvf release/$(BINARY)-$(VERSION)-darwin-arm64.tar.gz $(BINARY)-darwin-arm64
	@cd $(BUILD_DIR) && tar -czvf release/$(BINARY)-$(VERSION)-darwin-amd64.tar.gz $(BINARY)-darwin-amd64
	@echo "Release archives created in $(BUILD_DIR)/release/"

# Run tests
test:
	$(GOTEST) -v ./...

# Run benchmark (quick mode for testing)
run: build
	./$(BUILD_DIR)/$(BINARY) -quick

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f ethbench-*.json

# Show help
help:
	@echo "Ethereum Node Benchmark Tool - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-arm64    Build for Raspberry Pi 5 (ARM64 Linux)"
	@echo "  make build-amd64    Build for AMD64 Linux"
	@echo "  make build-all      Build for all platforms"
	@echo "  make release        Create release archives"
	@echo "  make test           Run tests"
	@echo "  make run            Build and run quick benchmark"
	@echo "  make clean          Remove build artifacts"
	@echo "  make deps           Download dependencies"
	@echo "  make tidy           Tidy dependencies"
	@echo "  make help           Show this help"
