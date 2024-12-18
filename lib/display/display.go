package display

import (
	"github.com/fatih/color"
)

func DisplayHeader() {
	print := color.New(color.FgHiCyan).PrintlnFunc()

	print("==============================")
	print("|                            |")
	print("|     Dev Setup Manager      |")
	print("|                     hsk-kr |")
	print("==============================")
}

func DisplayStep(stepName string) {
	print := color.New(color.FgHiBlue).PrintfFunc()

	print("%s\n", stepName)
	print("---\n")
}
