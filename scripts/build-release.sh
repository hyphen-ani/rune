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
GOOS=darwin GOARCH=arm64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-darwin-arm64 ./cmd/cli

GOOS=darwin GOARCH=amd64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-darwin-amd64 ./cmd/cli

GOOS=linux GOARCH=arm64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-linux-arm64 ./cmd/cli

GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-linux-amd64 ./cmd/cli

# Rune Server
GOOS=linux GOARCH=arm64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-server-linux-arm64 ./cmd/server

GOOS=linux GOARCH=amd64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-server-linux-amd64 ./cmd/server

GOOS=darwin GOARCH=arm64 \
  go build -ldflags="-s -w -X main.version=${VERSION}" \
  -o dist/rune-server-darwin-arm64 ./cmd/server

chmod +x dist/*

echo ""
echo "Build complete:"
ls -lh dist/