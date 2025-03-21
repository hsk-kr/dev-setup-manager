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
		case "Dotfiles":
			tools.SetupDotfiles()
			tools.SuccessMessage("Dotfiles setup complete.\n\n• Previous dotfiles have been deleted\n• New dotfiles installed in ~/dev-setup-manager/dotfiles\n• To apply zsh changes, run: source ~/.zshrc")
		case "Guide":
			Guide()
		default:
			NotSupported(choice)
		}

		display.DisplayHeader(false)
	}
}
