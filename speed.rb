class Speed < Formula
  desc "Internet speed test in your terminal — download & upload"
  homepage "https://github.com/Foxemsx/speed"
  url "https://github.com/Foxemsx/speed/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "PLACEHOLDER"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/speed --version 2>&1", 1)
  end
end
