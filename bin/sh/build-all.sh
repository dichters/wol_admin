#!/bin/bash
# CI 专用批量构建脚本：交叉编译所有平台并打包 release 压缩包
# 前提：前端已构建完成（dist/ 目录已存在）
# 用法：./bin/sh/build-all.sh <version>
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    VERSION="$(grep 'Version\s*=' version/version.go | head -1 | grep -oP '"\K[^"]+')"
fi
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')

# 检查前端是否已构建
if [ ! -d "dist" ]; then
    echo "ERROR: dist/ not found. Please build frontend first: cd frontend && npm run build"
    exit 1
fi

RELEASE_DIR="release"
BUILD_DIR="build"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR" "$BUILD_DIR"

# 平台定义：名称|GOOS|GOARCH|输出文件名|压缩格式
PLATFORMS=(
    "linux-x64|linux|amd64|wol_admin|tar.gz"
    "linux-arm64|linux|arm64|wol_admin|tar.gz"
    "windows-x64|windows|amd64|wol_admin.exe|zip"
    "windows-arm|windows|arm|wol_admin.exe|zip"
    "macos-apple-silicon|darwin|arm64|wol_admin|tar.gz"
    "macos-intel|darwin|amd64|wol_admin|tar.gz"
)

echo "========================================"
echo "Building wol_admin v${VERSION}"
echo "Build time: ${BUILD_TIME}"
echo "========================================"

for platform_def in "${PLATFORMS[@]}"; do
    IFS='|' read -r PLATFORM_NAME GOOS GOARCH OUTPUT_NAME ARCHIVE_FMT <<< "$platform_def"

    echo ""
    echo "--- Building ${PLATFORM_NAME} (${GOOS}/${GOARCH}) ---"

    # 交叉编译
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags "-s -w -X wol_admin/version.Version=${VERSION} -X wol_admin/version.Arch=${GOARCH} -X 'wol_admin/version.BuildTime=${BUILD_TIME}'" \
        -o "build/${OUTPUT_NAME}" .

    # 创建临时打包目录
    PKG_NAME="wol_admin-${VERSION}-${PLATFORM_NAME}"
    PKG_DIR="${RELEASE_DIR}/${PKG_NAME}"
    mkdir -p "$PKG_DIR"

    # 复制文件
    cp "build/${OUTPUT_NAME}" "${PKG_DIR}/"
    cp config.template.json "${PKG_DIR}/"
    cp wol_admin.service "${PKG_DIR}/"
    cp README.md "${PKG_DIR}/"

    # 打包
    pushd "$RELEASE_DIR" > /dev/null
    if [ "$ARCHIVE_FMT" = "zip" ]; then
        zip -r "${PKG_NAME}.zip" "$PKG_NAME"
    else
        tar -czf "${PKG_NAME}.tar.gz" "$PKG_NAME"
    fi
    popd > /dev/null

    # 清理临时目录
    rm -rf "$PKG_DIR"

    echo "Done: ${RELEASE_DIR}/${PKG_NAME}.${ARCHIVE_FMT}"
done

# 清理 build 目录中的二进制
rm -f build/wol_admin build/wol_admin.exe

echo ""
echo "========================================"
echo "All builds completed!"
echo "Release files in ${RELEASE_DIR}/:"
ls -la "$RELEASE_DIR"
echo "========================================"