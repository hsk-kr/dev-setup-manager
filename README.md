# Dev Setup Manager

A TUI (Terminal User Interface) tool that automates development environment setup after a fresh macOS reset. Select and install your tools, sync your dotfiles, and get coding — all from a single command.

## Target Environment

| Category | Value |
| --- | --- |
| Device | MacBook Air M1 |
| OS | macOS (darwin/arm64) |
| Shell | zsh |

## Features

### 1. Tools Installation

Browse and install development tools through an interactive menu with real-time installation status and a progress spinner. Tools are defined in a YAML config file with version pinning support.

Default tools included:

| Tool | Install Method | Category |
| --- | --- | --- |
| Homebrew | Manual | Package Manager |
| Git | `brew install` | Development |
| Neovim | `brew install` | Editor |
| Go | `brew install` | Language |
| nvm | `brew install` | Version Manager |
| Docker | `brew install --cask` | Development |
| WezTerm | `brew install --cask` | Terminal |
| tmux | `brew install` | Terminal |
| zsh-vi-mode | `brew install` | Shell Plugin |
| ripgrep | `brew install` | Search |
| fzf | `brew install` | Search |
| btop | `brew install` | System Monitor |
| AeroSpace | `brew install --cask` | Window Manager |
| Homerow | `brew install --cask` | Keyboard Navigation |
| Karabiner Elements | `brew install --cask` | Key Remapper |
| Snipaste | `brew install --cask` | Screenshot |

> **Note:** Homebrew must be installed manually first. The app will show the installation command.

### 2. Dotfiles

Clones a dotfiles repository and creates symlinks under `~/.config`. The repo URL is configurable via the YAML config file (defaults to [hsk-kr/dotfiles](https://github.com/hsk-kr/dotfiles)).

Default symlinks:
- Config directories: aerospace, devdeck, karabiner, nvim, tmux, zsh, alacritty
- Scripts directory (`~/scripts`)
- Claude Code skills (`~/.claude/skills`)
- Zsh configuration sourced from `~/.zshrc`

### 3. Guide

Built-in setup notes for configurations and installations that need manual attention.

## Configuration

Tools and dotfiles are defined in a YAML config file. The app uses a built-in default config, but you can override it by creating:

```
~/.config/dev-setup-manager/config.yaml
```

### Config Format

```yaml
dotfiles:
  repo: "git@github.com:your-user/dotfiles.git"
  config_links:
    - nvim
    - tmux
  home_links:
    scripts: scripts
  extra_links:
    - source: claude/skills
      target: ~/.claude/skills
  zsh_source: "source ~/.config/zsh/zshrc"

tools:
  - name: Go
    install_type: brew        # brew, cask, or manual
    package: go
    version: "1.23"           # optional, pins to go@1.23
    detect_type: command      # command, application, or brew_package
    detect_value: go

  - name: Docker
    install_type: cask
    package: docker
    detect_type: command
    detect_value: docker

  - name: nvm
    install_type: brew
    package: nvm
    detect_type: brew_package
    detect_value: nvm
    zsh_source: |             # optional, added to ~/.zshrc after install
      export NVM_DIR="$HOME/.nvm"
    post_install_dirs:        # optional, directories to create after install
      - ~/.nvm
    post_install_warning: "Run source ~/.zshrc"  # optional
```

## Usage

### Quick Install

Copy and paste this into your terminal after a fresh reset:

```bash
curl -sL https://raw.githubusercontent.com/hsk-kr/dev-setup-manager/main/install.sh | bash
```

This downloads the latest release, removes the macOS quarantine flag, and runs it.

### From Source

```bash
# Prerequisites: Go 1.23+
git clone https://github.com/hsk-kr/dev-setup-manager.git
cd dev-setup-manager
go run .
```

### Navigation

| Key | Action |
| --- | --- |
| `j` / `J` / `h` / `H` | Move down |
| `k` / `K` / `l` / `L` | Move up |
| `Enter` | Select |
| `ESC` | Back / Exit |

## Build

```bash
GOOS=darwin GOARCH=arm64 go build -o dev-setup-manager
```
