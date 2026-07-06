#!/bin/sh
set -e

REPO="Foxemsx/speed"
BINARY="speed"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  linux*)  OS="linux" ;;
  darwin*) OS="darwin" ;;
  *)
    echo "Error: unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)   ARCH="amd64" ;;
  aarch64|arm64)   ARCH="arm64" ;;
  armv7l|armhf)    ARCH="armv7" ;;
  *)
    echo "Error: unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Detect libc (for Linux)
LIBC="musl"
if [ "$OS" = "linux" ]; then
  if ldd --version 2>&1 | grep -q glibc; then
    LIBC="libc"
  fi
fi

ASSET="${BINARY}_${OS}_${ARCH}"
if [ "$OS" = "linux" ]; then
  ASSET="${BINARY}_${OS}_${LIBC}_${ARCH}"
fi

# Fetch latest release tag
echo "Fetching latest release..."
TAG=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$TAG" ]; then
  echo "Error: could not determine latest release."
  echo "You can install with: go install github.com/${REPO}@latest"
  exit 1
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET}"

echo "Downloading ${ASSET} ${TAG}..."
TMPFILE=$(mktemp)
curl -sL "$DOWNLOAD_URL" -o "$TMPFILE"

if [ ! -s "$TMPFILE" ]; then
  echo "Error: download failed. The release may not include this platform."
  echo "You can install with: go install github.com/${REPO}@latest"
  rm -f "$TMPFILE"
  exit 1
fi

chmod +x "$TMPFILE"

# Install
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMPFILE" "${INSTALL_DIR}/${BINARY}"
else
  echo "Installing to ${INSTALL_DIR} (requires sudo)..."
  sudo mv "$TMPFILE" "${INSTALL_DIR}/${BINARY}"
fi

echo "Installed ${BINARY} ${TAG} to ${INSTALL_DIR}/${BINARY}"
echo "Run '${BINARY}' to start a speed test."
