package display

import (
	"fmt"

	"github.com/hsk-kr/licokit/lib/styles"
	"github.com/hsk-kr/licokit/lib/terminal"
)

func DisplayHeader(clearConsole bool) {
	if clearConsole {
		terminal.ClearConsole()
	}
	header := styles.HeaderBox.Render("licokit\n          lico's dev loadout")
	fmt.Println(header)
}
