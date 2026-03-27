# tmux-md

Markdown viewer for tmux — renders `.md` files in a popup using neovim with treesitter syntax highlighting.

## Usage

### 4 ways to open

| Method | How |
|--------|-----|
| **tmux keybind** | `prefix + m` — fzf picks `.md` files in cwd, opens in popup |
| **neovim command** | `:TmuxMd` — opens current buffer in popup |
| **neovim visual** | Select a file path, `<leader>md` — opens that file in popup |
| **CLI** | `tmux-md file.md` — from shell, scripts, or Claude Code |

Press `q`, `<Esc>`, or `:q` to close the viewer.

## Setup

### Prerequisites

- tmux
- neovim with treesitter + tokyonight installed via lazy.nvim
- fzf + fzf-tmux (for the picker)

### Install

Already integrated into the dotfiles setup:

- `tmux.conf` loads the TPM entry point
- `nvim/init.lua` loads the neovim commands
- CLI symlink: `ln -sf ~/dev-setup-manager/dotfiles/tmux-md/scripts/tmux-md.sh /usr/local/bin/tmux-md`

### Manual reload

```bash
# Reload tmux config
tmux source-file ~/.config/tmux/tmux.conf

# In neovim, the module loads on startup
```

## How it works

The viewer uses `nvim --clean -u viewer-init.lua -R` for fast startup (~50ms). It reuses your existing treesitter parsers and tokyonight theme from `~/.local/share/nvim/lazy/` without loading your full config.
