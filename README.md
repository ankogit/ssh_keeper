# SSH Keeper 🔐

<div align="center">

**A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface**

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey.svg)](https://github.com/ankogit/ssh_keeper/releases)

_Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and Go_

</div>

---

## ✨ Features

- 🎨 **Beautiful TUI Interface** - Modern terminal user interface with colors and smooth animations
- 🔐 **Secure Password Storage** - Master password with system keyring integration (macOS Keychain, Linux Secret Service, Windows Credential Manager)
- 🔑 **Dual Authentication** - Support for both password and SSH key authentication
- 📁 **Connection Management** - Add, edit, delete, and organize your SSH connections
- 🔍 **Smart Search** - Quick connection search and filtering
- 📤 **Export/Import** - Full compatibility with OpenSSH config format
- ⚡ **Fast & Lightweight** - Built with Go for optimal performance
- 🌍 **Cross-Platform** - Works on macOS, Linux, and Windows
- 🔒 **Open Source** - MIT licensed, community-driven development

## 🚀 Quick Start

### One-line Installation (Recommended)

**macOS & Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex
```

### Manual Download & Install

Download the latest release for your platform:

- **macOS Intel**: [ssh-keeper-0.1.0-darwin-amd64.tar.gz](https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-amd64.tar.gz) (3.08 MiB)
- **macOS Apple Silicon**: [ssh-keeper-0.1.0-darwin-arm64.tar.gz](https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-darwin-arm64.tar.gz) (2.94 MiB)
- **Linux**: [ssh-keeper-0.1.0-linux-amd64.tar.gz](https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-linux-amd64.tar.gz) (3.84 MiB)
- **Windows**: [ssh-keeper-0.1.0-windows-amd64.zip](https://github.com/ankogit/ssh_keeper/releases/download/v0.1.0/ssh-keeper-0.1.0-windows-amd64.zip) (3.33 MiB)

### Extract and Run

```bash
# Extract the archive
tar -xzf ssh-keeper-0.1.0-*.tar.gz  # Linux/macOS
# or unzip ssh-keeper-0.1.0-windows-amd64.zip  # Windows

# Make executable (Linux/macOS)
chmod +x ssh-keeper*

# Run
./ssh-keeper
```

> **📖 Подробная инструкция по установке**: [INSTALL.md](INSTALL.md) | [macOS Apple Silicon](INSTALL_MACOS.md)

## 📸 Screenshots

> **Note**: Screenshots will be added here once the application is running

### Main Menu

![Main Menu](docs/screenshots/main-menu.png)
_Beautiful main menu with intuitive navigation_

### Connection List

![Connection List](docs/screenshots/connections.png)
_View and manage your SSH connections_

### Add Connection

![Add Connection](docs/screenshots/add-connection.png)
_Easy connection setup with form validation_

### Settings

![Settings](docs/screenshots/settings.png)
_Configure your SSH Keeper preferences_

## 🛠️ Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ankogit/ssh_keeper.git
cd ssh_keeper

# Build the application
make build

# Run the application
make run
```

### Using Make

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run in development mode
make run-dev

# Install system-wide
make install
```

### Prerequisites

- **Go 1.24+** - [Download Go](https://golang.org/dl/)
- **Make** (optional, for using Makefile)

## 📖 Usage

### First Run

When you first run SSH Keeper, you'll be prompted to set up a master password. This password is used to encrypt your connection data and is stored securely in your system's keyring.

```bash
ssh-keeper
```

### Main Menu Navigation

The application provides an intuitive main menu with the following options:

- **🔍 View Connections** - Browse and search your SSH connections
- **➕ Add Connection** - Add a new SSH connection
- **⚙️ Settings** - Configure application settings
- **📤 Export** - Export connections to OpenSSH config
- **📥 Import** - Import connections from OpenSSH config
- **❌ Quit** - Exit the application

### Keyboard Shortcuts

| Key      | Action              |
| -------- | ------------------- |
| `↑/↓`    | Navigate menu items |
| `Enter`  | Select item         |
| `Ctrl+S` | Search connections  |
| `Esc`    | Go back             |
| `Q`      | Quit application    |

### Adding Connections

1. Select "➕ Add Connection" from the main menu
2. Fill in the connection details:
   - **Name**: A friendly name for your connection
   - **Host**: Server hostname or IP address
   - **Port**: SSH port (default: 22)
   - **User**: Username for SSH connection
   - **Authentication**: Choose between password or SSH key
3. Save your connection

### Authentication Methods

#### Password Authentication

- Enter your password when adding the connection
- Password is encrypted and stored securely
- No need to enter password each time you connect

#### SSH Key Authentication

- Specify the path to your SSH private key
- Supports standard SSH key formats
- Works with existing SSH key infrastructure

## ⚙️ Configuration

SSH Keeper stores its configuration in `~/.ssh-keeper/`:

- `config.yaml` - Application settings
- `connections.yaml` - SSH connections (encrypted)
- Passwords are stored securely using system keyring

### Environment Variables

Create a `.env` file for development:

```bash
cp env.example .env
```

| Variable          | Description                          | Default                | Required |
| ----------------- | ------------------------------------ | ---------------------- | -------- |
| `DEBUG`           | Enable debug mode                    | `false`                | No       |
| `ENV`             | Environment (development/production) | `development`          | No       |
| `CONFIG_PATH`     | Path to application config file      | `~/.ssh-keeper/config` | No       |
| `APP_SIGNATURE`   | Application signature for security   | -                      | Yes      |
| `SSH_CONFIG_PATH` | Path to SSH config file              | `~/.ssh/config`        | No       |

## 🔒 Security

SSH Keeper prioritizes security and follows best practices:

- **Encrypted Storage**: All connection data is encrypted using AES-256
- **System Keyring Integration**: Master password stored in system keyring (Keychain/Secret Service/Credential Manager)
- **Memory Management**: Sensitive data cleared from memory after use
- **No Plain Text**: No passwords stored in plain text files
- **Open Source**: Full source code available for security audit

## 🏗️ Development

### Project Structure

```
ssh_keeper/
├── cmd/ssh-keeper/          # Main application entry point
├── internal/
│   ├── models/              # Data models (Connection, Config)
│   ├── ui/                  # TUI components and screens
│   │   ├── components/      # Reusable UI components
│   │   └── screens/         # Application screens
│   ├── services/            # Business logic services
│   ├── ssh/                 # SSH client integration
│   └── config/              # Configuration management
├── docs/                    # Documentation
├── scripts/                 # Build and utility scripts
├── Formula/                 # Homebrew formulas
└── Makefile                 # Build automation
```

### Development Setup

```bash
# Clone the repository
git clone https://github.com/ankogit/ssh_keeper.git
cd ssh_keeper

# Download dependencies
make deps

# Set up development environment
make dev-setup

# Run tests
make test

# Run with coverage
make test-coverage
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Create release packages
make release
```

## 🤝 Contributing

We welcome contributions! SSH Keeper is an open source project under the MIT license.

### How to Contribute

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure all tests pass before submitting PR

### Reporting Issues

Found a bug? Have a feature request? Please [open an issue](https://github.com/ankogit/ssh_keeper/issues)!

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### Why MIT License?

The MIT License is one of the most permissive open source licenses, allowing:

- ✅ **Commercial use** - Use in commercial projects
- ✅ **Modification** - Modify and create derivatives
- ✅ **Distribution** - Distribute copies
- ✅ **Private use** - Use privately
- ✅ **Patent use** - Use patented technology

This makes SSH Keeper accessible to everyone while maintaining the freedom to use, modify, and distribute the software.

## 🙏 Acknowledgments

SSH Keeper is built with amazing open source projects:

- [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) - A powerful TUI framework for Go
- [**Bubbles**](https://github.com/charmbracelet/bubbles) - Beautiful UI components
- [**Lip Gloss**](https://github.com/charmbracelet/lipgloss) - Styling for terminal applications
- [**go-keyring**](https://github.com/zalando/go-keyring) - Cross-platform keyring access
- [**Termenv**](https://github.com/muesli/termenv) - Terminal environment detection

## 📊 Project Status

- ✅ **Core Features** - Connection management, secure storage
- ✅ **TUI Interface** - Beautiful terminal interface
- ✅ **Cross-Platform** - macOS, Linux, Windows support
- ✅ **Security** - Encrypted storage, system keyring integration
- 🔄 **Documentation** - Comprehensive docs and examples
- 🔄 **Testing** - Unit tests and integration tests
- 🔄 **Packaging** - Homebrew, Debian packages

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ankogit/ssh_keeper&type=Date)](https://star-history.com/#ankogit/ssh_keeper&Date)

---

<div align="center">

**Made with ❤️ by the SSH Keeper community**

[⭐ Star this project](https://github.com/ankogit/ssh_keeper) | [🐛 Report Bug](https://github.com/ankogit/ssh_keeper/issues) | [💡 Request Feature](https://github.com/ankogit/ssh_keeper/issues) | [📖 Documentation](https://github.com/ankogit/ssh_keeper/tree/main/docs)

</div>
