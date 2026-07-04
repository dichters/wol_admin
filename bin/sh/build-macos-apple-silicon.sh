#!/bin/bash
# Build wol_admin for macOS Apple Silicon
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')}"
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

echo "Building frontend..."
cd frontend && npm run build && cd ..

echo "Building wol_admin v${VERSION} darwin/arm64"
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
  -ldflags "-s -w -X wol_admin/version.Version=${VERSION} -X wol_admin/version.Arch=arm64 -X wol_admin/version.BuildTime=${BUILD_TIME}" \
  -o build/wol_admin .
