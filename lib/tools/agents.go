package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SelectProjectDir uses fzf to let the user pick a git project directory.
// Returns the selected path or an error.
func SelectProjectDir() (string, error) {
	if !ExistCommand("fzf") {
		return "", fmt.Errorf("fzf is required. Install it from the Tools menu")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Find git repos up to 4 levels deep, pipe to fzf
	findCmd := fmt.Sprintf("find %s -maxdepth 4 -type d -name .git 2>/dev/null | sed 's/\\/.git$//' | sort", home)
	cmd := exec.Command("bash", "-c", findCmd+" | fzf --prompt='Select project: '")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("project selection cancelled")
	}

	dir := strings.TrimSpace(string(output))
	if dir == "" {
		return "", fmt.Errorf("no directory selected")
	}

	// Validate .git exists
	gitPath := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		return "", fmt.Errorf("selected directory does not contain a git repository")
	}

	return dir, nil
}
