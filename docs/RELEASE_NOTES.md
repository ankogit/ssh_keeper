# SSH Keeper v0.1.6 Release Notes

## 🎉 CI/CD Pipeline Release

This release introduces a fully automated CI/CD pipeline with GitHub Actions, making SSH Keeper development more robust and reliable.

## ✨ What's New

### CI/CD Pipeline

- 🚀 **Automated Builds** - GitHub Actions builds for all platforms automatically
- 🔄 **Continuous Integration** - Automated testing and validation on every commit
- 📦 **Automatic Releases** - Tag-based releases with artifacts
- 🛡️ **Security Scanning** - Gosec security scanner integration (temporarily disabled)
- ⚡ **Fast Feedback** - Quick build validation for developers

### Previous Features (v0.1.0)

- 🎨 **Beautiful TUI Interface** - Modern terminal user interface built with Bubble Tea
- 🔐 **Secure Password Storage** - Master password with system keyring integration
- 🔑 **Dual Authentication** - Support for both password and SSH key authentication
- 📁 **Connection Management** - Add, edit, delete, and organize SSH connections
- 🔍 **Smart Search** - Quick connection search and filtering
- 📤 **Export/Import** - Full compatibility with OpenSSH config format

## ✨ What's New

### Core Features

- 🎨 **Beautiful TUI Interface** - Modern terminal user interface built with Bubble Tea
- 🔐 **Secure Password Storage** - Master password with system keyring integration
- 🔑 **Dual Authentication** - Support for both password and SSH key authentication
- 📁 **Connection Management** - Add, edit, delete, and organize SSH connections
- 🔍 **Smart Search** - Quick connection search and filtering
- 📤 **Export/Import** - Full compatibility with OpenSSH config format

### Security Features

- **Encrypted Storage** - All connection data encrypted using AES-256
- **System Keyring Integration** - Master password stored in system keyring
- **Memory Management** - Sensitive data cleared from memory after use
- **No Plain Text** - No passwords stored in plain text files

### Cross-Platform Support

- ✅ **macOS** (Intel & Apple Silicon)
- ✅ **Linux** (x86_64)
- ✅ **Windows** (x86_64)

## 🚀 Installation

### One-Line Install (Recommended)

**macOS & Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.sh | bash
```

**Windows:**

```powershell
iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex
```

### Manual Download

Download the appropriate package for your platform:

- **macOS Intel**: `ssh-keeper-0.1.6-darwin-amd64.tar.gz`
- **macOS Apple Silicon**: `ssh-keeper-0.1.6-darwin-arm64.tar.gz`
- **Linux**: `ssh-keeper-0.1.6-linux-amd64.tar.gz`
- **Windows**: `ssh-keeper-0.1.6-windows-amd64.zip`

### Extract and Run

```bash
# Extract the archive
tar -xzf ssh-keeper-0.1.6-*.tar.gz  # Linux/macOS
# or unzip ssh-keeper-0.1.6-windows-amd64.zip  # Windows

# Make executable (Linux/macOS)
chmod +x ssh-keeper*

# Run
./ssh-keeper
```

## 🔧 Requirements

- **Go 1.24+** (for building from source)
- **Terminal** with color support
- **System Keyring** (Keychain/Secret Service/Credential Manager)

## 📖 Usage

1. **First Run**: Set up your master password
2. **Add Connections**: Use the intuitive form to add SSH connections
3. **Connect**: Select a connection and connect instantly
4. **Manage**: Edit, delete, or organize your connections
5. **Export/Import**: Compatible with OpenSSH config format

## 🛠️ Development

### CI/CD Pipeline

SSH Keeper now uses GitHub Actions for automated CI/CD:

- ✅ **Automated Testing** - Tests run on every push/PR
- ✅ **Multi-platform Builds** - Linux, macOS, Windows
- ✅ **Automatic Releases** - Tag-based releases with artifacts
- ✅ **Security Scanning** - Code and dependency vulnerability checks
- ✅ **Code Quality** - Linting and formatting checks

**Create a release:**

```bash
git tag v0.1.7
git push origin v0.1.7
```

### Building from Source

```bash
git clone https://github.com/ankogit/ssh_keeper.git
cd ssh_keeper
make build
```

### Available Make Targets

- `make build` - Build for current platform
- `make build-all` - Build for all platforms
- `make release` - Create release packages
- `make test` - Run tests
- `make install` - Install system-wide

## 🔒 Security

SSH Keeper follows security best practices:

- All data encrypted with AES-256
- Master password stored in system keyring
- No sensitive data in plain text
- Memory cleared after use
- Open source for security audit

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 🐛 Known Issues

- Screenshots in README are placeholders (will be added in next release)
- Some terminal themes may not display colors optimally
- Windows keyring integration requires additional testing

## 🔮 Roadmap

- [ ] Screenshots and better documentation
- [ ] Additional authentication methods
- [ ] Connection groups and tags
- [ ] Plugin system
- [ ] Configuration sync across devices
- [ ] Advanced search and filtering

## 🙏 Acknowledgments

Built with amazing open source projects:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [go-keyring](https://github.com/zalando/go-keyring) - Keyring access
- [Termenv](https://github.com/muesli/termenv) - Terminal environment

---

**Download now and start managing your SSH connections with style! 🚀**
