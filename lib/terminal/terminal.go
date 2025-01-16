package terminal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/mattn/go-tty"
)

type SelectItem struct {
	Name        string
	Render      func(name string, disabled bool)
	Run         func()
	GetDisabled func() bool
	Disabled    bool
}

func (si *SelectItem) UpdateDisabled() {
	if si.GetDisabled != nil {
		si.Disabled = si.GetDisabled()
	}
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ClearConsole() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func MoveCursor(x, relativeY int) {
	if relativeY < 0 {
		fmt.Printf("\033[%dA\033[%dG", -relativeY, x)
	} else {
		fmt.Printf("\033[%dB\033[%dG", relativeY, x)
	}
}

/*
Display items and returns the name of the item
If user presses esc, it returns empty string with an error
*/
func Select(items []SelectItem) (string, error) {
	print := color.New(color.FgWhite).PrintfFunc()
	cursor := color.New(color.FgGreen).Add(color.Bold).PrintFunc()
	itemLength := len(items)
	currentIndex := 0

	for _, item := range items {
		print("   ")

		if item.Render != nil {
			item.Render(item.Name, item.Disabled)
		} else {
			print(item.Name)
		}

		print("\n")
	}

	eraseCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		cursor(" ")
		MoveCursor(1, -(-itemLength + currentIndex))
	}

	drawCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		cursor(">")
		MoveCursor(1, -(-itemLength + currentIndex))
	}

	drawCurrentCursor()

	t, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	for {
		r, _ := t.ReadRune()

		switch r {
		case '\x1b':
			return "", errors.New("Escape")
		case 'j', 'J', 'h', 'H':
			if currentIndex >= itemLength-1 {
				break
			}
			eraseCurrentCursor()
			currentIndex++
			drawCurrentCursor()
		case 'k', 'K', 'l', 'L':
			if currentIndex <= 0 {
				break
			}
			eraseCurrentCursor()
			currentIndex--
			drawCurrentCursor()
		}

		if (r == '\r' || r == '\n') && !items[currentIndex].Disabled {
			break
		}
	}

	return items[currentIndex].Name, nil
}
