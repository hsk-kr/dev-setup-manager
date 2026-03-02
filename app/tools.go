package app

import (
	"fmt"
	"sync"

	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/styles"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
	"github.com/hsk-kr/dev-setup-manager/lib/tools"
)

func GetSelectItems() []terminal.SelectItem {
	return []terminal.SelectItem{
		{
			Name: "Homebrew",
		},
		{
			Name: "Git",
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
			Name: "go",
		},
		{
			Name: "nvm",
		},
		{
			Name: "btop",
		},
	}
}

func Tools() {
	fmt.Println(styles.LoadingText.Render("Reading installed software..."))

	items := GetSelectItems()

	// Initialize item properties
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
				if err := tools.Install(items[i].Name); err != nil {
					tools.WarningMessage(err.Error())
				}
			}
		}()
	}
	wg.Wait()

	display.DisplayHeader(true)
	fmt.Println(styles.SectionTitle.Render("Select a tool to install"))

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
