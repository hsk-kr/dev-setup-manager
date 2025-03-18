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
	},
	}

	display.DisplayHeader()

	for {
		choice, err := terminal.Select(items)

		if err != nil {
			return
		}

		switch choice {
		case "Tools":
			Tools()
		case "Zshrc":
			tools.SetupZshrc()
		default:
			NotSupported(choice)
			continue
		}

		display.DisplayHeader()

		switch choice {
		case "Zshrc":
			tools.SuccessMessage("Done .zshrc setup.\nCheck ~/dev-setup-manager/dev.zsh file.\nThe setup could have duplicated zsh as the setup's updated.")
		}
	}
}
