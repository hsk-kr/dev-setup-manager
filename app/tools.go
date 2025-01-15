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

	items := []terminal.SelectItem{
		{
			Name:     "Homebrew",
			Render:   RenderItem(installedSoftwareList.Homebrew),
			Disabled: installedSoftwareList.Homebrew,
			Install:  tools.InstallHomebrew,
		},
		{
			Name:     "WezTerm",
			Render:   RenderItem(installedSoftwareList.WezTerm),
			Disabled: installedSoftwareList.WezTerm,
			Install:  tools.InstallWezTerm,
		},
		{
			Name:     "Neovim",
			Render:   RenderItem(installedSoftwareList.Neovim),
			Disabled: installedSoftwareList.Neovim,
			Install:  tools.InstallNeovim,
		},
		{
			Name: "tmux",
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
		{
			Name: "docker",
		},
		{
			Name: "nvm",
		},
		{
			Name: "gvm",
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

		for i, software := range items {
			if software.Name == choice {
				if software.Install == nil {
					NotSupported(choice)
				} else if !software.Disabled {
					software.Install()
					items[i].Render = RenderItem(true)
					items[i].Disabled = true
				}

				break
			}
		}
	}
}

func RenderItem(isInstalled bool) func(string) {
	return func(name string) {
		installedPen := color.New(color.FgHiGreen).PrintFunc()
		notInstalledPen := color.New(color.FgRed).PrintFunc()
		softwareNamePen := color.New(color.FgWhite).PrintFunc()

		softwareNamePen(name)
		if isInstalled {
			installedPen(" - Installed")
		} else {
			notInstalledPen(" - Not Installed")
		}
	}
}
