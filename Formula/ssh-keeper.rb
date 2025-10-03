class SshKeeper < Formula
  desc "A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface"
  homepage "https://github.com/yourusername/ssh-keeper"
  url "https://github.com/yourusername/ssh-keeper/archive/v0.1.0.tar.gz"
  sha256 "" # This will be filled when you create the actual release
  license "MIT"
  head "https://github.com/yourusername/ssh-keeper.git", branch: "main"

  depends_on "go" => :build

  def install
    # Set up Go environment
    ENV["GOPATH"] = buildpath
    ENV["GO111MODULE"] = "on"
    
    # Build the application
    system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "ssh-keeper", "./cmd/ssh-keeper"
    bin.install "ssh-keeper"
  end

  test do
    system "#{bin}/ssh-keeper", "--version"
  end
end
