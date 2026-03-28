# Claude Agents Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a "Claude Agents" menu to licokit that launches multiple proactive Claude Code sessions in tmux, each with a role-specific system prompt.

**Architecture:** New `agents.go` config module loads embedded `agents.yaml` + `prompts/*.md` with user override at `~/.config/licokit/`. A new `MultiSelect()` in the terminal lib lets the user pick roles. The `app/agents.go` flow shells out to `fzf` for directory selection and `tmux` for session/pane creation, running `claude --append-system-prompt` in each pane.

**Tech Stack:** Go 1.23, lipgloss (styling), go-tty (input), fzf (directory picker), tmux (session management), claude CLI

---

### Task 1: Agent Config Types and Loading

**Files:**
- Create: `lib/config/agents.go`
- Create: `lib/config/agents.yaml`
- Create: `lib/config/agents_test.go`

- [ ] **Step 1: Write the failing test for AgentConfig parsing**

```go
// lib/config/agents_test.go
package config

import (
	"testing"
)

func TestParseAgentsYAML(t *testing.T) {
	yaml := `
agents:
  - name: Developer
    description: "Writes code"
    prompt_file: developer.md
  - name: Free
    description: "Plain session"
    prompt_file: ""
`
	cfg, err := parseAgentsYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("parseAgentsYAML error: %v", err)
	}

	if len(cfg.Agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(cfg.Agents))
	}

	if cfg.Agents[0].Name != "Developer" {
		t.Errorf("agent 0 name = %q, want Developer", cfg.Agents[0].Name)
	}
	if cfg.Agents[0].Description != "Writes code" {
		t.Errorf("agent 0 description = %q, want 'Writes code'", cfg.Agents[0].Description)
	}
	if cfg.Agents[0].PromptFile != "developer.md" {
		t.Errorf("agent 0 prompt_file = %q, want developer.md", cfg.Agents[0].PromptFile)
	}
	if cfg.Agents[1].PromptFile != "" {
		t.Errorf("agent 1 prompt_file = %q, want empty", cfg.Agents[1].PromptFile)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -run TestParseAgentsYAML -v`
Expected: FAIL — `parseAgentsYAML` undefined

- [ ] **Step 3: Write the types and parser**

```go
// lib/config/agents.go
package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed agents.yaml
var defaultAgentsConfig embed.FS

//go:embed prompts/*
var embeddedPrompts embed.FS

type AgentConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	PromptFile  string `yaml:"prompt_file"`
}

type AgentsConfig struct {
	Agents []AgentConfig `yaml:"agents"`
}

func parseAgentsYAML(data []byte) (*AgentsConfig, error) {
	var cfg AgentsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// LoadAgents reads agents config from user path or falls back to embedded default.
func LoadAgents() (*AgentsConfig, error) {
	homePath, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(homePath, ".config", "licokit", "agents.yaml")
		if data, err := os.ReadFile(userConfigPath); err == nil {
			cfg, err := parseAgentsYAML(data)
			if err != nil {
				return nil, fmt.Errorf("failed to parse user agents config: %w", err)
			}
			return cfg, nil
		}
	}

	data, err := defaultAgentsConfig.ReadFile("agents.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded agents config: %w", err)
	}

	return parseAgentsYAML(data)
}

// LoadPrompt resolves a prompt file: user override at ~/.config/licokit/prompts/{file},
// then embedded default. Returns empty string if promptFile is empty.
func LoadPrompt(promptFile string) (string, error) {
	if promptFile == "" {
		return "", nil
	}

	// Check user override
	homePath, err := os.UserHomeDir()
	if err == nil {
		userPromptPath := filepath.Join(homePath, ".config", "licokit", "prompts", promptFile)
		if data, err := os.ReadFile(userPromptPath); err == nil {
			return string(data), nil
		}
	}

	// Fall back to embedded
	data, err := embeddedPrompts.ReadFile(filepath.Join("prompts", promptFile))
	if err != nil {
		return "", fmt.Errorf("failed to read embedded prompt %q: %w", promptFile, err)
	}

	return string(data), nil
}

// BuildPrompt loads common prefix + role prompt and returns the full prompt string.
func BuildPrompt(agent AgentConfig) (string, error) {
	if agent.PromptFile == "" {
		return "", nil
	}

	common, err := LoadPrompt("common.md")
	if err != nil {
		return "", fmt.Errorf("failed to load common prompt: %w", err)
	}

	role, err := LoadPrompt(agent.PromptFile)
	if err != nil {
		return "", err
	}

	return common + "\n\n" + role, nil
}
```

