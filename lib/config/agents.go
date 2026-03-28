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
