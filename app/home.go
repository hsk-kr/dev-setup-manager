package app

import (
	"fmt"

	"github.com/hsk-kr/dev-setup-manager/lib/display"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func Home() {
	items := []string{"Tools", "Shell", "Dotfiles"}

	display.DisplayHeader()
	index := terminal.Select(items)

	fmt.Printf("Selected %d index\n", index)
}
