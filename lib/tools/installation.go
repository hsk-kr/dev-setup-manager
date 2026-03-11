package tools

import (
	"fmt"
	"os"

	"github.com/hsk-kr/dev-setup-manager/lib/config"
	"github.com/hsk-kr/dev-setup-manager/lib/spinner"
)

func Install(tool config.ToolConfig) error {
	switch tool.InstallType {
	case "manual":
		SuccessMessage(tool.ManualMessage)
		return nil
	case "brew":
		sp := spinner.New(fmt.Sprintf("Installing %s...", tool.Name))
		sp.Start()
		err := ExecCommandQuiet("brew", "install", tool.BrewPackage())
		sp.Stop()
		if err != nil {
			return err
		}
	case "cask":
		sp := spinner.New(fmt.Sprintf("Installing %s...", tool.Name))
		sp.Start()
		err := ExecCommandQuiet("brew", "install", "--cask", tool.BrewPackage())
		sp.Stop()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown install type: %s", tool.InstallType)
	}

	// Create post-install directories
	for _, dir := range tool.PostInstallDirs {
		expanded := config.ExpandPath(dir)
		if err := os.MkdirAll(expanded, 0755); err != nil {
			return err
		}
	}

	// Add zsh source if specified
	if tool.ZshSource != "" {
		if err := AddZshSource(tool.ZshSource); err != nil {
			return err
		}
	}

	// Show post-install warning if specified
	if tool.PostInstallWarning != "" {
		WarningMessage(tool.PostInstallWarning)
	}

	return nil
}

func IsInstalled(tool config.ToolConfig) (bool, error) {
	switch tool.DetectType {
	case "command":
		return ExistCommand(tool.DetectValue), nil
	case "application":
		return ExistApplication(tool.DetectValue), nil
	case "brew_package":
		return ExistBrewPackage(tool.DetectValue), nil
	default:
		return false, fmt.Errorf("unknown detect type: %s", tool.DetectType)
	}
}
