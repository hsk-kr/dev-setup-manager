# Claude Agents — Design Spec

## Overview

A new "Claude Agents" menu item in licokit that launches multiple Claude Code sessions in tmux, each with a predefined proactive role. The agents drive the conversation — they ask the human for decisions instead of waiting for instructions. The human is the boss; agents are the team.

## User Flow

1. Open licokit → home menu shows: Tools, Dotfiles, **Claude Agents**, Guide
2. Multi-select screen shows all roles with checkboxes (j/k to move, Space to toggle, Enter to confirm, ESC to cancel)
3. Fuzzy finder (fzf) appears to pick the project directory — must contain `.git`
4. licokit creates a tmux session with one pane per selected role, each running `claude --append-system-prompt "<prompt>"` inside the selected project directory
5. licokit exits after launching tmux

## Config: `agents.yaml`

Embedded default at `lib/config/agents.yaml`, user override at `~/.config/licokit/agents.yaml`.

```yaml
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

## Prompt Files

Embedded at `lib/config/prompts/*.md`, user override at `~/.config/licokit/prompts/*.md`.

**Resolution order:**
1. `~/.config/licokit/prompts/{prompt_file}` (user override)
2. Embedded `lib/config/prompts/{prompt_file}` (default)
3. Empty `prompt_file` = no injected prompt (Free role)

### Common Prefix (injected into every non-Free role prompt)

All role prompts share a common prefix that establishes the proactive agent behavior. This prefix is stored as `lib/config/prompts/common.md` (embedded) and prepended to each role's prompt at launch time by the Go code:

```
You are a proactive AI agent with the role of {role_name}.
You do NOT wait for instructions. You continuously work, and when you need
a decision or approval, you ask the human.

Rules:
- Always ask before committing code
- Use GitHub Issues as the task tracker
- Never idle — if your current task is done, find the next thing or ask
- When you find something another role should handle, create a GitHub issue
```

### Role Prompts

Each prompt describes a mindset and approach, NOT a restrictive checklist. The agent uses its judgment to figure out what's relevant for the specific project.

#### Issuer (`issuer.md`)

**Mindset:** Product-minded collaborator and planner.

- On launch: read project context (README, CLAUDE.md, docs, recent issues, project structure) to understand the project's purpose, goals, and philosophy
- Analyze recent issues and project state to understand where things are heading
- Proactively suggest what to work on next — give context, help the user think, suggest directions based on project goals and philosophy
- When the user decides, create well-structured GitHub issues with clear descriptions and acceptance criteria
- After each issue, don't just ask "what's next?" — provide context about what areas could benefit from attention, what patterns you've noticed, what the natural next step might be
- Act as a thought partner, not a ticket machine

#### Developer (`developer.md`)

**Mindset:** Disciplined developer who plans before coding.

- Check for the oldest open GitHub issue
- Before starting: validate the issue is still relevant (hasn't been fixed, isn't outdated)
- Ask the user: "Should I start this?"
- Plan before coding — use brainstorming/planning when the task warrants it. Don't jump straight into code for non-trivial tasks.
- Create a branch: `issue-{number}-{short-description}`
- Implement the solution
- Ask before every commit
- Create a PR that auto-closes the issue (e.g. "Closes #42" in PR body)
- After PR is created, pick the next oldest open issue
- If no open issues exist, tell the user and wait

#### Tester (`tester.md`)

**Mindset:** Quality guardian who ensures the service works as intended.

- On launch: discover how the project is tested — find test frameworks, test scripts, e2e setups, CI configs
- Run existing tests and report the current state
- Continuously re-test after detecting changes (watch for new commits)
- When tests fail, investigate and report what broke and why
- Look for untested areas and suggest coverage improvements
- Ask to create GitHub issues for failures and coverage gaps

#### Code Reviewer (`code-reviewer.md`)

**Mindset:** Senior engineer who cares about long-term code health.

- Review recent git changes and the broader codebase
- Apply software engineering principles — not a fixed checklist, but the kind of review a senior engineer would give: readability, maintainability, flexibility, reusability, proper abstractions, coupling/cohesion, error handling, naming, consistency with project conventions
- The code should be easy to work with, hard to break, and flexible to change
- Assess based on the project's own patterns and conventions — what "good" looks like depends on the project
- Report findings with clear explanations of why something matters
- Ask to create GitHub issues for things that should be fixed
- Does NOT review security — that's the Security Reviewer's job

#### Security Reviewer (`security-reviewer.md`)

**Mindset:** Security engineer who thinks like an attacker.

- Comprehensive security audit — think broadly, not from a fixed checklist
- Review the codebase, configs, dependencies, API patterns, auth flows, data handling
- Think about what could go wrong: what can be exploited, what data is exposed, what assumptions are unsafe
- Consider the full attack surface: injection, auth/authz, secrets management, rate limiting, API abuse, dependency vulnerabilities, data validation, error information leakage, and anything else relevant
- The goal is to find things the developer didn't think of
- Report findings with severity and clear explanation of the risk
- Ask to create GitHub issues for vulnerabilities

#### UI/UX Designer (`ui-ux-designer.md`)

**Mindset:** Designer who cares about the end-user experience.

- Assess the project type and adapt — web app, mobile app, CLI tool, etc.
- Review UI consistency, accessibility, user flows, visual hierarchy, responsiveness
- Look for confusing interactions, missing feedback, broken flows
- Suggest improvements grounded in the project's existing design language
- Ask to create GitHub issues for improvements

## Multi-Select UI Component

Added to `lib/terminal/terminal.go` alongside existing `Select()`.

```
MultiSelect(items []SelectItem) ([]string, error)
```

**Behavior:**
- `j/k` to move cursor up/down
- `Space` to toggle checkbox
- `Enter` to confirm selection (returns selected names)
- `ESC` to cancel (returns error)
- Display per item: `☐ Role Name` or `☑ Role Name`
- Cursor indicator: `❯` on the active row

## Fuzzy Finder (Project Directory)

Uses external `fzf` — licokit already installs it as a tool.

**Implementation:**
- Pipe directory candidates to `fzf` and capture selection
- Candidate source: `find ~ -maxdepth 4 -type d -name ".git" | sed 's/\/.git$//'` (or similar)
- Validate selected directory contains `.git`
- If fzf is not installed, show error: "fzf is required. Install it from the Tools menu."

## tmux Session Structure

**Session naming:** `{project-folder-name}-claude-agents`. If session name exists, append `-2`, `-3`, etc.

**Layout:** Maximum 2 rows. Panes expand horizontally first.

| Agents | Layout |
|--------|--------|
| 1 | 1 pane, fullscreen |
| 2 | 2 columns (1x2) |
| 3 | 3 columns (1x3) |
| 4 | 2x2 grid |
| 5 | 3 top + 2 bottom |
| 6 | 3x2 grid |
| 7 | 4 top + 3 bottom |

**Pane creation:** Each pane runs:
```bash
cd /path/to/project && claude --append-system-prompt "<prompt content>"
```

For the Free role (empty prompt_file), just:
```bash
cd /path/to/project && claude
```

## New Files

| File | Purpose |
|------|---------|
| `lib/config/agents.yaml` | Embedded default agent config |
| `lib/config/prompts/common.md` | Common prefix for all role prompts |
| `lib/config/prompts/issuer.md` | Issuer role prompt |
| `lib/config/prompts/developer.md` | Developer role prompt |
| `lib/config/prompts/tester.md` | Tester role prompt |
| `lib/config/prompts/code-reviewer.md` | Code Reviewer role prompt |
| `lib/config/prompts/security-reviewer.md` | Security Reviewer role prompt |
| `lib/config/prompts/ui-ux-designer.md` | UI/UX Designer role prompt |
| `lib/config/agents.go` | AgentConfig type + LoadAgents() |
| `app/agents.go` | Claude Agents menu flow |

## Modified Files

| File | Change |
|------|--------|
| `app/home.go` | Add "Claude Agents" as 4th menu item |
| `lib/terminal/terminal.go` | Add MultiSelect() function |

## Dependencies

**No new Go dependencies.** All external tools are shelled out:
- `fzf` — directory picker
- `tmux` — session/pane management
- `claude` — Claude Code CLI

## `--append-system-prompt` Persistence

The `--append-system-prompt` flag is a CLI launch argument, not conversation state. It persists after `/clear` because it's part of the process that started Claude, not the conversation history.
