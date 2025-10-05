# SSH Keeper 🔐

A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface.

## Configuration

SSH Keeper uses environment variables for configuration. Create a `.env` file in the project root:

```bash
# Copy the example configuration
cp env.example .env
```

### Environment Variables

| Variable          | Description                          | Default                | Required |
| ----------------- | ------------------------------------ | ---------------------- | -------- |
| `DEBUG`           | Enable debug mode                    | `false`                | No       |
| `ENV`             | Environment (development/production) | `development`          | No       |
| `CONFIG_PATH`     | Path to application config file      | `~/.ssh-keeper/config` | No       |
| `APP_SIGNATURE`   | Application signature for security   | -                      | Yes      |
| `SSH_CONFIG_PATH` | Path to SSH config file              | `~/.ssh/config`        | No       |
| `APP_NAME`        | Application name                     | `ssh-keeper`           | No       |
| `APP_VERSION`     | Application version                  | `1.0.0`                | No       |
| `LOG_LEVEL`       | Logging level                        | `info`                 | No       |
| `LOG_FORMAT`      | Logging format                       | `text`                 | No       |

### Security Configuration

The `APP_SIGNATURE` is used to protect the master password stored in the system keyring. In development mode, you can use `ssh-keeper-sig-dev`. In production, this should be a unique signature generated during the build process.

## Features

- 🔍 **Browse & Search** - View and search your SSH connections
- ➕ **Add Connections** - Easily add new SSH connections
- 🔐 **Secure Storage** - Passwords stored securely using go-keyring
- 📤 **Export/Import** - OpenSSH config compatibility
- 🎨 **Beautiful UI** - Modern TUI with colors and animations
- ⚡ **Fast & Lightweight** - Built with Go for performance

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/ssh-keeper.git
cd ssh-keeper

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

## Usage

```bash
# Start the application
ssh-keeper

# Run in development mode
make run-dev
```

### Main Menu

The application provides a beautiful main menu with the following options:

- **🔍 View Connections** - Browse and search your SSH connections
- **➕ Add Connection** - Add a new SSH connection
- **⚙️ Settings** - Configure application settings
- **📤 Export** - Export connections to OpenSSH config
- **📥 Import** - Import connections from OpenSSH config
- **❌ Quit** - Exit the application

### Keyboard Shortcuts

- `↑/↓` - Navigate menu items
- `Enter` - Select item
- `Ctrl+S` - Search connections
- `Esc` - Go back
- `Q` - Quit application

## Configuration

The application stores its configuration in `~/.ssh-keeper/`:

- `config.yaml` - Application settings
- `connections.yaml` - SSH connections (encrypted)
- Passwords are stored securely using go-keyring

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)

### Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/ssh-keeper.git
cd ssh-keeper

# Download dependencies
make deps

# Set up development environment
make dev-setup

# Run tests
make test

# Run with coverage
make test-coverage
```

### Project Structure

```
ssh_keeper/
├── cmd/ssh-keeper/          # Main application
├── internal/
│   ├── models/              # Data models
│   ├── ui/                  # TUI components
│   ├── storage/             # Data storage
│   ├── ssh/                 # SSH integration
│   └── auth/                # Authentication
├── pkg/                     # Public packages
├── go.mod                   # Go modules
├── Makefile                 # Build automation
└── README.md               # This file
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

## Security

- Passwords are stored securely using go-keyring
- Master key is stored in memory with timeout
- No sensitive data is stored in plain text files
- Full compatibility with OpenSSH security standards

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [go-keyring](https://github.com/99designs/go-keyring) - Secure storage
- [Termenv](https://github.com/muesli/termenv) - Terminal environment

#возможности ssh подключения

1. передача пароля через sshpass
2. SSH_ASKPASS
3. github.com/Netflix/go-expect
4. golang.org/x/crypto/ssh
