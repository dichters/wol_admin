#!/bin/bash
# Build & package wol-panel for Linux x86-64
# 用法：./bin/sh/build-linux-x86-64.sh [version]
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')}"
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')

# 前端：dist/ 不存在时才构建
if [ ! -d "dist" ]; then
    echo "Building frontend..."
    cd frontend && npm run build && cd ..
fi

echo "Building wol-panel v${VERSION} linux/amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w -X wol-panel/version.Version=${VERSION} -X wol-panel/version.Arch=amd64 -X 'wol-panel/version.BuildTime=${BUILD_TIME}'" \
    -o build/wol-panel .

RELEASE_DIR="release"
mkdir -p "$RELEASE_DIR"

PKG_NAME="wol-panel-${VERSION}-linux-x86-64"
PKG_DIR="${RELEASE_DIR}/${PKG_NAME}"
rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR"

cp build/wol-panel "${PKG_DIR}/"
cp config.template.json "${PKG_DIR}/"
cp wol-panel.service "${PKG_DIR}/"

rm -f "${RELEASE_DIR}/${PKG_NAME}.zip"
pushd "$RELEASE_DIR" > /dev/null
zip -r "${PKG_NAME}.zip" "$PKG_NAME"
popd > /dev/null
rm -rf "$PKG_DIR"

echo "Done: ${RELEASE_DIR}/${PKG_NAME}.zip"
