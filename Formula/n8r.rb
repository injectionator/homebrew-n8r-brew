class N8r < Formula
  desc "The n8r CLI - Injectionator command-line tool"
  homepage "https://injectionator.com"
  url "https://github.com/injectionator/homebrew-n8r-brew/releases/download/v0.1.0/n8r-0.1.0.tar.gz"
  sha256 "6f35c0450232838556b4a42c82441db896517375b2e8f22bd24ccc6b14737c19"
  license "Copyright 2026 Steve Chambers, Injectionator, Viewyonder"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/injectionator/n8r/internal/config.Version=#{version}"), "./cmd/n8r"
  end

  test do
    system "#{bin}/n8r", "--version"
  end
end
