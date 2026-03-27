#!/usr/bin/env bash
# tmux-md-fzf: Find and pick a markdown file with fzf, then open in tmux-md viewer
# Usage: tmux-md-fzf.sh [search_path]
#   search_path: directory to search in (default: current pane's working directory)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TMUX_MD="$SCRIPT_DIR/tmux-md.sh"

if [[ -n "${1:-}" ]]; then
  BASE_PATH="$1"
else
  BASE_PATH="$(tmux display-message -p '#{pane_current_path}')"
fi

SELECTED=$(fd --type f --extension md --max-depth 5 . "$BASE_PATH" 2>/dev/null | \
  sed "s|^$BASE_PATH/||" | \
  sort | \
  fzf-tmux -p 60%,50% --preview "head -40 $BASE_PATH/{}" --preview-window=right:50%:wrap) || true

if [[ -n "$SELECTED" ]]; then
  tmux run-shell -b "sleep 0.1 && '$TMUX_MD' '$BASE_PATH/$SELECTED'"
fi
