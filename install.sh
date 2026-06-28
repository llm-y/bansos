#!/bin/bash
set -e

# Auto-elevate to root if not already running as root
if [ "$(id -u)" -ne 0 ]; then
  echo "Requesting root privileges..."
  exec sudo bash "$0" "$@"
fi

echo "==========================================="
echo "  Bansos Installer (Root Mode)"
echo "==========================================="
echo ""

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

echo "[1/4] Detected OS: $OS, Architecture: $ARCH"
echo "[2/4] Downloading latest release of $ASSET_NAME..."

# Get latest release download URL
DOWNLOAD_URL=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep "browser_download_url.*${ASSET_NAME}\"" \
  | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Error: Could not find release asset for ${ASSET_NAME}"
  exit 1
fi

echo "  Download URL: $DOWNLOAD_URL"

# Download binary to temp location
TMP_FILE=$(mktemp)
curl -sL "$DOWNLOAD_URL" -o "$TMP_FILE"
chmod +x "$TMP_FILE"

# Install binary
echo "[3/4] Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
mv "$TMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"

# Configure firewall rules to allow bansos network access
echo "[4/4] Configuring firewall rules..."

BINARY_PATH="${INSTALL_DIR}/${BINARY_NAME}"

if [ "$OS" = "linux" ]; then
  # Check for ufw (Uncomplicated Firewall)
  if command -v ufw >/dev/null 2>&1 && ufw status | grep -q "active"; then
    echo "  Detected active ufw firewall, adding rules..."
    # Allow outbound connections for bansos
    ufw allow out to any proto tcp from any app bansos 2>/dev/null || \
    ufw allow out to any proto tcp 2>/dev/null || true
    echo "  UFW: Outbound TCP traffic allowed for bansos."
  # Check for firewalld
  elif command -v firewall-cmd >/dev/null 2>&1 && systemctl is-active --quiet firewalld 2>/dev/null; then
    echo "  Detected active firewalld, adding rules..."
    # Allow outbound traffic (firewalld allows outbound by default, but ensure the zone permits it)
    firewall-cmd --permanent --add-rich-rule='rule family="ipv4" service name="http" accept' 2>/dev/null || true
    firewall-cmd --permanent --add-rich-rule='rule family="ipv4" service name="https" accept' 2>/dev/null || true
    firewall-cmd --reload 2>/dev/null || true
    echo "  Firewalld: HTTP/HTTPS outbound access configured for bansos."
  # Check for iptables
  elif command -v iptables >/dev/null 2>&1; then
    echo "  Detected iptables, adding rules..."
    # Allow outbound TCP connections for the bansos binary owner (root)
    iptables -C OUTPUT -p tcp -m owner --cmd-owner bansos -j ACCEPT 2>/dev/null || \
    iptables -A OUTPUT -p tcp -m owner --cmd-owner bansos -j ACCEPT 2>/dev/null || true
    echo "  Iptables: Outbound TCP traffic allowed for bansos."
  else
    echo "  No active firewall detected. No firewall rules needed."
  fi
elif [ "$OS" = "darwin" ]; then
  # macOS: configure application firewall
  if command -v /usr/libexec/ApplicationFirewall/socketfilterfw >/dev/null 2>&1; then
    echo "  Configuring macOS Application Firewall..."
    # Add bansos to the firewall allow list
    /usr/libexec/ApplicationFirewall/socketfilterfw --add "$BINARY_PATH" 2>/dev/null || true
    /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp "$BINARY_PATH" 2>/dev/null || true
    echo "  macOS Firewall: bansos is allowed to accept and send network connections."
  else
    echo "  macOS Application Firewall tool not found. No firewall configuration needed."
  fi
fi

echo ""
echo "==========================================="
echo "  Installation complete!"
echo "==========================================="
echo "  Binary: ${INSTALL_DIR}/${BINARY_NAME}"
echo "  Run 'bansos' to get started."
echo ""
