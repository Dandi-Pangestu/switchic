// Package agent holds agent definitions loaded from bundled YAMLs.
package agent

// Definition mirrors the agents/<name>.yaml shape.
// Core fields are cross-platform. Platform-specific config lives under its
// own key (e.g. claude:) so adapters read only what they own.
type Definition struct {
	Name           string        `yaml:"name"`
	Description    string        `yaml:"description"`
	RequiredSkills []string      `yaml:"required_skills,omitempty"`
	Instructions   string        `yaml:"instructions,omitempty"`
	Claude         *ClaudeConfig `yaml:"claude,omitempty"`
	Copilot        *CopilotConfig `yaml:"copilot,omitempty"`
}

// CopilotConfig holds GitHub Copilot-specific frontmatter fields.
// These are emitted into .github/agents/<name>.agent.md when generating
// for the github-copilot platform.
type CopilotConfig struct {
	// Tool access — Copilot tool aliases: execute, read, edit, search, agent, web, todo.
	// Omit to grant all tools.
	Tools []string `yaml:"tools,omitempty"`

	// Model and execution
	Model string `yaml:"model,omitempty"`
	// Target restricts the agent to a specific environment ("vscode" or "github-copilot").
	// Defaults to "github-copilot" when generating for this platform.
	Target string `yaml:"target,omitempty"`

	// Invocation control
	DisableModelInvocation bool  `yaml:"disable-model-invocation,omitempty"`
	UserInvocable          *bool `yaml:"user-invocable,omitempty"`

	// Advanced integrations — arbitrary YAML structures passed through as-is
	McpServers interface{}       `yaml:"mcp-servers,omitempty"`
	Metadata   map[string]string `yaml:"metadata,omitempty"`
}

// ClaudeConfig holds Claude-specific frontmatter fields.
// These are emitted verbatim into .claude/agents/<name>.md when generating
// for the Claude platform. snake_case keys here map to Claude's camelCase keys.
type ClaudeConfig struct {
	// Tool access
	Tools           []string `yaml:"tools,omitempty"`
	DisallowedTools []string `yaml:"disallowed_tools,omitempty"`

	// Model and execution
	Model          string `yaml:"model,omitempty"`
	PermissionMode string `yaml:"permission_mode,omitempty"`
	MaxTurns       int    `yaml:"max_turns,omitempty"`
	Effort         string `yaml:"effort,omitempty"`
	Isolation      string `yaml:"isolation,omitempty"`
	Background     bool   `yaml:"background,omitempty"`

	// Context and memory
	Memory        string `yaml:"memory,omitempty"`
	InitialPrompt string `yaml:"initial_prompt,omitempty"`
	Color         string `yaml:"color,omitempty"`

	// Advanced integrations — arbitrary YAML structures passed through as-is
	McpServers interface{} `yaml:"mcp_servers,omitempty"`
	Hooks      interface{} `yaml:"hooks,omitempty"`
}
