package app

import (
	"fmt"

	"github.com/hsk-kr/licokit/lib/config"
	"github.com/hsk-kr/licokit/lib/display"
	"github.com/hsk-kr/licokit/lib/styles"
	"github.com/hsk-kr/licokit/lib/terminal"
	"github.com/hsk-kr/licokit/lib/tools"
)

func Agents() {
	agentsCfg, err := config.LoadAgents()
	if err != nil {
		tools.WarningMessage(fmt.Sprintf("Failed to load agents config: %v", err))
		return
	}

	items := make([]terminal.SelectItem, len(agentsCfg.Agents))
	for i, agent := range agentsCfg.Agents {
		agent := agent // capture loop variable
		items[i] = terminal.SelectItem{
			Name: agent.Name,
			Render: func(name string, disabled bool) {
				fmt.Printf("%s %s",
					styles.ItemName.Render(name),
					styles.ItemNameDisabled.Render("— "+agent.Description),
				)
			},
		}
	}

	display.DisplayHeader(true)
	fmt.Println(styles.SectionTitle.Render("Select agent roles (Space to toggle, Enter to confirm)"))

	selected, err := terminal.MultiSelect(items)
	if err != nil {
		return
	}

	// Select project directory via fzf
	terminal.ShowCursor()
	terminal.ClearConsole()
	fmt.Println(styles.SectionTitle.Render("Select project directory"))

	projectDir, err := tools.SelectProjectDir()
	if err != nil {
		tools.WarningMessage(err.Error())
		terminal.HideCursor()
		return
	}

	// Build agent panes
	var panes []tools.AgentPane
	for _, name := range selected {
		for _, agent := range agentsCfg.Agents {
			if agent.Name == name {
				prompt, err := config.BuildPrompt(agent)
				if err != nil {
					tools.WarningMessage(fmt.Sprintf("Failed to load prompt for %s: %v", name, err))
					terminal.HideCursor()
					return
				}
				panes = append(panes, tools.AgentPane{
					Name:   name,
					Prompt: prompt,
				})
				break
			}
		}
	}

	// Launch tmux session
	fmt.Println(styles.LoadingText.Render("Launching tmux session..."))
	if err := tools.LaunchAgentSession(projectDir, panes); err != nil {
		tools.WarningMessage(err.Error())
		terminal.HideCursor()
		return
	}
}
