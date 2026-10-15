#!/usr/bin/env bash

set -e

VERSION="${1:-dev}"

rm -rf dist
mkdir -p dist

echo "Building Rune ${VERSION}..."

# Build UI
echo "Building web UI..."
npm --prefix ui install
npm --prefix ui run build

# Rune CLI
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune_darwin_arm64        ./cmd/cli
GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune_darwin_amd64        ./cmd/cli
GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune_linux_arm64         ./cmd/cli
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune_linux_amd64         ./cmd/cli
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune_windows_amd64.exe   ./cmd/cli

# Rune Server
GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune-server_darwin_arm64       ./cmd/server
GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune-server_darwin_amd64       ./cmd/server
GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune-server_linux_arm64        ./cmd/server
GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune-server_linux_amd64        ./cmd/server
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o dist/rune-server_windows_amd64.exe  ./cmd/server

chmod +x dist/*

echo ""
echo "Build complete:"
ls -lh dist/
