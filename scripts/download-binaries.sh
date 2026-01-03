#!/bin/bash

set -e

BINARIES_DIR="internal/bootstrap/binaries"

echo "Downloading yt-dlp and ffmpeg binaries for all platforms..."

download_with_retry() {
    local url=$1
    local dest=$2
    local max_attempts=3
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -L -f -o "$dest" "$url" 2>/dev/null; then
            return 0
        fi
        echo "Attempt $attempt failed, retrying..."
        attempt=$((attempt + 1))
        sleep 2
    done
    return 1
}

echo "Downloading yt-dlp..."
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos" "${BINARIES_DIR}/darwin_amd64/yt-dlp" || echo "Warning: Failed to download yt-dlp for darwin_amd64"
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos_legacy" "${BINARIES_DIR}/darwin_arm64/yt-dlp" || download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos" "${BINARIES_DIR}/darwin_arm64/yt-dlp" || echo "Warning: Failed to download yt-dlp for darwin_arm64"
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux" "${BINARIES_DIR}/linux_amd64/yt-dlp" || echo "Warning: Failed to download yt-dlp for linux_amd64"
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64" "${BINARIES_DIR}/linux_arm64/yt-dlp" || echo "Warning: Failed to download yt-dlp for linux_arm64"
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" "${BINARIES_DIR}/windows_amd64/yt-dlp.exe" || echo "Warning: Failed to download yt-dlp for windows_amd64"
download_with_retry "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" "${BINARIES_DIR}/windows_arm64/yt-dlp.exe" || echo "Warning: Failed to download yt-dlp for windows_arm64"

chmod +x ${BINARIES_DIR}/*/yt-dlp* 2>/dev/null || true

echo "Downloading ffmpeg..."

echo "Downloading ffmpeg for macOS..."
if download_with_retry "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" "${BINARIES_DIR}/darwin_amd64/ffmpeg.zip"; then
    unzip -q -j -o "${BINARIES_DIR}/darwin_amd64/ffmpeg.zip" -d "${BINARIES_DIR}/darwin_amd64/" 2>/dev/null || true
    find "${BINARIES_DIR}/darwin_amd64/" -name "ffmpeg" -type f -exec mv {} "${BINARIES_DIR}/darwin_amd64/ffmpeg" \; 2>/dev/null || true
    rm -f "${BINARIES_DIR}/darwin_amd64/ffmpeg.zip" "${BINARIES_DIR}/darwin_amd64/"*.dylib 2>/dev/null || true
fi

if download_with_retry "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" "${BINARIES_DIR}/darwin_arm64/ffmpeg.zip"; then
    unzip -q -j -o "${BINARIES_DIR}/darwin_arm64/ffmpeg.zip" -d "${BINARIES_DIR}/darwin_arm64/" 2>/dev/null || true
    find "${BINARIES_DIR}/darwin_arm64/" -name "ffmpeg" -type f -exec mv {} "${BINARIES_DIR}/darwin_arm64/ffmpeg" \; 2>/dev/null || true
    rm -f "${BINARIES_DIR}/darwin_arm64/ffmpeg.zip" "${BINARIES_DIR}/darwin_arm64/"*.dylib 2>/dev/null || true
fi

echo "Downloading ffmpeg for Linux..."
if download_with_retry "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz" "${BINARIES_DIR}/linux_amd64/ffmpeg.tar.xz"; then
    tar -xf "${BINARIES_DIR}/linux_amd64/ffmpeg.tar.xz" -C "${BINARIES_DIR}/linux_amd64/" --strip-components=1 --wildcards "*/ffmpeg" 2>/dev/null || true
    rm -f "${BINARIES_DIR}/linux_amd64/ffmpeg.tar.xz" "${BINARIES_DIR}/linux_amd64/"*.md "${BINARIES_DIR}/linux_amd64/"*.txt 2>/dev/null || true
fi

if download_with_retry "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz" "${BINARIES_DIR}/linux_arm64/ffmpeg.tar.xz"; then
    tar -xf "${BINARIES_DIR}/linux_arm64/ffmpeg.tar.xz" -C "${BINARIES_DIR}/linux_arm64/" --strip-components=1 --wildcards "*/ffmpeg" 2>/dev/null || true
    rm -f "${BINARIES_DIR}/linux_arm64/ffmpeg.tar.xz" "${BINARIES_DIR}/linux_arm64/"*.md "${BINARIES_DIR}/linux_arm64/"*.txt 2>/dev/null || true
fi

echo "Downloading ffmpeg for Windows..."
if download_with_retry "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" "${BINARIES_DIR}/windows_amd64/ffmpeg.zip"; then
    unzip -q -j -o "${BINARIES_DIR}/windows_amd64/ffmpeg.zip" "ffmpeg-*-essentials_build/bin/ffmpeg.exe" -d "${BINARIES_DIR}/windows_amd64/" 2>/dev/null || \
    unzip -q -o "${BINARIES_DIR}/windows_amd64/ffmpeg.zip" -d "${BINARIES_DIR}/windows_amd64/" 2>/dev/null && \
    find "${BINARIES_DIR}/windows_amd64/" -name "ffmpeg.exe" -type f -exec mv {} "${BINARIES_DIR}/windows_amd64/ffmpeg.exe" \; 2>/dev/null || true
    rm -f "${BINARIES_DIR}/windows_amd64/ffmpeg.zip" 2>/dev/null || true
    rm -rf "${BINARIES_DIR}/windows_amd64/ffmpeg-"* 2>/dev/null || true
fi

if download_with_retry "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" "${BINARIES_DIR}/windows_arm64/ffmpeg.zip"; then
    unzip -q -j -o "${BINARIES_DIR}/windows_arm64/ffmpeg.zip" "ffmpeg-*-essentials_build/bin/ffmpeg.exe" -d "${BINARIES_DIR}/windows_arm64/" 2>/dev/null || \
    unzip -q -o "${BINARIES_DIR}/windows_arm64/ffmpeg.zip" -d "${BINARIES_DIR}/windows_arm64/" 2>/dev/null && \
    find "${BINARIES_DIR}/windows_arm64/" -name "ffmpeg.exe" -type f -exec mv {} "${BINARIES_DIR}/windows_arm64/ffmpeg.exe" \; 2>/dev/null || true
    rm -f "${BINARIES_DIR}/windows_arm64/ffmpeg.zip" 2>/dev/null || true
    rm -rf "${BINARIES_DIR}/windows_arm64/ffmpeg-"* 2>/dev/null || true
fi

chmod +x ${BINARIES_DIR}/*/ffmpeg* 2>/dev/null || true

echo ""
echo "Download complete!"
echo "Verifying binaries..."
for platform in darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 windows_amd64 windows_arm64; do
    ytdlp_ok=false
    ffmpeg_ok=false
    
    if [ -f "${BINARIES_DIR}/${platform}/yt-dlp" ] || [ -f "${BINARIES_DIR}/${platform}/yt-dlp.exe" ]; then
        ytdlp_ok=true
        echo "✓ yt-dlp found for ${platform}"
    else
        echo "✗ yt-dlp missing for ${platform}"
    fi
    
    if [ -f "${BINARIES_DIR}/${platform}/ffmpeg" ] || [ -f "${BINARIES_DIR}/${platform}/ffmpeg.exe" ]; then
        ffmpeg_ok=true
        echo "✓ ffmpeg found for ${platform}"
    else
        echo "✗ ffmpeg missing for ${platform}"
    fi
done
