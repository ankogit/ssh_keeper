class SshKeeper < Formula
  desc "A beautiful and secure CLI tool for managing SSH connections with a modern TUI interface"
  homepage "https://github.com/ankogit/ssh_keeper"
  url "https://github.com/ankogit/ssh_keeper/archive/v0.1.0.tar.gz"
  sha256 "YOUR_SHA256_HERE" # This needs to be calculated
  license "MIT"
  head "https://github.com/ankogit/ssh_keeper.git", branch: "main"

  depends_on "go" => :build

  def install
    system "make", "build"
    bin.install "build/ssh-keeper"
  end

  test do
    system "#{bin}/ssh-keeper", "--version"
  end
end
