package tools

import (
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// TODO: Not tested if it installs homebrew correctly
func Homebrew() {
	print := color.New(color.FgHiRed).PrintlnFunc()
	output, err := exec.Command("brew", "-v").Output()

	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(string(output), "Homebrew ") {
		print("Homebrew is already installed.")
		return
	}

	cmd := exec.Command("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
