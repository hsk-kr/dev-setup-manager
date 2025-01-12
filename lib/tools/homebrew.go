package tools

import (
	"os"
	"os/exec"
)

func IsHomebrewInstalled() bool {
	return ExistCommand("brew")
}

// TODO: Not tested if it installs homebrew correctly
func InstallHomebrew() {
	cmd := exec.Command("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
