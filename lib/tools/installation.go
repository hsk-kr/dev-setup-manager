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
		AddZshSource("source $(brew --prefix)/opt/zsh-vi-mode/share/zsh-vi-mode/zsh-vi-mode.plugin.zsh")
		WarningMessage("Run source ~/.zshrc to use zsh-vi-mode without reopening the terminal.")
	case "docker":
		ExecCommand("brew", "install", "docker")
	case "nvm":
		ExecCommand("brew", "install", "nvm")
		ExecCommand("mkdir", "~/.nvm")
		AddZshSource("export NVM_DIR=\"$HOME/.nvm\"\n [ -s \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\" ] && \\. \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\"\n [ -s \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\" ] && \\. \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\"")
		WarningMessage("Run source ~/.zshrc to use nvm without reopening the terminal.")
	case "gvm":
		ExecCommandWithIgnoreError("xcode-select", "--install")
		ExecCommand("brew", "update")
		ExecCommand("brew", "install", "mercurial")
		ExecCommand("brew", "install", "bison")
		ExecCommand("bash", "-c", "bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)")
		WarningMessage("Follow the last instruction to complete gvm installation to use on this terminal session")
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
	case "docker":
		return ExistCommand("docker"), nil
	case "nvm":
		return ExistCommand("nvm"), nil
	case "gvm":
		return ExistCommand("gvm"), nil
	default:
		return false, errors.New(fmt.Sprintf("IsInstall does not support app:%s\n", app))
	}
}
