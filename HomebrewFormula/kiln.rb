# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kiln < Formula
  desc ""
  homepage ""
  version "0.80.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/kiln/releases/download/v0.80.0/kiln-darwin-amd64-0.80.0.tar.gz"
      sha256 "99f9e3f6bdf9b83da7c5361aa9eb55dd25d7507f76980fab4616b3c9c9a0269e"

      def install
        bin.install "kiln"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/pivotal-cf/kiln/releases/download/v0.80.0/kiln-darwin-arm64-0.80.0.tar.gz"
      sha256 "d498624238d092a050c26f6741899590ea8d4397e7a161635200492ebaf39e0d"

      def install
        bin.install "kiln"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/kiln/releases/download/v0.80.0/kiln-linux-amd64-0.80.0.tar.gz"
      sha256 "5e9bc76f043f891382a2bba765296d13a736ae3c8f1d2db429cb7c0938b8174c"

      def install
        bin.install "kiln"
      end
    end
  end

  test do
    system "#{bin}/kiln --version"
  end
end
