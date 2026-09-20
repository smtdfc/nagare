#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "Please run this script with sudo privileges (sudo bash install.sh)."
  exit 1
fi

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
else
  echo "This architecture is not supported: $ARCH"
  exit 1
fi

TARGET_DIR="dist/${OS}-${ARCH}"

if [ ! -d "$TARGET_DIR" ]; then
  echo "No asset: $TARGET_DIR"
  exit 1
fi

INSTALL_PATH="/usr/local/bin"

if [ -f "$TARGET_DIR/nagare" ]; then
  cp "$TARGET_DIR/nagare" "$INSTALL_PATH/"
  chmod +x "$INSTALL_PATH/nagare"
fi

if [ -f "$TARGET_DIR/nagare-gateway" ]; then
  cp "$TARGET_DIR/nagare-gateway" "$INSTALL_PATH/"
  chmod +x "$INSTALL_PATH/nagare-gateway"
fi

echo "Nagare Installed!"