package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func Install(app string) error {
	switch app {
	case "Homebrew":
		SuccessMessage(`You should install the homebrew manually. Use this command "/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)""`)
	case "Git":
		if err := ExecCommand("brew", "install", "git"); err != nil {
			return err
		}
	case "WezTerm":
		if err := ExecCommand("brew", "install", "--cask", "wezterm"); err != nil {
			return err
		}
	case "Neovim":
		if err := ExecCommand("brew", "install", "neovim"); err != nil {
			return err
		}
	case "tmux":
		if err := ExecCommand("brew", "install", "tmux"); err != nil {
			return err
		}
	case "AeroSpace":
		if err := ExecCommand("brew", "install", "--cask", "nikitabobko/tap/aerospace"); err != nil {
			return err
		}
	case "Homerow":
		if err := ExecCommand("brew", "install", "--cask", "homerow"); err != nil {
			return err
		}
	case "Karabiner Elements":
		if err := ExecCommand("brew", "install", "--cask", "karabiner-elements"); err != nil {
			return err
		}
	case "Snipaste":
		if err := ExecCommand("brew", "install", "--cask", "snipaste"); err != nil {
			return err
		}
	case "ripgrep":
		if err := ExecCommand("brew", "install", "ripgrep"); err != nil {
			return err
		}
	case "fzf":
		if err := ExecCommand("brew", "install", "fzf"); err != nil {
			return err
		}
	case "zsh-vi-mode":
		if err := ExecCommand("brew", "install", "zsh-vi-mode"); err != nil {
			return err
		}
		if err := AddZshSource("source $(brew --prefix)/opt/zsh-vi-mode/share/zsh-vi-mode/zsh-vi-mode.plugin.zsh"); err != nil {
			return err
		}
		WarningMessage("Run source ~/.zshrc to use zsh-vi-mode without reopening the terminal.")
	case "docker":
		if err := ExecCommand("brew", "install", "--cask", "docker"); err != nil {
			return err
		}
	case "go":
		if err := ExecCommand("brew", "install", "go"); err != nil {
			return err
		}
	case "nvm":
		if err := ExecCommand("brew", "install", "nvm"); err != nil {
			return err
		}
		homePath, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		nvmDir := filepath.Join(homePath, ".nvm")
		if err := os.MkdirAll(nvmDir, 0755); err != nil {
			return err
		}
		if err := AddZshSource("export NVM_DIR=\"$HOME/.nvm\"\n [ -s \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\" ] && \\. \"$HOMEBREW_PREFIX/opt/nvm/nvm.sh\"\n [ -s \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\" ] && \\. \"$HOMEBREW_PREFIX/opt/nvm/etc/bash_completion.d/nvm\""); err != nil {
			return err
		}
		WarningMessage("Run source ~/.zshrc to use nvm without reopening the terminal.")
	case "btop":
		if err := ExecCommand("brew", "install", "btop"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Install does not support app:%s", app)
	}

	return nil
}

func IsInstalled(app string) (bool, error) {
	switch app {
	case "Homebrew":
		return ExistCommand("brew"), nil
	case "Git":
		return ExistBrewPackage("git"), nil
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
	case "go":
		return ExistCommand("go"), nil
	case "nvm":
		return ExistCommand("nvm"), nil
	case "btop":
		return ExistCommand("btop"), nil
	default:
		return false, fmt.Errorf("IsInstall does not support app:%s", app)
	}
}
