#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INIT="$SCRIPT_DIR/init.lua"

# Capture the current window so we can target its panes from inside the popup
TARGET_WINDOW=$(tmux display-message -p '#{session_name}:#{window_index}')

CS_TARGET_WINDOW="$TARGET_WINDOW" tmux display-popup -w 80% -h 40% -E \
  "nvim --clean -u '$INIT'"
