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

func TestLoadAgents_DefaultConfig(t *testing.T) {
	cfg, err := LoadAgents()
	if err != nil {
		t.Fatalf("LoadAgents() error: %v", err)
	}

	if len(cfg.Agents) != 7 {
		t.Errorf("expected 7 agents, got %d", len(cfg.Agents))
	}

	for _, agent := range cfg.Agents {
		if agent.Name == "" {
			t.Error("agent name should not be empty")
		}
		if agent.Description == "" {
			t.Errorf("agent %q: description should not be empty", agent.Name)
		}
	}

	free := cfg.Agents[len(cfg.Agents)-1]
	if free.Name != "Free" {
		t.Errorf("last agent = %q, want Free", free.Name)
	}
	if free.PromptFile != "" {
		t.Errorf("Free agent prompt_file = %q, want empty", free.PromptFile)
	}
}

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
