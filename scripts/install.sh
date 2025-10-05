#!/bin/bash

# SSH Keeper - One-line installer
# Usage: curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="ankogit/ssh_keeper"
VERSION="0.1.0"
VERSION_TAG="v0.1.0"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="ssh-keeper"

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""
    
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) os="windows" ;;
        *)          echo -e "${RED}Unsupported operating system: $(uname -s)${NC}" >&2; exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l)     arch="arm" ;;
        *)          echo -e "${RED}Unsupported architecture: $(uname -m)${NC}" >&2; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Download and install
install_ssh_keeper() {
    local platform=$(detect_platform)
    local download_url="https://github.com/${REPO}/releases/download/${VERSION_TAG}/ssh-keeper-${VERSION}-${platform}"
    
    # Add extension for Windows
    if [[ "$platform" == "windows-amd64" ]]; then
        download_url="${download_url}.zip"
        local filename="ssh-keeper-${VERSION}-${platform}.zip"
        local binary_name="ssh-keeper-${VERSION}-${platform}.exe"
    else
        download_url="${download_url}.tar.gz"
        local filename="ssh-keeper-${VERSION}-${platform}.tar.gz"
        local binary_name="ssh-keeper-${VERSION}-${platform}"
    fi
    
    echo -e "${BLUE}🚀 Installing SSH Keeper...${NC}"
    echo -e "${YELLOW}Platform: ${platform}${NC}"
    echo -e "${YELLOW}Version: ${VERSION}${NC}"
    echo -e "${YELLOW}Download URL: ${download_url}${NC}"
    
    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # Download
    echo -e "${BLUE}📥 Downloading SSH Keeper...${NC}"
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$filename"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$download_url" -O "$filename"
    else
        echo -e "${RED}Error: curl or wget is required to download SSH Keeper${NC}" >&2
        exit 1
    fi
    
    # Extract
    echo -e "${BLUE}📦 Extracting archive...${NC}"
    if [[ "$filename" == *.zip ]]; then
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$filename"
        else
            echo -e "${RED}Error: unzip is required to extract the archive${NC}" >&2
            exit 1
        fi
    else
        tar -xzf "$filename"
    fi
    
    # Make executable (Unix systems)
    if [[ "$platform" != "windows-amd64" ]]; then
        chmod +x "$binary_name"
    fi
    
    # Install
    echo -e "${BLUE}🔧 Installing SSH Keeper...${NC}"
    if [[ "$platform" == "windows-amd64" ]]; then
        # Windows: copy to a directory in PATH
        local windows_install_dir="$HOME/bin"
        mkdir -p "$windows_install_dir"
        cp "$binary_name" "$windows_install_dir/$BINARY_NAME.exe"
        echo -e "${GREEN}✅ SSH Keeper installed to ${windows_install_dir}/${BINARY_NAME}.exe${NC}"
        echo -e "${YELLOW}💡 Make sure ${windows_install_dir} is in your PATH${NC}"
    else
        # Unix systems: install to system directory
        if [[ -w "$INSTALL_DIR" ]]; then
            cp "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
        else
            echo -e "${YELLOW}🔐 Installing to system directory (requires sudo)...${NC}"
            sudo cp "$binary_name" "$INSTALL_DIR/$BINARY_NAME"
        fi
        echo -e "${GREEN}✅ SSH Keeper installed to ${INSTALL_DIR}/${BINARY_NAME}${NC}"
    fi
    
    # Cleanup
    cd /
    rm -rf "$temp_dir"
}

# Verify installation
verify_installation() {
    echo -e "${BLUE}🔍 Verifying installation...${NC}"
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version=$("$BINARY_NAME" --version 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✅ SSH Keeper is installed successfully!${NC}"
        echo -e "${GREEN}   Version: ${version}${NC}"
        echo -e "${GREEN}   Location: $(which $BINARY_NAME)${NC}"
        echo ""
        echo -e "${BLUE}🎉 Ready to use! Run '${BINARY_NAME}' to start.${NC}"
    else
        echo -e "${RED}❌ Installation verification failed${NC}" >&2
        echo -e "${YELLOW}💡 Try running '${BINARY_NAME}' manually${NC}"
        exit 1
    fi
}

# Show usage information
show_usage() {
    echo -e "${BLUE}📖 SSH Keeper Usage:${NC}"
    echo -e "${YELLOW}   ${BINARY_NAME}                    # Start SSH Keeper${NC}"
    echo -e "${YELLOW}   ${BINARY_NAME} --version          # Show version${NC}"
    echo -e "${YELLOW}   ${BINARY_NAME} --help             # Show help${NC}"
    echo ""
    echo -e "${BLUE}📚 Documentation:${NC}"
    echo -e "${YELLOW}   https://github.com/${REPO}${NC}"
    echo -e "${YELLOW}   https://github.com/${REPO}/releases${NC}"
}

# Main installation process
main() {
    echo -e "${BLUE}"
    echo "  ____  ____  _  _    _  __  __  _____  _____  _____  _____  _____  _____ "
    echo " / ___||  _ \| || |  | |/ /|  \/  | __||  _  ||  _  ||  _  ||  _  ||  _  |"
    echo " \___ \| |_) | || |_ | ' / | |\/| | _| | | | || | | || | | || | | || | | |"
    echo "  ___) |  __/|__   _|| . \ | |  | | |_ | |_| || |_| || |_| || |_| || |_| |"
    echo " |____/|_|      |_|  |_|\_\|_|  |_|\__| \___/  \___/  \___/  \___/  \___/ "
    echo -e "${NC}"
    echo -e "${BLUE}🔐 Secure SSH Connection Manager${NC}"
    echo ""
    
    # Check if already installed
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  SSH Keeper is already installed at: $(which $BINARY_NAME)${NC}"
        read -p "Do you want to reinstall? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Installation cancelled.${NC}"
            exit 0
        fi
    fi
    
    # Install
    install_ssh_keeper
    verify_installation
    show_usage
}

# Run main function
main "$@"