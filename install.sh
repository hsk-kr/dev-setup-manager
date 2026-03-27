#!/bin/bash
set -e

REPO="hsk-kr/licokit"
INSTALL_DIR="$HOME/licokit"
BINARY="licokit"

echo "Downloading latest release..."
DOWNLOAD_URL=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep "browser_download_url.*licokit" | head -1 | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Error: Could not find latest release"
  exit 1
fi

mkdir -p "$INSTALL_DIR"
curl -sL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"
xattr -d com.apple.quarantine "$INSTALL_DIR/$BINARY" 2>/dev/null || true

echo "Installed to $INSTALL_DIR/$BINARY"
echo "Running..."
"$INSTALL_DIR/$BINARY"
