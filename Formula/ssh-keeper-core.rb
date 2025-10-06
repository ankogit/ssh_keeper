class SshKeeper < Formula
  desc "A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface"
  homepage "https://github.com/yourusername/ssh-keeper"
  url "https://github.com/yourusername/ssh-keeper/archive/v0.1.0.tar.gz"
  sha256 "abc123..."  # SHA256 архива с исходным кодом
  license "MIT"
  head "https://github.com/yourusername/ssh-keeper.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-X main.version=#{version}", "-o", "ssh-keeper", "./cmd/ssh-keeper"
    bin.install "ssh-keeper"
  end

  test do
    assert_match "SSH Keeper", shell_output("#{bin}/ssh-keeper --help", 1)
  end
end




