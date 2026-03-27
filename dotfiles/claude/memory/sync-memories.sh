#!/usr/bin/env bash
# Syncs global Claude memories into all existing project memory directories.
# Run after pulling dotfiles to ensure every project has the shared memories.
#
# How it works:
# 1. Finds all project dirs under ~/.claude/projects/
# 2. For each project that has a memory/ dir, copies _global memories into it
# 3. Rebuilds MEMORY.md index to include the global entries
#
# Usage: bash sync-memories.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
GLOBAL_DIR="$SCRIPT_DIR/_global"
PROJECTS_DIR="$HOME/.claude/projects"

if [ ! -d "$GLOBAL_DIR" ]; then
  echo "No global memories found at $GLOBAL_DIR"
  exit 0
fi

if [ ! -d "$PROJECTS_DIR" ]; then
  echo "No Claude projects directory at $PROJECTS_DIR"
  exit 0
fi

global_files=("$GLOBAL_DIR"/*.md)
if [ ${#global_files[@]} -eq 0 ]; then
  echo "No global memory files found"
  exit 0
fi

synced=0

for project_dir in "$PROJECTS_DIR"/*/; do
  memory_dir="$project_dir/memory"

  # Create memory dir if project exists but has no memories yet
  mkdir -p "$memory_dir"

  # Copy global memories (overwrite to keep them up to date)
  for global_file in "${global_files[@]}"; do
    cp "$global_file" "$memory_dir/"
  done

  # Rebuild MEMORY.md — preserve existing non-global entries, add global ones
  memory_index="$memory_dir/MEMORY.md"
  global_section="## Global (synced from dotfiles)"

  # Collect global entries
  global_entries=""
  for global_file in "${global_files[@]}"; do
    filename="$(basename "$global_file")"
    # Extract description from frontmatter
    desc=$(grep '^description:' "$global_file" | head -1 | sed 's/^description: *//')
    global_entries="$global_entries\n- [$filename]($filename) — $desc"
  done

  if [ -f "$memory_index" ]; then
    # Remove old global section if present, then append new one
    tmp=$(mktemp)
    sed '/^## Global (synced from dotfiles)/,/^## /{ /^## Global/d; /^## /!d; }' "$memory_index" > "$tmp"
    # Remove trailing blank lines
    sed -i '' -e :a -e '/^\n*$/{$d;N;ba' -e '}' "$tmp" 2>/dev/null || true
    mv "$tmp" "$memory_index"
    echo -e "\n$global_section$global_entries" >> "$memory_index"
  else
    echo -e "# Memory Index\n\n$global_section$global_entries" > "$memory_index"
  fi

  synced=$((synced + 1))
done

echo "Synced ${#global_files[@]} global memories to $synced projects"