- [ ] **Step 4: Create the default agents.yaml**

```yaml
# lib/config/agents.yaml
agents:
  - name: Issuer
    description: "Planner and issue creator"
    prompt_file: issuer.md

  - name: Developer
    description: "Writes code"
    prompt_file: developer.md

  - name: Tester
    description: "Tests the service"
    prompt_file: tester.md

  - name: Code Reviewer
    description: "Reviews code quality"
    prompt_file: code-reviewer.md

  - name: Security Reviewer
    description: "Audits security"
    prompt_file: security-reviewer.md

  - name: UI/UX Designer
    description: "Reviews UI/UX"
    prompt_file: ui-ux-designer.md

  - name: Free
    description: "Plain session"
    prompt_file: ""
```

- [ ] **Step 5: Create placeholder prompt files so embed compiles**

Create these files under `lib/config/prompts/` with placeholder content `TODO` (they'll be filled in Task 7):

- `common.md`
- `issuer.md`
- `developer.md`
- `tester.md`
- `code-reviewer.md`
- `security-reviewer.md`
- `ui-ux-designer.md`

- [ ] **Step 6: Run test to verify it passes**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -run TestParseAgentsYAML -v`
Expected: PASS

- [ ] **Step 7: Write test for LoadAgents with embedded default**

```go
// append to lib/config/agents_test.go

func TestLoadAgents_DefaultConfig(t *testing.T) {
	cfg, err := LoadAgents()
	if err != nil {
		t.Fatalf("LoadAgents() error: %v", err)
	}

	if len(cfg.Agents) != 7 {
		t.Errorf("expected 7 agents, got %d", len(cfg.Agents))
	}

	// Verify all agents have required fields
	for _, agent := range cfg.Agents {
		if agent.Name == "" {
			t.Error("agent name should not be empty")
		}
		if agent.Description == "" {
			t.Errorf("agent %q: description should not be empty", agent.Name)
		}
	}

	// Free agent should have empty prompt_file
	free := cfg.Agents[len(cfg.Agents)-1]
	if free.Name != "Free" {
		t.Errorf("last agent = %q, want Free", free.Name)
	}
	if free.PromptFile != "" {
		t.Errorf("Free agent prompt_file = %q, want empty", free.PromptFile)
	}
}
```

- [ ] **Step 8: Run test to verify it passes**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -run TestLoadAgents -v`
Expected: PASS

- [ ] **Step 9: Write test for LoadPrompt**

```go
// append to lib/config/agents_test.go

func TestLoadPrompt_EmptyFile(t *testing.T) {
	prompt, err := LoadPrompt("")
	if err != nil {
		t.Fatalf("LoadPrompt(\"\") error: %v", err)
	}
	if prompt != "" {
		t.Errorf("expected empty prompt, got %q", prompt)
	}
}

func TestLoadPrompt_EmbeddedFile(t *testing.T) {
	prompt, err := LoadPrompt("common.md")
	if err != nil {
		t.Fatalf("LoadPrompt(\"common.md\") error: %v", err)
	}
	if prompt == "" {
		t.Error("expected non-empty prompt for common.md")
	}
}

func TestLoadPrompt_NonexistentFile(t *testing.T) {
	_, err := LoadPrompt("nonexistent.md")
	if err == nil {
		t.Error("expected error for nonexistent prompt file")
	}
}
```

- [ ] **Step 10: Run tests to verify they pass**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -run TestLoadPrompt -v`
Expected: PASS

- [ ] **Step 11: Write test for BuildPrompt**

```go
// append to lib/config/agents_test.go

func TestBuildPrompt_FreeAgent(t *testing.T) {
	agent := AgentConfig{Name: "Free", PromptFile: ""}
	prompt, err := BuildPrompt(agent)
	if err != nil {
		t.Fatalf("BuildPrompt error: %v", err)
	}
	if prompt != "" {
		t.Errorf("expected empty prompt for Free agent, got %q", prompt)
	}
}

func TestBuildPrompt_RoleAgent(t *testing.T) {
	agent := AgentConfig{Name: "Issuer", PromptFile: "issuer.md"}
	prompt, err := BuildPrompt(agent)
	if err != nil {
		t.Fatalf("BuildPrompt error: %v", err)
	}
	if prompt == "" {
		t.Error("expected non-empty prompt for Issuer agent")
	}
}
```

- [ ] **Step 12: Run all config tests**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -v`
Expected: ALL PASS

- [ ] **Step 13: Commit**

```bash
git add lib/config/agents.go lib/config/agents.yaml lib/config/agents_test.go lib/config/prompts/
git commit -m "feat: add agent config types, loading, and prompt resolution"
```

---

### Task 2: MultiSelect UI Component

**Files:**
- Modify: `lib/terminal/terminal.go`

- [ ] **Step 1: Write MultiSelect function**

Add to the bottom of `lib/terminal/terminal.go`:

```go
// MultiSelect displays items with checkboxes and returns selected names.
// j/k to move, Space to toggle, Enter to confirm, ESC to cancel.
func MultiSelect(items []SelectItem) ([]string, error) {
	itemLength := len(items)
	currentIndex := 0
	selected := make([]bool, itemLength)

	checkbox := func(checked bool) string {
		if checked {
			return styles.Cursor.Render("☑")
		}
		return "☐"
	}

	// Initial render
	for i, item := range items {
		fmt.Printf("   %s ", checkbox(selected[i]))
		if item.Render != nil {
			item.Render(item.Name, item.Disabled)
		} else {
			fmt.Printf("%s", item.Name)
		}
		fmt.Printf("\n")
	}

	eraseCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		fmt.Print(" ")
		MoveCursor(1, -(-itemLength + currentIndex))
	}

	drawCurrentCursor := func() {
		MoveCursor(1, -itemLength+currentIndex)
		fmt.Print(styles.Cursor.Render("❯"))
		MoveCursor(1, -(-itemLength + currentIndex))
	}

	redrawCheckbox := func(index int) {
		MoveCursor(4, -itemLength+index)
		fmt.Print(checkbox(selected[index]))
		MoveCursor(1, -(-itemLength + index))
	}

	drawCurrentCursor()

	t, err := tty.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open TTY: %w", err)
	}
	defer t.Close()

	for {
		r, err := t.ReadRune()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		switch r {
		case '\x1b':
			return nil, errors.New("Escape")
		case 'j', 'J', 'h', 'H':
			if currentIndex >= itemLength-1 {
				break
			}
			eraseCurrentCursor()
			currentIndex++
			drawCurrentCursor()
		case 'k', 'K', 'l', 'L':
			if currentIndex <= 0 {
				break
			}
			eraseCurrentCursor()
			currentIndex--
			drawCurrentCursor()
		case ' ':
			selected[currentIndex] = !selected[currentIndex]
			redrawCheckbox(currentIndex)
		case '\r', '\n':
			var result []string
			for i, item := range items {
				if selected[i] {
					result = append(result, item.Name)
				}
			}
			if len(result) == 0 {
				break // don't confirm with nothing selected
			}
			return result, nil
		}
	}
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/hsk.coder/dev/licokit && go build ./...`
Expected: Success (no errors)

- [ ] **Step 3: Run existing tests to ensure nothing broke**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./... -v`
Expected: ALL PASS

- [ ] **Step 4: Commit**

```bash
git add lib/terminal/terminal.go
git commit -m "feat: add MultiSelect component with checkbox toggling"
```

---

### Task 3: Styles for Multi-Select

**Files:**
- Modify: `lib/styles/styles.go`

- [ ] **Step 1: Add checkbox styles**

Add to `lib/styles/styles.go`:

```go
	Checkbox = lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // gray

	CheckboxSelected = lipgloss.NewStyle().
		Foreground(lipgloss.Color("2")). // green
		Bold(true)
```

- [ ] **Step 2: Update MultiSelect to use the new styles**

In `lib/terminal/terminal.go`, update the `checkbox` closure inside `MultiSelect`:

```go
	checkbox := func(checked bool) string {
		if checked {
			return styles.CheckboxSelected.Render("☑")
		}
		return styles.Checkbox.Render("☐")
	}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd /Users/hsk.coder/dev/licokit && go build ./...`
Expected: Success

- [ ] **Step 4: Commit**

```bash
git add lib/styles/styles.go lib/terminal/terminal.go
git commit -m "feat: add checkbox styles for multi-select"
```

---

### Task 4: Agents Menu Flow — FZF Directory Picker

**Files:**
- Create: `lib/tools/agents.go`

- [ ] **Step 1: Write the fzf directory picker function**

```go
// lib/tools/agents.go
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
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/hsk.coder/dev/licokit && go build ./...`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add lib/tools/agents.go
git commit -m "feat: add fzf-based project directory picker"
```

---

### Task 5: Agents Menu Flow — tmux Launcher

**Files:**
- Modify: `lib/tools/agents.go`

- [ ] **Step 1: Write the tmux session launcher function**

Append to `lib/tools/agents.go`:

```go
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
	_ = ExecCommandQuiet("tmux", "send-prefix", "-t", sessionName)
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

// AgentPane holds the info needed to launch one pane.
type AgentPane struct {
	Name   string
	Prompt string // full prompt text, empty for Free role
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
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /Users/hsk.coder/dev/licokit && go build ./...`
Expected: Success

- [ ] **Step 3: Commit**

```bash
git add lib/tools/agents.go
git commit -m "feat: add tmux session launcher for claude agents"
```

---

### Task 6: Agents App Screen

**Files:**
- Create: `app/agents.go`
- Modify: `app/home.go`

- [ ] **Step 1: Create the agents app screen**

```go
// app/agents.go
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
```

- [ ] **Step 2: Add "Claude Agents" to the home menu**

In `app/home.go`, update the items slice to add the new menu item. Replace:

```go
	items := []terminal.SelectItem{{
		Name: "Tools",
	}, {
		Name: "Dotfiles",
	}, {
		Name: "Guide",
	},
	}
```

With:

```go
	items := []terminal.SelectItem{{
		Name: "Tools",
	}, {
		Name: "Dotfiles",
	}, {
		Name: "Claude Agents",
	}, {
		Name: "Guide",
	},
	}
```

- [ ] **Step 3: Add the case in the switch**

In `app/home.go`, add this case inside the `switch choice` block, before `case "Guide"`:

```go
		case "Claude Agents":
			Agents()
			return // exit licokit after launching
```

- [ ] **Step 4: Verify it compiles**

Run: `cd /Users/hsk.coder/dev/licokit && go build ./...`
Expected: Success

- [ ] **Step 5: Run all tests**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./... -v`
Expected: ALL PASS

- [ ] **Step 6: Commit**

```bash
git add app/agents.go app/home.go
git commit -m "feat: add Claude Agents menu screen with full launch flow"
```

---

### Task 7: Write Role Prompts

**Files:**
- Modify: `lib/config/prompts/common.md`
- Modify: `lib/config/prompts/issuer.md`
- Modify: `lib/config/prompts/developer.md`
- Modify: `lib/config/prompts/tester.md`
- Modify: `lib/config/prompts/code-reviewer.md`
- Modify: `lib/config/prompts/security-reviewer.md`
- Modify: `lib/config/prompts/ui-ux-designer.md`

- [ ] **Step 1: Write common.md**

```markdown
You are a proactive AI agent with the role of {role_name}.
You do NOT wait for instructions. You continuously work, and when you need a decision or approval, you ask the human.

Rules:
- Always ask before committing code
- Use GitHub Issues as the task tracker for this project
- Never idle — if your current task is done, find the next thing or ask
- When you find something that falls outside your role, create a GitHub issue so the appropriate agent picks it up
- Be concise and direct in your communication
- Read the project's README, CLAUDE.md, and any docs to understand the project context before starting
```

Note: `{role_name}` is a literal placeholder — `BuildPrompt` in `agents.go` needs to replace it. Update `BuildPrompt` to do string replacement:

In `lib/config/agents.go`, update `BuildPrompt`:

```go
func BuildPrompt(agent AgentConfig) (string, error) {
	if agent.PromptFile == "" {
		return "", nil
	}

	common, err := LoadPrompt("common.md")
	if err != nil {
		return "", fmt.Errorf("failed to load common prompt: %w", err)
	}

	role, err := LoadPrompt(agent.PromptFile)
	if err != nil {
		return "", err
	}

	full := common + "\n\n" + role
	full = strings.ReplaceAll(full, "{role_name}", agent.Name)
	return full, nil
}
```

Add `"strings"` to the imports in `agents.go`.

- [ ] **Step 2: Write issuer.md**

```markdown
# Issuer

You are a product-minded collaborator and planner. Your job is to help the human decide what to work on and turn those decisions into well-structured GitHub issues.

## On Launch

- Read the project thoroughly: README, CLAUDE.md, docs folder, recent git history, open/closed issues
- Understand the project's purpose, goals, philosophy, and current state
- Analyze recent issues and commits to understand the trajectory — what's been done, what's in progress, what's missing

## How You Work

- Proactively suggest what to work on next based on the project's goals and current state
- Don't just ask "what's next?" — give context. Explain what areas could benefit from attention, what patterns you've noticed, what the natural next step might be
- Help the human think through priorities and trade-offs
- When the human decides on something, create a well-structured GitHub issue with:
  - Clear title
  - Description with context and motivation
  - Acceptance criteria when applicable
- After creating an issue, analyze what would logically follow and suggest it — be a thought partner, not a ticket machine
- If you notice technical debt, missing features, or improvement opportunities while reviewing the project, surface them as suggestions

## What You Don't Do

- Don't write code
- Don't make big decisions without the human's approval
- Don't create issues without confirming with the human first
```

- [ ] **Step 3: Write developer.md**

```markdown
# Developer

You are a disciplined developer who plans before coding. Your job is to pick up GitHub issues and turn them into working code.

## How You Work

- Check for the oldest open GitHub issue assigned to this project
- Before starting: validate the issue is still relevant — check if it hasn't already been fixed, if the code it references still exists, if the described problem still reproduces
- Ask the human: "Should I start this?" — briefly summarize the issue and your approach
- For non-trivial tasks, plan before coding. Think through the approach, identify affected files, consider edge cases. Use brainstorming/planning tools when the task warrants it.
- Create a branch: `issue-{number}-{short-description}`
- Implement the solution
- Ask before every commit — show what you're about to commit and why
- When done, create a PR that auto-closes the issue (include "Closes #{number}" in the PR body)
- After the PR is created, pick the next oldest open issue and repeat

## When There Are No Issues

- Tell the human there are no open issues
- Let them know you're available and will start when a new issue appears
- Periodically check for new issues

## What You Don't Do

- Don't create issues — that's another role's job
- Don't make architectural decisions without asking
- Don't commit without approval
```

- [ ] **Step 4: Write tester.md**

```markdown
# Tester

You are a quality guardian. Your job is to make sure the service works as intended and that changes don't break existing functionality.

## On Launch

- Discover how the project is tested: find test frameworks, test scripts, e2e setups, CI configuration, package.json scripts, Makefiles, etc.
- Run the existing test suite and report the current state — what passes, what fails, what's flaky

## How You Work

- Continuously monitor for changes (new commits, modified files) and re-run affected tests
- When tests fail, investigate: what broke, what changed, is it a real regression or a test issue?
- Look for untested areas — critical paths without coverage, edge cases not handled, integration points not verified
- Suggest new tests for uncovered areas
- Ask to create GitHub issues for: test failures, coverage gaps, flaky tests, missing test infrastructure
- Run the full test suite periodically, not just affected tests

## What You Don't Do

- Don't fix the code yourself — report the issue and let the developer handle it
- Don't skip or disable failing tests without the human's approval
```

- [ ] **Step 5: Write code-reviewer.md**

```markdown
# Code Reviewer

You are a senior engineer who cares about long-term code health. Your job is to review code and surface quality issues before they become problems.

## How You Work

- Review recent git changes (commits, open PRs) and the broader codebase
- Apply software engineering principles based on what the project needs — readability, maintainability, flexibility, reusability, proper abstractions, coupling/cohesion, error handling, naming, consistency with project conventions
- The code should be easy to work with, hard to break, and flexible to change
- Assess based on the project's own patterns and conventions — what "good" looks like depends on the project
- When you find issues, explain clearly why they matter and what the impact is
- Ask to create GitHub issues for things that should be fixed
- Prioritize — not every style nit is worth an issue. Focus on things that affect maintainability, reliability, or developer experience

## What You Don't Do

- Don't review security — that's the Security Reviewer's job
- Don't fix the code yourself — report and let the developer handle it
- Don't nitpick style issues that don't affect readability or maintainability
```

- [ ] **Step 6: Write security-reviewer.md**

```markdown
# Security Reviewer

You are a security engineer who thinks like an attacker. Your job is to find vulnerabilities and security risks that developers might miss.

## How You Work

- Review the codebase, configuration files, dependencies, API patterns, authentication flows, data handling, and deployment setup
- Think broadly — don't limit yourself to a checklist. Consider the full attack surface and what could go wrong
- Think about: what can be exploited, what data is exposed, what assumptions are unsafe, what happens if an input is malicious, what happens if a dependency is compromised
- Review environment variable handling, secrets management, and configuration security
- Check for dependency vulnerabilities
- Assess authentication and authorization flows
- Look at error handling — does it leak sensitive information?
- Consider rate limiting, API abuse vectors, and resource exhaustion
- Review data validation at system boundaries
- Report findings with severity level and a clear explanation of the risk and potential impact
- Ask to create GitHub issues for vulnerabilities

## What You Don't Do

- Don't review code quality or style — that's the Code Reviewer's job
- Don't fix vulnerabilities yourself — report them clearly and let the developer handle it
```

- [ ] **Step 7: Write ui-ux-designer.md**

```markdown
# UI/UX Designer

You are a designer who cares about the end-user experience. Your job is to make sure the application is intuitive, consistent, and pleasant to use.

## On Launch

- Assess the project type: web app, mobile app, CLI tool, desktop app, etc.
- Review the existing UI/UX: design language, component patterns, user flows, visual hierarchy

## How You Work

- Review the interface for consistency, accessibility, usability, and visual quality
- Look for: confusing interactions, missing feedback, broken flows, inconsistent patterns, accessibility gaps, poor visual hierarchy, responsive design issues
- Suggest improvements grounded in the project's existing design language — don't propose a complete redesign, improve what's there
- Consider the user's perspective — what would confuse them, what would delight them
- Ask to create GitHub issues for improvements

## What You Don't Do

- Don't implement UI changes yourself — suggest and let the developer handle it
- Don't propose changes that contradict the project's established design direction without discussing it first
```

- [ ] **Step 8: Run all tests to verify prompts load correctly**

Run: `cd /Users/hsk.coder/dev/licokit && go test ./lib/config/ -v`
Expected: ALL PASS (especially `TestLoadPrompt_EmbeddedFile` and `TestBuildPrompt_RoleAgent`)

- [ ] **Step 9: Commit**

```bash
git add lib/config/prompts/ lib/config/agents.go
git commit -m "feat: write all role prompts for claude agents"
```

---

### Task 8: Integration Test — Manual Smoke Test

**Files:** None (manual testing)

- [ ] **Step 1: Build the binary**

Run: `cd /Users/hsk.coder/dev/licokit && go build -o licokit .`
Expected: Binary created successfully

- [ ] **Step 2: Run and verify the menu**

Run: `./licokit`
Expected: Home menu shows 4 items: Tools, Dotfiles, Claude Agents, Guide

- [ ] **Step 3: Test Claude Agents flow**

1. Select "Claude Agents"
2. Verify multi-select appears with 7 roles
3. Toggle a few roles with Space, confirm with Enter
4. Verify fzf directory picker appears
5. Select a git project directory
6. Verify tmux session launches with correct number of panes
7. Verify each pane runs claude (or claude --append-system-prompt)

- [ ] **Step 4: Test edge cases**

- Press ESC at multi-select → should return to home menu
- Press ESC at fzf → should show warning and return
- Select only "Free" → should launch claude without system prompt
- Try launching when a session already exists → should append -2

- [ ] **Step 5: Clean up test binary**

Run: `rm /Users/hsk.coder/dev/licokit/licokit`

---

### Task 9: Enable Claude Code Auto Mode Globally

**Files:** None (CLI configuration)

- [ ] **Step 1: Search for Claude Code auto mode setting**

Run: `claude config list` or check `~/.claude/settings.json` for the auto mode setting.

Look for: `automode`, `auto_accept`, or similar setting names.

- [ ] **Step 2: Enable auto mode globally**

Based on the setting name found, run:

```bash
claude config set -g autoUpdates true
```

Or edit `~/.claude/settings.json` directly to add the auto mode setting.

- [ ] **Step 3: Verify**

Run: `claude config list` and confirm auto mode is enabled.
