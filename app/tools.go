package app

import (
	"github.com/fatih/color"
	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
	"github.com/hsk-kr/dev-setup-manager/lib/tools"
)

func Tools() {
	loadingPen := color.New(color.FgHiGreen).PrintfFunc()

	installedSoftwareListChan := make(chan *tools.InstalledSoftwareList)

	loadingPen("Reading installed software...\n")

	go func() {
		installedSoftwareList, _ := tools.CreateInstalledSoftwareList()
		installedSoftwareListChan <- installedSoftwareList
	}()

	installedSoftwareList := <-installedSoftwareListChan

	printSoftwareAvailability := func(name string, isInstalled bool) {
		installedPen := color.New(color.FgHiGreen).PrintFunc()
		notInstalledPen := color.New(color.FgRed).PrintFunc()
		softwareNamePen := color.New(color.FgWhite).PrintFunc()

		softwareNamePen(name)
		if isInstalled {
			installedPen("- Installed")
		} else {
			notInstalledPen("- Not Installed")
		}
	}

	items := []terminal.SelectItem{
		{
			Name: "Homebrew",
			Render: func(name string) {
				printSoftwareAvailability(name, installedSoftwareList.Homebrew)
			},
			Disabled: installedSoftwareList.Homebrew,
		},
		{
			Name: "WezTerm",
		},

		{
			Name: "Neovim",
		},

		{
			Name: "Aerospace",
		},

		{
			Name: "Homerow",
		},

		{
			Name: "Karabiner Elements",
		},

		{
			Name: "Snipaste",
		},

		{
			Name: "fzf",
		},
		{
			Name: "zsh-vi-mode",
		},
	}

	display.DisplayHeader()
	print := color.New(color.FgGreen).PrintlnFunc()
	print("Select item you install")

	for {
		choice, err := terminal.Select(items)

		if err != nil {
			return
		}

		switch choice {
		case "Homebrew":
			tools.InstallHomebrew()
		default:
			NotSupported(choice)
			continue
		}
	}
}
