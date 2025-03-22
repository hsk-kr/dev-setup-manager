package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	ls := exec.Command("ls", "/Applications")
	grep := exec.Command("grep", appName)

	pipe, err := ls.StdoutPipe()
	defer pipe.Close()

	if err != nil {
		panic(err)
	}

	grep.Stdin = pipe

	err = ls.Start()

	if err != nil {
		panic(err)
	}

	output, err := grep.Output()

	if err != nil {
		return false
	}

	return len(string(output)) >= len(appName)
}

func ExecCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
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

	devSetupManagerHomePath := homePath + "/dev-setup-manager"
	devSetupManagerZshPath := devSetupManagerHomePath + "/dev.zsh"
	zshrcPath := homePath + "/.zshrc"
	ExecCommand("mkdir", "-p", devSetupManagerHomePath)

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
		return nil
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
	f.WriteString(content)
	return nil
}

func WarningMessage(message string) {
	print := color.New(color.FgRed).PrintlnFunc()
	print(message)
}

func SuccessMessage(message string) {
	print := color.New(color.FgGreen).PrintlnFunc()
	print(message)
}
