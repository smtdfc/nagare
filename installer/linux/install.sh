#!/bin/bash

# Check for root/sudo privileges
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script with sudo privileges (sudo bash install.sh)"
  exit 1
fi

REPO="smtdfc/nagare"
INSTALL_PATH="/usr/local/bin"

# Detect Operating System
OS="$(uname | tr '[:upper:]' '[:lower:]')"
if [ "$OS" != "linux" ] && [ "$OS" != "darwin" ]; then
  echo "Unsupported operating system: $OS"
  exit 1
fi

# Detect CPU Architecture
ARCH="$(uname -m)"
if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "Unsupported CPU architecture: $ARCH"
  exit 1
fi

echo "Detected system: OS=$OS, Arch=$ARCH"

# Fetch the latest version from GitHub Releases
echo "Checking for the latest version from GitHub..."
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
  echo "Failed to fetch version from GitHub. Please check your network connection."
  exit 1
fi

echo "Selected version for installation: $VERSION"

# Strip leading 'v' to match the asset naming convention if necessary
CLEAN_VERSION="${VERSION#v}"

# Define binary filenames based on semantic-release configuration
CLI_FILE="nagare-${OS}-${ARCH}-${CLEAN_VERSION}"
GATEWAY_FILE="nagare-gateway-${OS}-${ARCH}-${CLEAN_VERSION}"

BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

echo "Downloading binaries..."

# Download CLI
CLI_URL="$BASE_URL/$CLI_FILE"
echo "Downloading CLI from: $CLI_URL"
if curl -sL -f "$CLI_URL" -o "$INSTALL_PATH/nagare"; then
  chmod +x "$INSTALL_PATH/nagare"
  echo "Successfully installed CLI: nagare"
else
  echo "Warning: Could not download CLI asset for this configuration from GitHub Release."
fi

# Download Gateway
GATEWAY_URL="$BASE_URL/$GATEWAY_FILE"
echo "Downloading Gateway from: $GATEWAY_URL"
if curl -sL -f "$GATEWAY_URL" -o "$INSTALL_PATH/nagare-gateway"; then
  chmod +x "$INSTALL_PATH/nagare-gateway"
  echo "Successfully installed Gateway: nagare-gateway"
else
  echo "Warning: Could not download Gateway asset for this configuration from GitHub Release."
fi

echo "========================================="
echo " Nagare installation completed ($VERSION)!"
echo "========================================="