package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func RenderItem(name string, disabled bool) {
	installedPen := color.New(color.FgHiGreen).PrintFunc()
	notInstalledPen := color.New(color.FgRed).PrintFunc()
	softwareNamePen := color.New(color.FgWhite).PrintFunc()

	softwareNamePen(name)
	if disabled {
		installedPen(" - Installed")
	} else {
		notInstalledPen(" - Not Installed")
	}
}

func ExistCommand(cmd string) bool {
	_, err := exec.Command("zsh", "-l", "-c", fmt.Sprintf("which %s", cmd)).Output()

	if err != nil {
		return false
	}

	return true
}

func ExistBrewPackage(packageName string) bool {
	_, err := exec.Command("brew", "list", packageName).Output()

	if err != nil {
		return false
	}

	return true
}

func ExistApplication(appName string) bool {
	appPath := filepath.Join("/Applications", appName)
	_, err := os.Stat(appPath)
	return err == nil
}

func ExecCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w", command, err)
	}

	return nil
}

/*
Add the source to ~/dev-setup-manager/dev.zsh

It will insert the line to import dev.zsh from .zshrc if it's not setup.

If the source exists in the dev.zsh, it will be ignored.
*/
func AddZshSource(source string) error {
	homePath, err := os.UserHomeDir()

	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	devSetupManagerHomePath := filepath.Join(homePath, "dev-setup-manager")
	devSetupManagerZshPath := filepath.Join(devSetupManagerHomePath, "dev.zsh")
	zshrcPath := filepath.Join(homePath, ".zshrc")
	if err := ExecCommand("mkdir", "-p", devSetupManagerHomePath); err != nil {
		return err
	}

	// add source into dev.zsh. if it's already setup, it does nothing
	addSourceToDevZsh := func() error {
		exist, err := existFile(devSetupManagerZshPath)

		if err != nil {
			return err
		}

		if exist {
			contains, err := containInFile(devSetupManagerZshPath, source)

			if err != nil {
				return err
			}

			if contains {
				return nil
			}
		}

		err = appendFile(devSetupManagerZshPath, fmt.Sprintf("\n%s", source))
		return err
	}

	// add source dev zsh into .zshrc. if it's already setup, it does nothing
	addDevZshToZshrc := func() error {
		exist, err := existFile(zshrcPath)

		if err != nil {
			return err
		}

		if exist {
			contains, err := containInFile(zshrcPath, devSetupManagerZshPath)

			if err != nil {
				return err
			}

			if contains {
				return nil
			}
		}

		err = appendFile(zshrcPath, fmt.Sprintf("\nsource %s", devSetupManagerZshPath))
		return err
	}

	err = addSourceToDevZsh()
	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	return addDevZshToZshrc()
}

func existFile(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

/*
Returns if the str exists in the content of the file or not
*/
func containInFile(path, str string) (bool, error) {
	bContent, err := os.ReadFile(path)

	if err != nil {
		return false, err
	}

	return strings.Contains(string(bContent), str), nil
}

func appendFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(content)
	return err
}

func WarningMessage(message string) {
	print := color.New(color.FgRed).PrintlnFunc()
	print(message)
}

func SuccessMessage(message string) {
	print := color.New(color.FgGreen).PrintlnFunc()
	print(message)
}
