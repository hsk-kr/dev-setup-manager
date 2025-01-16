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
			Name:        "Homebrew",
			Render:      tools.RenderItem,
			Disabled:    installedSoftwareList.Homebrew,
			Run:         tools.InstallHomebrew,
			GetDisabled: tools.IsHomebrewInstalled,
		},
		{
			Name:        "WezTerm",
			Render:      tools.RenderItem,
			Disabled:    installedSoftwareList.WezTerm,
			Run:         tools.InstallWezTerm,
			GetDisabled: tools.IsWezTermInstalled,
		},
		{
			Name:        "Neovim",
			Render:      tools.RenderItem,
			Disabled:    installedSoftwareList.Neovim,
			Run:         tools.InstallNeovim,
			GetDisabled: tools.IsNeovimInstalled,
		},
		{
			Name:        "tmux",
			Render:      tools.RenderItem,
			Disabled:    installedSoftwareList.Tmux,
			Run:         tools.InstallTmux,
			GetDisabled: tools.IsTmuxInstalled,
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
				if software.Run == nil {
					NotSupported(choice)
				} else if !software.Disabled {
					software.Run()
					items[i].UpdateDisabled()
				}

				break
			}
		}
	}
}
