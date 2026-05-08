#!/usr/bin/env bash

set -e

VERSION="v0.2.0-dev"
REPO="https://github.com/hyphen-ani/rune/releases/download"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
fi

RUNE_BINARY="rune-${OS}-${ARCH}"
SERVER_BINARY="rune-server-${OS}-${ARCH}"

TMP_DIR=$(mktemp -d)

echo "Detected OS: $OS"
echo "Detected ARCH: $ARCH"

echo "Downloading Rune..."

curl -fL "$REPO/$VERSION/$RUNE_BINARY" -o "$TMP_DIR/rune"
curl -fL "$REPO/$VERSION/$SERVER_BINARY" -o "$TMP_DIR/rune-server"

chmod +x "$TMP_DIR/rune"
chmod +x "$TMP_DIR/rune-server"

echo "Installing Rune..."

sudo mv "$TMP_DIR/rune" /usr/local/bin/rune
sudo mv "$TMP_DIR/rune-server" /usr/local/bin/rune-server

rm -rf "$TMP_DIR"

echo
echo "Rune installed successfully"
echo
echo "Run:"
echo "  rune --version"
echo "  rune-server"
echo