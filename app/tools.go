package app

import (
	"github.com/fatih/color"
	"github.com/hsk-kr/dev-setup-manager/app/tools"
	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func Tools() {
	items := []string{"Homebrew", "WezTerm", "Neovim", "Aerospace", "Homerow", "Karabiner Elements", "Snipaste", "fzf", "zsh-vi-mode"}

	display.DisplayHeader()
	print := color.New(color.FgGreen).PrintlnFunc()
	print("Select item you install")

	choice, err := terminal.Select(items)

	if err != nil {
		return
	}

	switch choice {
	case "Homebrew":
		tools.Homebrew()
	default:
		NotSupported(choice)
	}
}
