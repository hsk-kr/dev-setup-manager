package tools

import (
	"os"
	"os/exec"
)

func IsWezTermInstalled() bool {
	return ExistApplication("WezTerm.app")
}

func InstallWezTerm() {
	cmd := exec.Command("brew install --cask wezterm")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
