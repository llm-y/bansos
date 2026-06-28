#!/bin/bash
set -e

REPO="llm-y/bansos"
BINARY_NAME="bansos"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  amd64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

ASSET_NAME="bansos-${OS}-${ARCH}"

echo "Detected OS: $OS, Architecture: $ARCH"
echo "Downloading latest release of $ASSET_NAME..."

# Get latest release download URL
DOWNLOAD_URL=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep "browser_download_url.*${ASSET_NAME}\"" \
  | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Error: Could not find release asset for ${ASSET_NAME}"
  exit 1
fi

echo "Download URL: $DOWNLOAD_URL"

# Download binary to temp location
TMP_FILE=$(mktemp)
curl -sL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

# Install binary
echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
else
  sudo mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
fi

echo "Successfully installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
echo "Run 'bansos' to get started."
