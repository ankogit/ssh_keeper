# SSH Keeper Makefile

# Variables
BINARY_NAME=ssh-keeper
VERSION=0.1.0
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test run install uninstall help

# Default target
all: clean build

# Build the application
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ssh-keeper
	@echo "$(GREEN)Build completed!$(NC)"

# Build for multiple platforms
build-all:
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	@echo "$(YELLOW)Building for Linux...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/ssh-keeper
	
	# macOS
	@echo "$(YELLOW)Building for macOS...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/ssh-keeper
	@GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/ssh-keeper
	
	# Windows
	@echo "$(YELLOW)Building for Windows...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/ssh-keeper
	
	@echo "$(GREEN)Multi-platform build completed!$(NC)"

# Run the application
run: build
	@echo "$(BLUE)Running $(BINARY_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run without building
run-dev:
	@echo "$(BLUE)Running in development mode...$(NC)"
	@go run ./cmd/ssh-keeper

# Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean completed!$(NC)"

# Install the application
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME)...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)Installation completed!$(NC)"

# Uninstall the application
uninstall:
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)Uninstallation completed!$(NC)"

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)Code formatted!$(NC)"

# Lint code
lint:
	@echo "$(BLUE)Linting code...$(NC)"
	@go vet ./...
	@echo "$(GREEN)Linting completed!$(NC)"

# Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Dependencies downloaded!$(NC)"

# Create release packages
release: build-all
	@echo "$(BLUE)Creating release packages...$(NC)"
	@mkdir -p $(BUILD_DIR)/release
	
	# Linux
	@tar -czf $(BUILD_DIR)/release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	
	# macOS
	@tar -czf $(BUILD_DIR)/release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	@tar -czf $(BUILD_DIR)/release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	
	# Windows
	@zip -j $(BUILD_DIR)/release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	@echo "$(GREEN)Release packages created in $(BUILD_DIR)/release/$(NC)"

# Create Debian package
deb: build
	@echo "$(BLUE)Creating Debian package...$(NC)"
	@./build-deb.sh
	@echo "$(GREEN)Debian package created!$(NC)"

# Create Homebrew formula
homebrew-formula:
	@echo "$(BLUE)Creating Homebrew formula...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cp Formula/ssh-keeper.rb $(BUILD_DIR)/
	@echo "$(GREEN)Homebrew formula created in $(BUILD_DIR)/$(NC)"

# Create GitHub release
github-release: release deb homebrew-formula
	@echo "$(BLUE)Creating GitHub release...$(NC)"
	@echo "$(YELLOW)Please create a GitHub release manually with the following files:$(NC)"
	@echo "$(BUILD_DIR)/release/*.tar.gz"
	@echo "$(BUILD_DIR)/release/*.zip"
	@echo "$(BUILD_DIR)/*.deb"
	@echo "$(BUILD_DIR)/*.changes"
	@echo "$(BUILD_DIR)/*.dsc"
	@echo "$(BUILD_DIR)/ssh-keeper.rb"

# Development setup
dev-setup: deps
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)Development setup completed!$(NC)"

# Show help
help:
	@echo "$(BLUE)SSH Keeper - Available commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Build commands:$(NC)"
	@echo "  build        Build the application"
	@echo "  build-all    Build for multiple platforms"
	@echo "  clean        Clean build artifacts"
	@echo ""
	@echo "$(YELLOW)Run commands:$(NC)"
	@echo "  run          Build and run the application"
	@echo "  run-dev      Run in development mode"
	@echo ""
	@echo "$(YELLOW)Test commands:$(NC)"
	@echo "  test         Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo ""
	@echo "$(YELLOW)Install commands:$(NC)"
	@echo "  install      Install the application"
	@echo "  uninstall    Uninstall the application"
	@echo ""
	@echo "$(YELLOW)Development commands:$(NC)"
	@echo "  fmt          Format code"
	@echo "  lint         Lint code"
	@echo "  deps         Download dependencies"
	@echo "  dev-setup    Set up development environment"
	@echo ""
	@echo "$(YELLOW)Release commands:$(NC)"
	@echo "  release      Create release packages"
	@echo "  deb          Create Debian package"
	@echo "  homebrew-formula Create Homebrew formula"
	@echo "  github-release Create all release files"
	@echo ""
	@echo "$(YELLOW)Other commands:$(NC)"
	@echo "  help         Show this help message"
