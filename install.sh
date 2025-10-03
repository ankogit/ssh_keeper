#!/bin/bash

# SSH Keeper Installation Script
# Supports macOS (Homebrew), Ubuntu/Debian (apt), and manual installation

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="yourusername/ssh-keeper"
BINARY_NAME="ssh-keeper"
INSTALL_DIR="/usr/local/bin"

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if command -v apt-get >/dev/null 2>&1; then
            echo "ubuntu"
        elif command -v yum >/dev/null 2>&1; then
            echo "centos"
        else
            echo "linux"
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "macos"
    else
        echo "unknown"
    fi
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install binary
install_binary() {
    local version=$1
    local os=$2
    local arch=$3
    
    echo -e "${BLUE}Downloading ssh-keeper $version for $os-$arch...${NC}"
    
    # Determine file extension
    if [[ "$os" == "windows" ]]; then
        ext="zip"
        binary_name="${BINARY_NAME}-windows-amd64.exe"
    else
        ext="tar.gz"
        binary_name="${BINARY_NAME}-${os}-${arch}"
    fi
    
    # Download URL
    local url="https://github.com/$REPO/releases/download/$version/${BINARY_NAME}-${version}-${os}-${arch}.${ext}"
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download and extract
    if [[ "$ext" == "zip" ]]; then
        curl -L "$url" -o "${BINARY_NAME}.zip"
        unzip "${BINARY_NAME}.zip"
    else
        curl -L "$url" | tar -xz
    fi
    
    # Install binary
    sudo mv "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    # Cleanup
    cd /
    rm -rf "$temp_dir"
    
    echo -e "${GREEN}ssh-keeper installed successfully!${NC}"
}

# Install via Homebrew (macOS)
install_homebrew() {
    echo -e "${BLUE}Installing via Homebrew...${NC}"
    
    if ! command -v brew >/dev/null 2>&1; then
        echo -e "${YELLOW}Homebrew not found. Installing Homebrew first...${NC}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi
    
    # Add tap and install
    brew tap "$REPO/ssh-keeper"
    brew install ssh-keeper
    
    echo -e "${GREEN}ssh-keeper installed via Homebrew!${NC}"
}

# Install via apt (Ubuntu/Debian)
install_apt() {
    echo -e "${BLUE}Installing via apt...${NC}"
    
    # Download and install .deb package
    local version=$(get_latest_version)
    local url="https://github.com/$REPO/releases/download/$version/ssh-keeper_${version#v}_amd64.deb"
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download .deb package
    curl -L "$url" -o "ssh-keeper.deb"
    
    # Install package
    sudo dpkg -i ssh-keeper.deb
    sudo apt-get install -f  # Fix any dependencies
    
    # Cleanup
    cd /
    rm -rf "$temp_dir"
    
    echo -e "${GREEN}ssh-keeper installed via apt!${NC}"
}

# Main installation function
main() {
    echo -e "${BLUE}SSH Keeper Installation Script${NC}"
    echo -e "${BLUE}==============================${NC}"
    
    local os=$(detect_os)
    echo -e "${YELLOW}Detected OS: $os${NC}"
    
    # Check if already installed
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${YELLOW}ssh-keeper is already installed.${NC}"
        read -p "Do you want to update it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Installation cancelled.${NC}"
            exit 0
        fi
    fi
    
    # Get latest version
    local version=$(get_latest_version)
    echo -e "${YELLOW}Latest version: $version${NC}"
    
    # Choose installation method
    case "$os" in
        "macos")
            if command -v brew >/dev/null 2>&1; then
                echo -e "${BLUE}Homebrew detected. Installing via Homebrew...${NC}"
                install_homebrew
            else
                echo -e "${BLUE}Installing binary directly...${NC}"
                # Detect architecture
                if [[ $(uname -m) == "arm64" ]]; then
                    install_binary "$version" "darwin" "arm64"
                else
                    install_binary "$version" "darwin" "amd64"
                fi
            fi
            ;;
        "ubuntu")
            echo -e "${BLUE}Installing via apt...${NC}"
            install_apt
            ;;
        "linux"|"centos")
            echo -e "${BLUE}Installing binary directly...${NC}"
            install_binary "$version" "linux" "amd64"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $os${NC}"
            echo -e "${YELLOW}Please install manually from: https://github.com/$REPO/releases${NC}"
            exit 1
            ;;
    esac
    
    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${GREEN}Installation completed successfully!${NC}"
        echo -e "${BLUE}Run 'ssh-keeper' to start the application.${NC}"
    else
        echo -e "${RED}Installation failed. Please try manual installation.${NC}"
        exit 1
    fi
}

# Run main function
main "$@"

