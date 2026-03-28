#!/bin/bash
# Sets the tmux pane border status for Claude Code agent panels.
# Called from Claude Code hooks with: bash agent-status.sh <status>
# Status: "working" (green) or "attention" (yellow)
#
# NOT WORKING YET: The script works when called manually, but Claude Code
# hooks don't seem to trigger tmux pane option updates reliably.
# The pane-border-format conditional coloring is set up in lib/tools/agents.go
# and works when @agent-status is set manually. The issue is getting the
# hooks to fire and update the status in practice.
#
# To enable, add these to settings.json hooks:
#   Notification idle_prompt:    bash ~/.claude/hooks/agent-status.sh attention
#   Notification permission_prompt: bash ~/.claude/hooks/agent-status.sh attention
#   PreToolUse (empty matcher):  bash ~/.claude/hooks/agent-status.sh working

STATUS="${1:-working}"

# Only run inside tmux
[ -z "$TMUX_PANE" ] && exit 0

tmux set-option -t "$TMUX_PANE" -p @agent-status "$STATUS" 2>/dev/null
tmux refresh-client 2>/dev/null
exit 0
