#!/bin/bash

REPO="smtdfc/nagare"
INSTALL_PATH="/usr/local/bin"
install_binary() {
    local src="$1"
    local dest="$2"
    local name="$3"

    if cp "$src" "$dest"; then
        chmod +x "$dest"
        echo "Successfully installed $name: $(basename "$dest")"
        return 0
    else
        echo "Error: Failed to install $name to $dest."
        return 1
    fi
}

detect_os() {
    local os
    os="$(uname | tr '[:upper:]' '[:lower:]')"
    if [ "$os" != "linux" ] && [ "$os" != "darwin" ]; then
        echo "Unsupported operating system: $os"
        exit 1
    fi
    echo "$os"
}

detect_arch() {
    local arch
    arch="$(uname -m)"
    if [ "$arch" = "x86_64" ]; then
        echo "amd64"
    elif [ "$arch" = "aarch64" ] || [ "$arch" = "arm64" ]; then
        echo "arm64"
    else
        echo "Unsupported CPU architecture: $arch"
        exit 1
    fi
}


OS="$(detect_os)"
ARCH="$(detect_arch)"
echo "Detected system: OS=$OS, Arch=$ARCH"

LOCAL_DIST="dist/${OS}-${ARCH}"

if [ -d "$LOCAL_DIST" ] && [ -f "$LOCAL_DIST/nagare" ]; then
    echo "Found local build in $LOCAL_DIST. Installing from local files..."

    install_binary "$LOCAL_DIST/nagare" "$INSTALL_PATH/nagare" "CLI (local)" || exit 1
    if [ -f "$LOCAL_DIST/nagare-gateway" ]; then
        install_binary "$LOCAL_DIST/nagare-gateway" "$INSTALL_PATH/nagare-gateway" "Gateway (local)"
    else
        echo "Notice: Local Gateway binary not found in $LOCAL_DIST, skipping."
    fi

    echo "========================================="
    echo " Nagare local installation completed!"
    echo "========================================="
    exit 0
fi


if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo privileges (sudo bash install.sh) to download and install from GitHub."
    exit 1
fi

echo "Local build not found in $LOCAL_DIST. Fetching the latest version from GitHub..."
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to fetch version from GitHub. Please check your network connection."
    exit 1
fi

echo "Selected version for installation: $VERSION"

CLEAN_VERSION="${VERSION#v}"
CLI_FILE="nagare-${OS}-${ARCH}-${CLEAN_VERSION}"
GATEWAY_FILE="nagare-gateway-${OS}-${ARCH}-${CLEAN_VERSION}"
BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

echo "Downloading binaries from GitHub..."
CLI_URL="$BASE_URL/$CLI_FILE"
echo "Downloading CLI from: $CLI_URL"
if curl -sL -f "$CLI_URL" -o "$INSTALL_PATH/nagare"; then
    chmod +x "$INSTALL_PATH/nagare"
    echo "Successfully installed CLI: nagare"
else
    echo "Warning: Could not download CLI asset for this configuration from GitHub Release."
fi

GATEWAY_URL="$BASE_URL/$GATEWAY_FILE"
echo "Downloading Gateway from: $GATEWAY_URL"
if curl -sL -f "$GATEWAY_URL" -o "$INSTALL_PATH/nagare-gateway"; then
    chmod +x "$INSTALL_PATH/nagare-gateway"
    echo "Successfully installed Gateway: nagare-gateway"
else
    echo "Warning: Could not download Gateway asset for this configuration from GitHub Release."
fi

echo "========================================="
echo " Nagare GitHub installation completed ($VERSION)!"
echo "========================================="