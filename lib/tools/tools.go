package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hsk-kr/licokit/lib/styles"
)

func RenderItem(name string, disabled bool) {
	if disabled {
		fmt.Print(styles.ItemNameDisabled.Render(name))
		fmt.Print(styles.StatusInstalled.Render(" ✓ Installed"))
	} else {
		fmt.Print(styles.ItemName.Render(name))
		fmt.Print(styles.StatusNotInstalled.Render(" ✗ Not Installed"))
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

// ExecCommandQuiet runs a command without printing stdout/stderr.
// Used when a spinner is showing progress instead.
func ExecCommandQuiet(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w", command, err)
	}

	return nil
}

/*
Add the source to ~/licokit/dev.zsh

It will insert the line to import dev.zsh from .zshrc if it's not setup.

If the source exists in the dev.zsh, it will be ignored.
*/
func AddZshSource(source string) error {
	homePath, err := os.UserHomeDir()

	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	licokitHomePath := filepath.Join(homePath, "licokit")
	licokitZshPath := filepath.Join(licokitHomePath, "dev.zsh")
	zshrcPath := filepath.Join(homePath, ".zshrc")
	if err := ExecCommand("mkdir", "-p", licokitHomePath); err != nil {
		return err
	}

	// add source into dev.zsh. if it's already setup, it does nothing
	addSourceToDevZsh := func() error {
		exist, err := existFile(licokitZshPath)

		if err != nil {
			return err
		}

		if exist {
			contains, err := containInFile(licokitZshPath, source)

			if err != nil {
				return err
			}

			if contains {
				return nil
			}
		}

		err = appendFile(licokitZshPath, fmt.Sprintf("\n%s", source))
		return err
	}

	// add source dev zsh into .zshrc. if it's already setup, it does nothing
	addDevZshToZshrc := func() error {
		exist, err := existFile(zshrcPath)

		if err != nil {
			return err
		}

		if exist {
			contains, err := containInFile(zshrcPath, licokitZshPath)

			if err != nil {
				return err
			}

			if contains {
				return nil
			}
		}

		err = appendFile(zshrcPath, fmt.Sprintf("\nsource %s", licokitZshPath))
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
	fmt.Println(styles.WarningBox.Render("⚠ " + message))
}

func SuccessMessage(message string) {
	fmt.Println(styles.SuccessBox.Render("✓ " + message))
}
