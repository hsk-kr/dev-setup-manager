package display

import (
	"github.com/fatih/color"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func DisplayHeader() {
	terminal.ClearConsole()
	print := color.New(color.FgHiCyan).PrintlnFunc()

	print("==============================")
	print("|                            |")
	print("|     Dev Setup Manager      |")
	print("|                     hsk-kr |")
	print("==============================")
}
