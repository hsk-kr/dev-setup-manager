package terminal

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-tty"
)

func MoveCursor(x, relativeY int) {
	if relativeY < 0 {
		fmt.Printf("\033[%dA\033[%dG", -relativeY, x)
	} else {
		fmt.Printf("\033[%dB\033[%dG", relativeY, x)
	}
}

/*
*

	Display items and returns the index of the item, which starts from zero
*/
func Select(items []string) int {
	print := color.New(color.FgWhite).PrintfFunc()
	cursor := color.New(color.FgGreen).Add(color.Bold).PrintFunc()
	itemLength := len(items)
	currentIndex := 0

	for _, item := range items {
		print("   %s\n", item)
	}

	eraseCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		cursor(" ")
		MoveCursor(1, -(-itemLength+currentIndex)+itemLength)
	}

	drawCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		cursor(">")
		MoveCursor(1, -(-itemLength+currentIndex)+itemLength)
	}

	drawCurrentCursor()

	t, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	for {
		r, err := t.ReadRune()
		if err != nil {
			panic(err)
		}

		switch r {
		case 'j', 'J', 'h', 'H':
			eraseCurrentCursor()
			if currentIndex < itemLength-1 {
				currentIndex++
			}
			drawCurrentCursor()
		case 'k', 'K', 'l', 'L':
			eraseCurrentCursor()
			if currentIndex > 0 {
				currentIndex--
			}
			drawCurrentCursor()
		}

		if r == '\r' || r == '\n' {
			break
		}
	}

	return currentIndex
}
