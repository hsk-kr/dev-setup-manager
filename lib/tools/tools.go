package tools

import (
	"os/exec"
	"strings"
)

type InstalledSoftwareList struct {
	Homebrew bool
	WezTerm  bool
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
	}

	updateInstalledSoftware()

	return &installedSoftwareList, updateInstalledSoftware
}

func ExistCommand(cmd string) bool {
	output, err := exec.Command(cmd, "-v").Output()

	if err != nil {
		panic(err)
	}

	return !strings.Contains(string(output), "command not found")
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
