# SSH Keeper Installation Guide

## Quick Installation

### One-line installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/ssh-keeper/main/scripts/install.sh | bash
```

## Manual Installation Methods

### macOS (Homebrew)

```bash
# Add tap and install
brew tap yourusername/ssh-keeper
brew install ssh-keeper
```

### Ubuntu/Debian (apt)

```bash
# Download and install .deb package
wget https://github.com/yourusername/ssh-keeper/releases/latest/download/ssh-keeper_0.1.0_amd64.deb
sudo dpkg -i ssh-keeper_0.1.0_amd64.deb
sudo apt-get install -f  # Fix dependencies if needed
```

### Manual Binary Installation

#### Linux

```bash
# Download latest release
wget https://github.com/yourusername/ssh-keeper/releases/latest/download/ssh-keeper-0.1.0-linux-amd64.tar.gz
tar -xzf ssh-keeper-0.1.0-linux-amd64.tar.gz
sudo mv ssh-keeper-linux-amd64 /usr/local/bin/ssh-keeper
sudo chmod +x /usr/local/bin/ssh-keeper
```

#### macOS

```bash
# For Intel Macs
wget https://github.com/yourusername/ssh-keeper/releases/latest/download/ssh-keeper-0.1.0-darwin-amd64.tar.gz
tar -xzf ssh-keeper-0.1.0-darwin-amd64.tar.gz
sudo mv ssh-keeper-darwin-amd64 /usr/local/bin/ssh-keeper
sudo chmod +x /usr/local/bin/ssh-keeper

# For Apple Silicon Macs
wget https://github.com/yourusername/ssh-keeper/releases/latest/download/ssh-keeper-0.1.0-darwin-arm64.tar.gz
tar -xzf ssh-keeper-0.1.0-darwin-arm64.tar.gz
sudo mv ssh-keeper-darwin-arm64 /usr/local/bin/ssh-keeper
sudo chmod +x /usr/local/bin/ssh-keeper
```

#### Windows

```bash
# Download and extract
wget https://github.com/yourusername/ssh-keeper/releases/latest/download/ssh-keeper-0.1.0-windows-amd64.zip
unzip ssh-keeper-0.1.0-windows-amd64.zip
# Move ssh-keeper-windows-amd64.exe to your PATH
```

## Verification

After installation, verify that ssh-keeper is working:

```bash
ssh-keeper --version
```

## Uninstallation

### Homebrew

```bash
brew uninstall ssh-keeper
brew untap yourusername/ssh-keeper
```

### apt

```bash
sudo apt-get remove ssh-keeper
```

### Manual

```bash
sudo rm /usr/local/bin/ssh-keeper
```
