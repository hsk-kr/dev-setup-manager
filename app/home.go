package app

import (
	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func Home() {
	items := []string{"Tools", "Shell", "Dotfiles"}

	display.DisplayHeader()

	for {
		choice, err := terminal.Select(items)

		if err != nil {
			return
		}

		switch choice {
		case "Tools":
			Tools()
		default:
			NotSupported(choice)
		}
	}
}
