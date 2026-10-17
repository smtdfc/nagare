#!/bin/bash

REPO="smtdfc/nagare"
INSTALL_PATH="$HOME/.local/bin"

NAGARE_DATA_DIR="$HOME/.nagare"
NAGARE_ASSETS_DIR="$NAGARE_DATA_DIR/assets"
BUILTIN_PLUGINS_STORE_DIR="$NAGARE_ASSETS_DIR/plugins/builtins"
CORE_PLUGINS_STORE_DIR="$NAGARE_ASSETS_DIR/plugins/core"
WEB_STORE_DIR="$NAGARE_ASSETS_DIR/web"

mkdir -p "$NAGARE_ASSETS_DIR"
mkdir -p "$INSTALL_PATH"
mkdir -p "$BUILTIN_PLUGINS_STORE_DIR"
mkdir -p "$CORE_PLUGINS_STORE_DIR"
mkdir -p "$WEB_STORE_DIR"

install_binary() {
    local src="$1"
    local dest="$2"
    local name="$3"

    rm -f "$dest"

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

LOCAL_DIST="dist"

if [ -d "$LOCAL_DIST" ] && [ -f "$LOCAL_DIST/nagare" ]; then
    echo "Found local build in $LOCAL_DIST. Installing from local files..."

    install_binary "$LOCAL_DIST/nagare" "$INSTALL_PATH/nagare" "CLI (local)" || exit 1

    if [ -f "$LOCAL_DIST/nagare-gateway" ]; then
        install_binary "$LOCAL_DIST/nagare-gateway" "$INSTALL_PATH/nagare-gateway" "Gateway (local)"
    else
        echo "Notice: Local Gateway binary not found in $LOCAL_DIST, skipping."
    fi

    if [ -f "$LOCAL_DIST/nagare-desktop" ]; then
        install_binary "$LOCAL_DIST/nagare-desktop" "$INSTALL_PATH/nagare-desktop" "Nagare Desktop (local)"
    else
        echo "Notice: Local Nagare desktop binary not found in $LOCAL_DIST, skipping."
    fi

    echo "Copying plugins and web assets..."
    [ -d "$LOCAL_DIST/plugins/builtins" ] && cp -r "$LOCAL_DIST/plugins/builtins/." "$BUILTIN_PLUGINS_STORE_DIR/"
    [ -d "$LOCAL_DIST/plugins/core" ] && cp -r "$LOCAL_DIST/plugins/core/." "$CORE_PLUGINS_STORE_DIR/"
    [ -d "$LOCAL_DIST/web" ] && cp -r "$LOCAL_DIST/web/." "$WEB_STORE_DIR/"

    echo "========================================="
    echo " Nagare local installation completed!"
    echo "========================================="
    exit 0
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

    if [ -f "$TMP_DIR/nagare-desktop" ]; then
        install_binary "$TMP_DIR/nagare-desktop" "$INSTALL_PATH/nagare-desktop" "Nagare Desktop"
    else
        echo "Notice: 'nagare-desktop' binary not found in the archive, skipping."
    fi

    echo "Copying plugins and web assets..."
    [ -d "$TMP_DIR/plugins/builtins" ] && cp -r "$TMP_DIR/plugins/builtins/." "$BUILTIN_PLUGINS_STORE_DIR/"
    [ -d "$TMP_DIR/plugins/core" ] && cp -r "$TMP_DIR/plugins/core/." "$CORE_PLUGINS_STORE_DIR/"
    [ -d "$TMP_DIR/plugins/web" ] && cp -r "$TMP_DIR/plugins/web/." "$WEB_STORE_DIR/"
else
    echo "Error: Could not download release archive for this configuration from GitHub."
    exit 1
fi

echo "========================================="
echo " Nagare installation completed!"
echo "========================================="
