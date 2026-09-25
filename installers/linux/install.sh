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
    echo " Nagare installation completed!"
    echo "========================================="
    exit 0
fi


if [ "$EUID" -ne 0 ]; then
    echo "Please run this script with sudo privileges (sudo bash install.sh) to download and install from GitHub."
    exit 1
fi

VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to fetch version from GitHub. Please check your network connection."
    exit 1
fi

echo "Selected version for installation: $VERSION"

ARCHIVE_NAME="nagare-${OS}-${ARCH}.tar.gz"
BASE_URL="https://github.com/$REPO/releases/download/$VERSION"
ARCHIVE_URL="$BASE_URL/$ARCHIVE_NAME"

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading archive from: $ARCHIVE_URL"
if curl -sL -f "$ARCHIVE_URL" -o "$TMP_DIR/$ARCHIVE_NAME"; then
    echo "Extracting archive..."
    tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"

    if [ -f "$TMP_DIR/nagare" ]; then
        install_binary "$TMP_DIR/nagare" "$INSTALL_PATH/nagare" "CLI"
    else
        echo "Warning: 'nagare' binary not found in the downloaded archive."
    fi

    if [ -f "$TMP_DIR/nagare-gateway" ]; then
        install_binary "$TMP_DIR/nagare-gateway" "$INSTALL_PATH/nagare-gateway" "Gateway"
    else
        echo "Notice: 'nagare-gateway' binary not found in the archive, skipping."
    fi
else
    echo "Error: Could not download release archive for this configuration from GitHub."
    exit 1
fi

echo "========================================="
echo " Nagare installation completed!"
echo "========================================="