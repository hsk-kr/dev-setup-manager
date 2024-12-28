package tools

import (
	"os"
	"os/exec"
	"strings"
)

func IsHomebrewInstalled() bool {
	output, err := exec.Command("brew", "-v").Output()

	if err != nil {
		panic(err)
	}

	return strings.HasPrefix(string(output), "Homebrew ")
}

// TODO: Not tested if it installs homebrew correctly
func InstallHomebrew() {
	cmd := exec.Command("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
