# SSH Keeper v0.1.0 Release Notes

## 🎉 Initial Release

This is the first public release of SSH Keeper, a beautiful and secure CLI tool for managing SSH connections with a modern TUI interface.

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

### Quick Install

Download the appropriate package for your platform:

- **macOS Intel**: `ssh-keeper-0.1.0-darwin-amd64.tar.gz`
- **macOS Apple Silicon**: `ssh-keeper-0.1.0-darwin-arm64.tar.gz`
- **Linux**: `ssh-keeper-0.1.0-linux-amd64.tar.gz`
- **Windows**: `ssh-keeper-0.1.0-windows-amd64.zip`

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
