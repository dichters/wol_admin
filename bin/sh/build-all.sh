#!/bin/bash
# CI 专用批量构建脚本：构建前端后依次调用各平台构建脚本
# 用法：./bin/sh/build-all.sh [version]
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')}"
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')

RELEASE_DIR="release"
mkdir -p "$RELEASE_DIR"

# 构建前端（仅构建一次，后续各脚本检测到 dist/ 存在会跳过）
if [ ! -d "frontend/dist" ]; then
    echo "========================================"
    echo "Building frontend..."
    echo "========================================"
    cd frontend && npm run build && cd ..
fi

echo ""
echo "========================================"
echo "Building wol-panel v${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "========================================"

# 依次调用各平台构建脚本（编译 + 打包）
bash bin/sh/build-linux-aarch64.sh   "$VERSION"
bash bin/sh/build-linux-x86-64.sh    "$VERSION"
bash bin/sh/build-macos-darwin.sh    "$VERSION"
bash bin/sh/build-macos-intel.sh     "$VERSION"
bash bin/sh/build-windows-x86-64.sh  "$VERSION"

echo ""
echo "========================================"
echo "All builds completed!"
echo "Release files in ${RELEASE_DIR}/:"
ls -la "$RELEASE_DIR"
echo "========================================"
