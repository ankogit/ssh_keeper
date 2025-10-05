# SSH Keeper - Windows One-line installer
# Usage: iwr -useb https://raw.githubusercontent.com/ankogit/ssh_keeper/main/scripts/install.ps1 | iex

param(
    [string]$Version = "v0.1.0",
    [string]$InstallDir = "$env:USERPROFILE\bin"
)

# Configuration
$REPO = "ankogit/ssh_keeper"
$BINARY_NAME = "ssh-keeper"

# Colors for output
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color = "White"
    )
    Write-Host $Message -ForegroundColor $Color
}

# Detect architecture
function Get-Architecture {
    if ([Environment]::Is64BitOperatingSystem) {
        return "amd64"
    } else {
        return "386"
    }
}

# Download and install
function Install-SSHKeeper {
    $arch = Get-Architecture
    $platform = "windows-$arch"
    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/ssh-keeper-$Version-$platform.zip"
    $filename = "ssh-keeper-$Version-$platform.zip"
    $binaryName = "ssh-keeper-$Version-$platform.exe"
    
    Write-ColorOutput "🚀 Installing SSH Keeper..." "Cyan"
    Write-ColorOutput "Platform: $platform" "Yellow"
    Write-ColorOutput "Version: $Version" "Yellow"
    Write-ColorOutput "Download URL: $downloadUrl" "Yellow"
    
    # Create temporary directory
    $tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    try {
        # Download
        Write-ColorOutput "📥 Downloading SSH Keeper..." "Cyan"
        $zipPath = Join-Path $tempDir $filename
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath -UseBasicParsing
        
        # Extract
        Write-ColorOutput "📦 Extracting archive..." "Cyan"
        $extractPath = Join-Path $tempDir "extracted"
        Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force
        
        # Install
        Write-ColorOutput "🔧 Installing SSH Keeper..." "Cyan"
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }
        
        $sourcePath = Join-Path $extractPath $binaryName
        $destPath = Join-Path $InstallDir "$BINARY_NAME.exe"
        
        Copy-Item -Path $sourcePath -Destination $destPath -Force
        
        Write-ColorOutput "✅ SSH Keeper installed to $destPath" "Green"
        
        # Add to PATH if not already there
        $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($currentPath -notlike "*$InstallDir*") {
            Write-ColorOutput "🔧 Adding to PATH..." "Cyan"
            [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$InstallDir", "User")
            Write-ColorOutput "✅ Added $InstallDir to PATH" "Green"
            Write-ColorOutput "💡 Please restart your terminal for PATH changes to take effect" "Yellow"
        }
        
    } finally {
        # Cleanup
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

# Verify installation
function Test-Installation {
    Write-ColorOutput "🔍 Verifying installation..." "Cyan"
    
    $binaryPath = Join-Path $InstallDir "$BINARY_NAME.exe"
    if (Test-Path $binaryPath) {
        Write-ColorOutput "✅ SSH Keeper is installed successfully!" "Green"
        Write-ColorOutput "   Location: $binaryPath" "Green"
        Write-ColorOutput "" "White"
        Write-ColorOutput "🎉 Ready to use! Run '$BINARY_NAME' to start." "Cyan"
        return $true
    } else {
        Write-ColorOutput "❌ Installation verification failed" "Red"
        Write-ColorOutput "💡 Try running '$BINARY_NAME' manually" "Yellow"
        return $false
    }
}

# Show usage information
function Show-Usage {
    Write-ColorOutput "📖 SSH Keeper Usage:" "Cyan"
    Write-ColorOutput "   $BINARY_NAME                    # Start SSH Keeper" "Yellow"
    Write-ColorOutput "   $BINARY_NAME --version          # Show version" "Yellow"
    Write-ColorOutput "   $BINARY_NAME --help             # Show help" "Yellow"
    Write-ColorOutput "" "White"
    Write-ColorOutput "📚 Documentation:" "Cyan"
    Write-ColorOutput "   https://github.com/$REPO" "Yellow"
    Write-ColorOutput "   https://github.com/$REPO/releases" "Yellow"
}

# Main installation process
function Main {
    Write-ColorOutput "" "White"
    Write-ColorOutput "  ____  ____  _  _    _  __  __  _____  _____  _____  _____  _____  _____ " "Cyan"
    Write-ColorOutput " / ___||  _ \| || |  | |/ /|  \/  | __||  _  ||  _  ||  _  ||  _  ||  _  |" "Cyan"
    Write-ColorOutput " \___ \| |_) | || |_ | ' / | |\/| | _| | | | || | | || | | || | | || | | |" "Cyan"
    Write-ColorOutput "  ___) |  __/|__   _|| . \ | |  | | |_ | |_| || |_| || |_| || |_| || |_| |" "Cyan"
    Write-ColorOutput " |____/|_|      |_|  |_|\_\|_|  |_|\__| \___/  \___/  \___/  \___/  \___/ " "Cyan"
    Write-ColorOutput "" "White"
    Write-ColorOutput "🔐 Secure SSH Connection Manager" "Cyan"
    Write-ColorOutput "" "White"
    
    # Check if already installed
    $binaryPath = Join-Path $InstallDir "$BINARY_NAME.exe"
    if (Test-Path $binaryPath) {
        Write-ColorOutput "⚠️  SSH Keeper is already installed at: $binaryPath" "Yellow"
        $response = Read-Host "Do you want to reinstall? (y/N)"
        if ($response -notmatch "^[Yy]$") {
            Write-ColorOutput "Installation cancelled." "Cyan"
            return
        }
    }
    
    # Install
    Install-SSHKeeper
    if (Test-Installation) {
        Show-Usage
    }
}

# Run main function
Main
