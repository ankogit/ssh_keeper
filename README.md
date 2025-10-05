# SSH Keeper

SSH Connection Manager with secure password storage and modern TUI interface.

## Quick Start

```bash
# Build
go build -o ssh-keeper ./cmd/ssh-keeper/

# Run
./ssh-keeper

# Or install using the install script
./scripts/install.sh
```

## Configuration

Create a `.env` file in the project root:

```bash
cp env.example .env
```

Set your `APP_SIGNATURE` in the `.env` file for development.

## Documentation

All documentation is located in the [`docs/`](./docs/) directory:

- **[README.md](./docs/README.md)** - Complete project documentation
- **[INSTALL.md](./docs/INSTALL.md)** - Installation instructions
- **[CONFIG_DOCUMENTATION.md](./docs/CONFIG_DOCUMENTATION.md)** - Configuration guide
- **[SSH_KEYS_DOCUMENTATION.md](./docs/SSH_KEYS_DOCUMENTATION.md)** - SSH keys setup
- **[PROJECT_DESIGN.md](./docs/PROJECT_DESIGN.md)** - Architecture overview
- **[SSH_CLIENTS_ARCHITECTURE.md](./docs/SSH_CLIENTS_ARCHITECTURE.md)** - SSH client architecture
- **[TESTING_SSH.md](./docs/TESTING_SSH.md)** - Testing guide
- **[PACKAGING_SETUP.md](./docs/PACKAGING_SETUP.md)** - Packaging instructions
- **[HOMEBREW_CORE_CHECKLIST.md](./docs/HOMEBREW_CORE_CHECKLIST.md)** - Homebrew setup
- **[DEPENDENCIES.md](./docs/DEPENDENCIES.md)** - Dependencies overview
- **[PROJECT_FEEDBACK.md](./docs/PROJECT_FEEDBACK.md)** - Project feedback

## Scripts

Utility scripts are located in the [`scripts/`](./scripts/) directory:

- **`install.sh`** - Install SSH Keeper system-wide
- **`build-deb.sh`** - Build Debian package
- **`check-formula.sh`** - Check Homebrew formula
- **`create-homebrew-tap.sh`** - Create Homebrew tap
- **`prepare-homebrew-core-pr.sh`** - Prepare Homebrew core PR
- **`update-homebrew-formula.sh`** - Update Homebrew formula
- **`test_navigation.sh`** - Test UI navigation

## Features

- 🔐 **Secure Password Storage** - Master password with system keyring integration
- 🎨 **Modern TUI** - Beautiful terminal interface with Bubble Tea
- 🔑 **SSH Key Support** - Both password and SSH key authentication
- 📁 **Connection Management** - Add, edit, delete, and organize connections
- 🔍 **Search & Filter** - Quick connection search
- ⚙️ **Configurable** - Flexible configuration via environment variables
- 🚀 **Fast & Lightweight** - Built with Go for performance

## License

MIT License - see [LICENSE](LICENSE) file for details.
