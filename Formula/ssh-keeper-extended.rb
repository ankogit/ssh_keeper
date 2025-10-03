class SshKeeper < Formula
  desc "A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface"
  homepage "https://github.com/yourusername/ssh-keeper"
  url "https://github.com/yourusername/ssh-keeper/archive/v0.1.0.tar.gz"
  sha256 "" # This will be filled when you create the actual release
  license "MIT"
  head "https://github.com/yourusername/ssh-keeper.git", branch: "main"

  # Build dependencies (needed only for building)
  depends_on "go" => :build

  # Runtime dependencies (needed for the app to work)
  # Note: Go applications are statically linked, so most dependencies are bundled
  # But we might need some system libraries for certain features
  
  # Optional dependencies for enhanced functionality
  depends_on "sshpass" => :optional  # For password-based SSH connections
  depends_on "gnupg" => :optional   # For GPG key management

  # Platform-specific dependencies
  depends_on "pkg-config" if OS.linux?

  def install
    # Set up Go environment
    ENV["GOPATH"] = buildpath
    ENV["GO111MODULE"] = "on"
    ENV["CGO_ENABLED"] = "0"  # Disable CGO for static linking
    
    # Build the application
    system "go", "build", 
           "-ldflags", "-X main.version=#{version} -s -w",  # Strip debug info
           "-o", "ssh-keeper", 
           "./cmd/ssh-keeper"
    
    # Install binary
    bin.install "ssh-keeper"
    
    # Install man page (if you have one)
    # man1.install "ssh-keeper.1"
    
    # Install bash completion (if you have one)
    # bash_completion.install "completions/bash/ssh-keeper"
    
    # Install zsh completion (if you have one)
    # zsh_completion.install "completions/zsh/_ssh-keeper"
  end

  test do
    # Test that the binary works
    assert_match "SSH Keeper", shell_output("#{bin}/ssh-keeper --help", 1)
  end

  # Caveats shown to user after installation
  def caveats
    <<~EOS
      SSH Keeper has been installed!
      
      To get started:
        1. Run: ssh-keeper
        2. Follow the setup wizard
      
      Configuration files will be stored in:
        ~/.ssh-keeper/
      
      Optional features:
        - Install sshpass for password-based connections: brew install sshpass
        - Install gnupg for GPG key management: brew install gnupg
      
      For more information, visit: #{homepage}
    EOS
  end
end

