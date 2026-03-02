package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupDotfiles() error {
	homePath, err := os.UserHomeDir()

	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	devSetupManagerHomePath := filepath.Join(homePath, "dev-setup-manager")
	configDirPath := filepath.Join(homePath, ".config")
	devSetupManagerDotfilesPath := filepath.Join(devSetupManagerHomePath, "dotfiles")

	if err := ExecCommand("mkdir", "-p", devSetupManagerHomePath); err != nil {
		return err
	}

	if exist, _ := existFile(devSetupManagerDotfilesPath); exist {
		if err := ExecCommand("rm", "-rf", devSetupManagerDotfilesPath); err != nil {
			return err
		}
	}

	if err := ExecCommand("git", "clone", "git@github.com:hsk-kr/dotfiles.git", devSetupManagerDotfilesPath); err != nil {
		return err
	}

	if err := ExecCommand("mkdir", "-p", configDirPath); err != nil {
		return err
	}

	copyItems := []string{
		"aerospace",
		"devdeck",
		"karabiner",
		"nvim",
		"tmux",
		"zsh",
		"alacritty",
	}

	for _, copyItem := range copyItems {
		if err := ExecCommand("ln", "-sfn", filepath.Join(devSetupManagerDotfilesPath, copyItem), filepath.Join(configDirPath, copyItem)); err != nil {
			return err
		}
	}

	if err := ExecCommand("ln", "-sfn", filepath.Join(devSetupManagerDotfilesPath, "scripts"), filepath.Join(homePath, "scripts")); err != nil {
		return err
	}

	return AddZshSource(fmt.Sprintf("source %s", filepath.Join(configDirPath, "zsh", "zshrc")))
}
