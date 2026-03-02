package display

import (
	"fmt"

	"github.com/hsk-kr/dev-setup-manager/lib/styles"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func DisplayHeader(clearConsole bool) {
	if clearConsole {
		terminal.ClearConsole()
	}
	header := styles.HeaderBox.Render("Dev Setup Manager\n                    hsk-kr")
	fmt.Println(header)
}
