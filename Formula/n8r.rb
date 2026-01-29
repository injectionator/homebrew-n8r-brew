class N8r < Formula
  desc "The n8r CLI - Injectionator command-line tool"
  homepage "https://injectionator.com"
  url "https://github.com/injectionator/n8r-brew/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "TODO"
  license "Copyright 2026 Steve Chambers, Injectionator, Viewyonder"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/injectionator/n8r/internal/config.Version=#{version}"), "./cmd/n8r"
  end

  test do
    system "#{bin}/n8r", "--version"
  end
end
