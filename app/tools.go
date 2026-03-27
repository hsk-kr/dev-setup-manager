package app

import (
	"fmt"
	"sync"

	"github.com/hsk-kr/licokit/lib/config"
	"github.com/hsk-kr/licokit/lib/display"
	"github.com/hsk-kr/licokit/lib/styles"
	"github.com/hsk-kr/licokit/lib/terminal"
	"github.com/hsk-kr/licokit/lib/tools"
)

func Tools(cfg *config.Config) {
	fmt.Println(styles.LoadingText.Render("Reading installed software..."))

	items := make([]terminal.SelectItem, len(cfg.Tools))
	for i, t := range cfg.Tools {
		items[i] = terminal.SelectItem{Name: t.Name}
	}

	// Initialize item properties
	var wg sync.WaitGroup
	for i := range items {
		wg.Add(1)

		go func() {
			defer wg.Done()
			toolCfg := cfg.Tools[i]
			items[i].Render = tools.RenderItem
			items[i].GetDisabled = func() bool {
				installed, _ := tools.IsInstalled(toolCfg)
				return installed
			}
			items[i].Disabled = items[i].GetDisabled()
			items[i].Run = func() {
				if err := tools.Install(toolCfg); err != nil {
					tools.WarningMessage(err.Error())
				}
			}
		}()
	}
	wg.Wait()

	display.DisplayHeader(true)
	fmt.Println(styles.SectionTitle.Render("Select a tool to install"))

	for {
		choice, err := terminal.Select(items)

		if err != nil {
			return
		}

		for i, software := range items {
			if software.Name == choice {
				if software.Run == nil {
					NotSupported(choice)
				} else if !software.Disabled {
					software.Run()
					items[i].UpdateDisabled()
				}

				break
			}
		}
	}
}
