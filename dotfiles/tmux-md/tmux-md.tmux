#!/usr/bin/env bash
# TPM entry point for tmux-md plugin

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FZF_SCRIPT="$CURRENT_DIR/scripts/tmux-md-fzf.sh"

# prefix + m c → search .md files in current pane directory
# prefix + m h → search .md files from home directory
tmux bind-key m switch-client -T tmux-md
tmux bind-key -T tmux-md c run-shell "$FZF_SCRIPT"
tmux bind-key -T tmux-md h run-shell "$FZF_SCRIPT $HOME"
