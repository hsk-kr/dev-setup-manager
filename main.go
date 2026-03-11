package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hsk-kr/dev-setup-manager/app"
	"github.com/hsk-kr/dev-setup-manager/lib/config"
	"github.com/hsk-kr/dev-setup-manager/lib/terminal"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	terminal.HideCursor()
	defer terminal.ShowCursor()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		terminal.ShowCursor()
		os.Exit(0)
	}()

	app.Home(cfg)
}
