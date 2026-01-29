# typed: false
# frozen_string_literal: true

# This file is auto-updated by goreleaser on each release.
# Manual changes will be overwritten.

class Notte < Formula
  desc "Browser automation CLI for notte.cc"
  homepage "https://github.com/nottelabs/notte-cli"
  license "MIT"
  version "0.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/nottelabs/notte-cli/releases/download/v#{version}/notte-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    else
      url "https://github.com/nottelabs/notte-cli/releases/download/v#{version}/notte-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/nottelabs/notte-cli/releases/download/v#{version}/notte-cli_#{version}_linux_arm64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    else
      url "https://github.com/nottelabs/notte-cli/releases/download/v#{version}/notte-cli_#{version}_linux_amd64.tar.gz"
      sha256 "SHA256_PLACEHOLDER"
    end
  end

  def install
    bin.install "notte"
  end

  test do
    system "#{bin}/notte", "version"
  end
end
