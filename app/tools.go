package app

import (
	"sync"

	"github.com/fatih/color"
	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
	"github.com/hsk-kr/dev-setup-manager/lib/tools"
)

func GetSelectItems() []terminal.SelectItem {
	return []terminal.SelectItem{
		{
			Name: "Homebrew",
		},
		{
			Name: "WezTerm",
		},
		{
			Name: "Neovim",
		},
		{
			Name: "tmux",
		},
		{
			Name: "AeroSpace",
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
			Name: "ripgrep",
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
}

func Tools() {
	loadingPen := color.New(color.FgHiGreen).PrintfFunc()

	loadingPen("Reading installed software...\n")

	items := GetSelectItems()

	// Initilaize item properties
	var wg sync.WaitGroup
	for i := range items {
		wg.Add(1)

		go func() {
			defer wg.Done()
			items[i].Render = tools.RenderItem
			items[i].GetDisabled = func() bool {
				installed, _ := tools.IsInstalled(items[i].Name)
				return installed
			}
			items[i].Disabled = items[i].GetDisabled()
			items[i].Run = func() {
				tools.Install(items[i].Name)
			}
		}()
	}
	wg.Wait()

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
