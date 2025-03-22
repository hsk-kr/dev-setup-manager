package tools

import (
	"fmt"
	"os"
)

func SetupDotfiles() error {
	homePath, err := os.UserHomeDir()

	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	devSetupManagerHomePath := homePath + "/dev-setup-manager"
	configDirPath := homePath + "/.config"
	devSetupManagerDotfilesPath := devSetupManagerHomePath + "/dotfiles"
	ExecCommand("mkdir", "-p", devSetupManagerHomePath)

	if exist, _ := existFile(devSetupManagerDotfilesPath); exist {
		ExecCommand("rm", "-rf", devSetupManagerDotfilesPath)
	}

	ExecCommand("git", "clone", "git@github.com:hsk-kr/dotfiles.git", devSetupManagerDotfilesPath)

	ExecCommand("mkdir", "-p", configDirPath)

	copyItems := []string{
		"aerospace",
		"devdeck",
		"karabiner",
		"nvim",
		"tmux",
		"zsh",
	}

	for _, copyItem := range copyItems {
		ExecCommand("ln", "-sfn", fmt.Sprintf("%s/%s", devSetupManagerDotfilesPath, copyItem), fmt.Sprintf("%s/%s", configDirPath, copyItem))
	}

	ExecCommand("ln", "-sfn", fmt.Sprintf("%s/%s", devSetupManagerDotfilesPath, "scripts"), fmt.Sprintf("%s/%s", homePath, "scripts"))

	AddZshSource(fmt.Sprintf("source %s/%s", configDirPath, "zsh/zshrc"))

	return nil
}
