#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SETTINGS_FILE="$SCRIPT_DIR/settings.json"

if ! command -v claude &>/dev/null; then
  echo "claude CLI not found, skipping plugin setup."
  exit 0
fi

if ! command -v jq &>/dev/null; then
  echo "Error: jq not found. Install with: brew install jq"
  exit 1
fi

if [[ ! -f "$SETTINGS_FILE" ]]; then
  echo "settings.json not found at $SETTINGS_FILE, skipping plugin setup."
  exit 0
fi

echo "=== Claude Code Plugin Setup ==="

# Add extra marketplaces
printf "\n--- Adding marketplaces ---\n"
jq -r '.extraKnownMarketplaces // {} | to_entries[] | .value.source.repo' "$SETTINGS_FILE" | while read -r repo; do
  echo "Adding marketplace: $repo"
  claude plugin marketplace add "$repo" 2>/dev/null || echo "  (already added or failed)"
done

# Install enabled plugins
printf "\n--- Installing plugins ---\n"
jq -r '.enabledPlugins // {} | to_entries[] | select(.value == true) | .key' "$SETTINGS_FILE" | while read -r plugin; do
  plugin_name="${plugin%%@*}"
  echo "Installing: $plugin_name"
  claude plugin install "$plugin" 2>/dev/null || echo "  (already installed or failed)"
done

printf "\n=== Done ===\n"
