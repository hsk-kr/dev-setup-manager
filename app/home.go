package app

import (
	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
	"github.com/hsk-kr/dev-setup-manager/lib/tools"
)

func Home() {
	items := []terminal.SelectItem{{
		Name: "Tools",
	}, {
		Name: "Zshrc",
	}, {
		Name: "Dotfiles",
	}, {
		Name: "Guide",
	},
	}

	display.DisplayHeader(true)

	for {
		choice, err := terminal.Select(items)

		if err != nil {
			return
		}

		switch choice {
		case "Tools":
			Tools()
			display.DisplayHeader(true)
			continue
		case "Zshrc":
			tools.SetupZshrc()
			tools.SuccessMessage("Done .zshrc setup.\nCheck ~/dev-setup-manager/dev.zsh file.\nThe setup could have duplicated zsh as the setup's updated.")
		case "Dotfiles":
			tools.SetupDotfiles()
			tools.SuccessMessage("Done dotfiles setup.\nCheck ~/dev-setup-manager/dotfiles.\nIf you already had the dotfiles before, it would be reinstalled.")
		case "Guide":
			Guide()
		default:
			NotSupported(choice)
		}

		display.DisplayHeader(false)
	}
}
