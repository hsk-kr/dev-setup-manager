package tools

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
)

type InstalledSoftwareList struct {
	Homebrew bool
	WezTerm  bool
	Neovim   bool
	Tmux     bool
}

func RenderItem(name string, disabled bool) {
	installedPen := color.New(color.FgHiGreen).PrintFunc()
	notInstalledPen := color.New(color.FgRed).PrintFunc()
	softwareNamePen := color.New(color.FgWhite).PrintFunc()

	softwareNamePen(name)
	if disabled {
		installedPen(" - Installed")
	} else {
		notInstalledPen(" - Not Installed")
	}
}

/*
first value contains if software is installed or not,
to update this first value, you should call the function given in the return value
*/
func CreateInstalledSoftwareList() (*InstalledSoftwareList, func()) {
	installedSoftwareList := InstalledSoftwareList{}

	updateInstalledSoftware := func() {
		installedSoftwareList.Homebrew = IsHomebrewInstalled()
		installedSoftwareList.WezTerm = IsWezTermInstalled()
		installedSoftwareList.Neovim = IsNeovimInstalled()
		installedSoftwareList.Tmux = IsTmuxInstalled()
	}

	updateInstalledSoftware()

	return &installedSoftwareList, updateInstalledSoftware
}

func ExistCommand(cmd string) bool {
	_, err := exec.Command("which", cmd).Output()

	if err != nil {
		return false
	}

	return true
}

func ExistApplication(appName string) bool {
	ls := exec.Command("ls", "/Applications")
	grep := exec.Command("grep", appName)

	pipe, err := ls.StdoutPipe()
	defer pipe.Close()

	if err != nil {
		panic(err)
	}

	grep.Stdin = pipe

	err = ls.Start()

	if err != nil {
		panic(err)
	}

	output, err := grep.Output()

	if err != nil {
		return false
	}

	return len(string(output)) >= len(appName)
}

func ExecCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}
