#!/bin/bash
# Build wol_admin for Linux x64
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')}"
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')

echo "Building frontend..."
cd frontend && npm run build && cd ..

echo "Building wol_admin v${VERSION} linux/amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "-s -w -X wol_admin/version.Version=${VERSION} -X wol_admin/version.Arch=amd64 -X 'wol_admin/version.BuildTime=${BUILD_TIME}'" \
  -o build/wol_admin .
