package tools

import (
	"errors"
	"fmt"
)

func Install(app string) error {
	switch app {
	case "Homebrew":
		ExecCommand("/bin/bash", "-c", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	case "WezTerm":
		ExecCommand("brew", "install", "--cask", "wezterm")
	case "Neovim":
		ExecCommand("brew", "install", "neovim")
	case "tmux":
		ExecCommand("brew", "install", "tmux")
	case "AeroSpace":
		ExecCommand("brew", "install", "--cask", "nikitabobko/tap/aerospace")
	case "Homerow":
		ExecCommand("brew", "install", "--cask", "homerow")
	case "Karabiner Elements":
		ExecCommand("brew", "install", "--cask", "karabiner-elements")
	case "Snipaste":
		ExecCommand("brew", "install", "--cask", "snipaste")
	case "ripgrep":
		ExecCommand("brew", "install", "ripgrep")
	case "fzf":
		ExecCommand("brew", "install", "fzf")
	case "zsh-vi-mode":
		ExecCommand("brew", "install", "zsh-vi-mode")
	default:
		return errors.New(fmt.Sprintf("Install does not support app:%s\n", app))
	}

	return nil
}

func IsInstalled(app string) (bool, error) {
	switch app {
	case "Homebrew":
		return ExistCommand("brew"), nil
	case "WezTerm":
		return ExistApplication("WezTerm.app"), nil
	case "Neovim":
		return ExistCommand("nvim"), nil
	case "tmux":
		return ExistCommand("tmux"), nil
	case "AeroSpace":
		return ExistApplication("AeroSpace.app"), nil
	case "Homerow":
		return ExistApplication("Homerow.app"), nil
	case "Karabiner Elements":
		return ExistApplication("Karabiner-Elements.app"), nil
	case "Snipaste":
		return ExistApplication("Snipaste.app"), nil
	case "ripgrep":
		return ExistCommand("rg"), nil
	case "fzf":
		return ExistCommand("fzf"), nil
	case "zsh-vi-mode":
		return ExistBrewPackage("zsh-vi-mode"), nil
	default:
		return false, errors.New(fmt.Sprintf("IsInstall does not support app:%s\n", app))
	}
}