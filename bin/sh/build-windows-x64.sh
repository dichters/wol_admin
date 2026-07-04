#!/bin/bash
# Build wol_admin for Windows x64
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')}"
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

echo "Building frontend..."
cd frontend && npm run build && cd ..

echo "Building wol_admin v${VERSION} windows/amd64"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
  -ldflags "-s -w -X wol_admin/version.Version=${VERSION} -X wol_admin/version.Arch=amd64 -X wol_admin/version.BuildTime=${BUILD_TIME}" \
  -o build/wol_admin.exe .
