#!/usr/bin/env bash
# tmux-md: Open markdown files in a tmux popup with neovim readonly viewer
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VIEWER_INIT="$SCRIPT_DIR/viewer-init.lua"

usage() {
  echo "Usage: tmux-md <file.md>"
  echo "  Opens a markdown file in a tmux popup with syntax highlighting."
  echo "  Falls back to plain neovim if not inside tmux."
  exit 1
}

[[ $# -lt 1 ]] && usage

FILE="$1"

# Resolve to absolute path
if [[ ! "$FILE" = /* ]]; then
  FILE="$(pwd)/$FILE"
fi

if [[ ! -f "$FILE" ]]; then
  echo "Error: File not found: $FILE" >&2
  exit 1
fi

NVIM_CMD="nvim --clean -u '$VIEWER_INIT' -R '$FILE'"

if [[ -n "${TMUX:-}" ]]; then
  tmux display-popup -w 90% -h 90% -E "$NVIM_CMD"
else
  eval "$NVIM_CMD"
fi
