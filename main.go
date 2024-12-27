package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hsk-kr/dev-setup-manager/app"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func main() {
	terminal.HideCursor()
	defer terminal.ShowCursor()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		terminal.ShowCursor()
		os.Exit(0)
	}()

	app.Home()
}
