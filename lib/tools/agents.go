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

// AgentPane holds the info needed to launch one pane.
type AgentPane struct {
	Name   string
	Prompt string // full prompt text, empty for Free role
}

// LaunchAgentSession creates a tmux session with panes for each agent.
// Each pane runs claude with the given prompt in the project directory.
func LaunchAgentSession(projectDir string, agents []AgentPane) error {
	if !ExistCommand("tmux") {
		return fmt.Errorf("tmux is required. Install it from the Tools menu")
	}
	if !ExistCommand("claude") {
		return fmt.Errorf("claude is required. Install it from the Tools menu")
	}

	sessionName := generateSessionName(projectDir)
	agentCount := len(agents)

	// Create session with first agent
	firstCmd := buildClaudeCommand(projectDir, agents[0])
	err := ExecCommandQuiet("tmux", "new-session", "-d", "-s", sessionName, "-x", "200", "-y", "50", firstCmd)
	if err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Rename first pane
	_ = ExecCommandQuiet("tmux", "select-pane", "-t", sessionName+":0.0", "-T", agents[0].Name)

	// Create remaining panes
	for i := 1; i < agentCount; i++ {
		cmd := buildClaudeCommand(projectDir, agents[i])
		err := ExecCommandQuiet("tmux", "split-window", "-t", sessionName, cmd)
		if err != nil {
			return fmt.Errorf("failed to create pane for %s: %w", agents[i].Name, err)
		}
		_ = ExecCommandQuiet("tmux", "select-pane", "-t", sessionName, "-T", agents[i].Name)
	}

	// Apply layout
	applyLayout(sessionName, agentCount)

	// Enable pane border status to show role names
	_ = ExecCommandQuiet("tmux", "set-option", "-t", sessionName, "pane-border-status", "top")
	_ = ExecCommandQuiet("tmux", "set-option", "-t", sessionName, "pane-border-format", "#{pane_title}")

	// Attach to session
	attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
	attachCmd.Stdin = os.Stdin
	attachCmd.Stdout = os.Stdout
	attachCmd.Stderr = os.Stderr
	return attachCmd.Run()
}

func buildClaudeCommand(projectDir string, agent AgentPane) string {
	if agent.Prompt == "" {
		return fmt.Sprintf("cd %s && claude", shellEscape(projectDir))
	}
	return fmt.Sprintf("cd %s && claude --append-system-prompt %s", shellEscape(projectDir), shellEscape(agent.Prompt))
}

func shellEscape(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}

func generateSessionName(projectDir string) string {
	base := filepath.Base(projectDir) + "-claude-agents"

	// Check if session already exists, append number if so
	name := base
	for i := 2; ; i++ {
		err := ExecCommandQuiet("tmux", "has-session", "-t", name)
		if err != nil {
			// Session doesn't exist, use this name
			return name
		}
		name = fmt.Sprintf("%s-%d", base, i)
	}
}

func applyLayout(sessionName string, count int) {
	// Layout: max 2 rows, expand horizontally first
	// 1: fullscreen, 2: 2 cols, 3: 3 cols, 4: 2x2, 5: 3+2, 6: 3x2, 7: 4+3
	if count <= 3 {
		// Single row: even-horizontal
		_ = ExecCommandQuiet("tmux", "select-layout", "-t", sessionName, "even-horizontal")
	} else {
		// Two rows: tiled gives us a grid
		_ = ExecCommandQuiet("tmux", "select-layout", "-t", sessionName, "tiled")
	}
}
