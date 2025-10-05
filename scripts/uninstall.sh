#!/bin/bash

# SSH Keeper - Uninstaller
# Usage: curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/uninstall.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="ssh-keeper"
CONFIG_DIR="$HOME/.ssh-keeper"

# Detect platform
detect_platform() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) echo "windows" ;;
        *)          echo "unknown" ;;
    esac
}

# Uninstall SSH Keeper
uninstall_ssh_keeper() {
    local platform=$(detect_platform)
    local found=false
    
    echo -e "${BLUE}🗑️  Uninstalling SSH Keeper...${NC}"
    
    # Check different installation locations
    local locations=()
    
    if [[ "$platform" == "windows" ]]; then
        locations=(
            "$HOME/bin/$BINARY_NAME.exe"
            "$HOME/AppData/Local/bin/$BINARY_NAME.exe"
            "C:/Program Files/$BINARY_NAME/$BINARY_NAME.exe"
            "C:/Program Files (x86)/$BINARY_NAME/$BINARY_NAME.exe"
        )
    else
        locations=(
            "/usr/local/bin/$BINARY_NAME"
            "/usr/bin/$BINARY_NAME"
            "$HOME/bin/$BINARY_NAME"
            "/opt/$BINARY_NAME/$BINARY_NAME"
        )
    fi
    
    # Remove binary files
    for location in "${locations[@]}"; do
        if [[ -f "$location" ]]; then
            echo -e "${YELLOW}📁 Found binary at: $location${NC}"
            if [[ -w "$(dirname "$location")" ]]; then
                rm -f "$location"
                echo -e "${GREEN}✅ Removed: $location${NC}"
            else
                echo -e "${YELLOW}🔐 Removing system file (requires sudo)...${NC}"
                sudo rm -f "$location"
                echo -e "${GREEN}✅ Removed: $location${NC}"
            fi
            found=true
        fi
    done
    
    # Remove configuration directory
    if [[ -d "$CONFIG_DIR" ]]; then
        echo -e "${YELLOW}📁 Found config directory: $CONFIG_DIR${NC}"
        read -p "Do you want to remove configuration files? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$CONFIG_DIR"
            echo -e "${GREEN}✅ Removed configuration directory${NC}"
        else
            echo -e "${YELLOW}💡 Configuration directory preserved: $CONFIG_DIR${NC}"
        fi
    fi
    
    # Remove from PATH (optional)
    if [[ "$platform" != "windows" ]]; then
        local shell_configs=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile" "$HOME/.profile")
        for config in "${shell_configs[@]}"; do
            if [[ -f "$config" ]] && grep -q "$HOME/bin" "$config" 2>/dev/null; then
                echo -e "${YELLOW}📝 Found PATH entry in: $config${NC}"
                read -p "Do you want to remove PATH entry? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    # Remove the PATH export line
                    sed -i.bak '/export PATH="\$HOME\/bin:\$PATH"/d' "$config"
                    echo -e "${GREEN}✅ Removed PATH entry from $config${NC}"
                fi
            fi
        done
    fi
    
    if [[ "$found" == "true" ]]; then
        echo -e "${GREEN}✅ SSH Keeper uninstalled successfully!${NC}"
    else
        echo -e "${YELLOW}⚠️  SSH Keeper not found in common locations${NC}"
        echo -e "${BLUE}💡 You may need to remove it manually${NC}"
    fi
}

# Verify uninstallation
verify_uninstallation() {
    echo -e "${BLUE}🔍 Verifying uninstallation...${NC}"
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  SSH Keeper is still available at: $(which $BINARY_NAME)${NC}"
        echo -e "${BLUE}💡 You may need to restart your terminal or run: source ~/.bashrc${NC}"
    else
        echo -e "${GREEN}✅ SSH Keeper is no longer available${NC}"
    fi
}

# Show cleanup information
show_cleanup_info() {
    echo -e "${BLUE}📚 Manual cleanup (if needed):${NC}"
    echo -e "${YELLOW}   • Configuration: $CONFIG_DIR${NC}"
    echo -e "${YELLOW}   • Shell configs: ~/.bashrc, ~/.zshrc${NC}"
    echo -e "${YELLOW}   • Check PATH: echo \$PATH${NC}"
    echo ""
    echo -e "${BLUE}📖 Documentation:${NC}"
    echo -e "${YELLOW}   https://github.com/ankogit/ssh_keeper${NC}"
}

# Main uninstallation process
main() {
    echo -e "${BLUE}"
    echo "  ╔══════════════════════════════════════════════════════════════╗"
    echo "  ║                                                              ║"
    echo "  ║                    🗑️  SSH KEEPER 🗑️                           ║"
    echo "  ║                                                              ║"
    echo "  ║                    Uninstaller                               ║"
    echo "  ║                                                              ║"
    echo "  ╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Check if running in non-interactive mode (curl | bash)
    if [[ ! -t 0 ]]; then
        echo -e "${YELLOW}⚠️  Running in non-interactive mode. Proceeding with uninstallation...${NC}"
        echo ""
    else
        # Confirm uninstallation only in interactive mode
        read -p "Are you sure you want to uninstall SSH Keeper? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Uninstallation cancelled.${NC}"
            exit 0
        fi
    fi
    
    # Uninstall
    uninstall_ssh_keeper
    verify_uninstallation
    show_cleanup_info
}

# Run main function
main "$@"
