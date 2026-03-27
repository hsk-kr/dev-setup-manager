<p align="center">
  <h1 align="center">licokit</h1>
  <p align="center">
    A TUI tool that bootstraps your entire macOS dev environment after a fresh reset.<br/>
    Select tools, sync dotfiles, get coding — all from one command.
  </p>
</p>

<p align="center">
  <a href="#quick-install">Quick Install</a> ·
  <a href="#features">Features</a> ·
  <a href="#dotfiles">Dotfiles</a> ·
  <a href="#configuration">Configuration</a> ·
  <a href="#build">Build</a>
</p>

---

## Target Environment

| Category | Value |
|----------|-------|
| Device   | MacBook Air M1 |
| OS       | macOS (darwin/arm64) |
| Shell    | zsh |

## Quick Install

```bash
curl -sL https://raw.githubusercontent.com/hsk-kr/licokit/main/install.sh | bash
```

> Downloads the latest release binary to `~/licokit`, removes the macOS quarantine flag, and runs it.

## Features

### 1. Tools Installation

Browse and install dev tools through an interactive TUI with real-time status detection and a progress spinner. Tools are defined in a YAML config with version pinning support.

<details>
<summary><strong>Default tools (18)</strong></summary>

| Tool | Install Method | Category |
|------|---------------|----------|
| Homebrew | Manual | Package Manager |
| Git | `brew install` | Development |
| WezTerm | `brew install --cask` | Terminal |
| Ghostty | `brew install --cask` | Terminal |
| Neovim | `brew install` | Editor |
| tmux | `brew install` | Terminal Multiplexer |
| AeroSpace | `brew install --cask` | Window Manager |
| Neru | `brew install --cask` | Keyboard Navigation |
| Karabiner Elements | `brew install --cask` | Key Remapper |
| Snipaste | `brew install --cask` | Screenshot |
| ripgrep | `brew install` | Search |
| fzf | `brew install` | Fuzzy Finder |
| zsh-vi-mode | `brew install` | Shell Plugin |
| Docker | `brew install --cask` | Containers |
| Go | `brew install` | Language |
| nvm | `brew install` | Version Manager |
| btop | `brew install` | System Monitor |
| Claude Code | Script (`curl`) | AI Assistant |

</details>

> **Note:** Homebrew must be installed first. The app will show the install command.

### 2. Dotfiles

Sets up dotfiles that live **inside this repository** (`dotfiles/` directory) — no separate repo needed. Runs `git pull` to sync, then creates symlinks under `~/.config` and `~`.

<details>
<summary><strong>What gets symlinked</strong></summary>

**Config directories** → `~/.config/<name>`

| Directory | Purpose |
|-----------|---------|
| `aerospace` | Tiling window manager |
| `alacritty` | Terminal emulator |
| `devdeck` | DevDeck dashboard |
| `ghostty` | Terminal emulator |
| `karabiner` | Key remapping |
| `neru` | Keyboard-driven navigation |
| `nvim` | Neovim (35+ plugins, LSP) |
| `tmux` | Terminal multiplexer |
| `zsh` | Shell configuration |

**Home directory links** → `~/<name>`

| Source | Target |
|--------|--------|
| `scripts` | `~/scripts` |

**Extra links**

| Source | Target |
|--------|--------|
| `claude/skills` | `~/.claude/skills` |
| `claude/agents` | `~/.claude/agents` |
| `claude/commands` | `~/.claude/commands` |
| `claude/docs` | `~/.claude/docs` |
| `claude/hooks` | `~/.claude/hooks` |
| `claude/settings.json` | `~/.claude/settings.json` |
| `claude/CLAUDE.md` | `~/.claude/CLAUDE.md` |
| `claude/statusline-command.sh` | `~/.claude/statusline-command.sh` |

**Post-setup scripts**

- `claude/memory/sync-memories.sh` — syncs cross-session memory
- `claude/setup-plugins.sh` — installs Claude Code plugins

</details>

### 3. Guide

Built-in setup notes for configurations that need manual attention.

## Dotfiles

The `dotfiles/` directory is a self-contained collection of configurations for the full dev environment:

```
dotfiles/
├── aerospace/        # Window tiling
├── alacritty/        # Terminal
├── claude/           # Claude Code (agents, skills, commands, hooks, memory)
├── devdeck/          # Dashboard
├── ghostty/          # Terminal
├── karabiner/        # Key remapping
├── neru/             # Keyboard navigation
├── nvim/             # Neovim (lazy.nvim, LSP, custom scripts)
├── opencode/         # OpenCode agents
├── scripts/          # Utility scripts
├── tmux/             # Terminal multiplexer
├── tmux-md/          # Markdown-based tmux manager
├── vscode/           # VS Code keybindings
├── wezterm/          # Terminal (with backgrounds)
└── zsh/              # Shell config
```

## Configuration

The app ships with a sensible default config. Override it by creating:

```
~/.config/licokit/config.yaml
```

### Config Format

```yaml
dotfiles:
  config_links:
    - nvim
    - tmux
  home_links:
    scripts: scripts
  extra_links:
    - source: claude/skills
      target: ~/.claude/skills
  post_scripts:
    - claude/setup-plugins.sh
  zsh_source: "~/.config/zsh/zshrc"

tools:
  - name: Go
    install_type: brew          # brew | cask | manual | script
    package: go
    version: "1.23"             # optional — pins to go@1.23
    detect_type: command        # command | application | brew_package
    detect_value: go

  - name: Docker
    install_type: cask
    package: docker
    detect_type: command
    detect_value: docker

  - name: Claude Code
    install_type: script
    install_command: "curl -fsSL https://claude.ai/install.sh | bash"
    detect_type: command
    detect_value: claude

  - name: nvm
    install_type: brew
    package: nvm
    detect_type: brew_package
    detect_value: nvm
    zsh_source: |               # optional — added to dev.zsh after install
      export NVM_DIR="$HOME/.nvm"
    post_install_dirs:          # optional — directories to create
      - ~/.nvm
    post_install_warning: "Run source ~/.zshrc"  # optional
```

## Usage

### From Source

```bash
# Prerequisites: Go 1.23+
git clone https://github.com/hsk-kr/licokit.git
cd licokit
go run .
```

### Navigation

| Key | Action |
|-----|--------|
| `j` / `J` / `h` / `H` | Move down |
| `k` / `K` / `l` / `L` | Move up |
| `Enter` | Select |
| `ESC` | Back / Exit |

## Build

```bash
GOOS=darwin GOARCH=arm64 go build -o licokit
```

## License

[MIT](LICENSE)
