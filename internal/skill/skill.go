// Package skill holds skill definitions loaded from bundled YAMLs.
package skill

// EmbeddedFile is an additional file bundled inside a folder-based skill.
// RelPath is relative to the skill root (e.g. "examples/basic.md").
type EmbeddedFile struct {
	RelPath string
	Content []byte
}

// Definition mirrors the skills/<name>.yaml (flat) or skills/<name>/skill.yaml (folder) shape.
// Core fields are cross-platform. Platform-specific config lives under its own key.
type Definition struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Prompt      string         `yaml:"prompt,omitempty"`
	Claude      *ClaudeConfig  `yaml:"claude,omitempty"`
	Copilot     *CopilotConfig `yaml:"copilot,omitempty"`
	Kiro        *KiroConfig    `yaml:"kiro,omitempty"`

	// Files holds extra files from a folder-based skill; populated by LoadAll, not by YAML.
	Files []EmbeddedFile `yaml:"-"`
}

// KiroConfig holds Kiro CLI-specific skill frontmatter fields.
// These are emitted into .kiro/skills/<name>/SKILL.md when generating for the kiro platform.
type KiroConfig struct {
	AllowedTools  string            `yaml:"allowed-tools,omitempty"`
	Compatibility string            `yaml:"compatibility,omitempty"`
	License       string            `yaml:"license,omitempty"`
	Metadata      map[string]string `yaml:"metadata,omitempty"`
}

// CopilotConfig holds GitHub Copilot-specific skill frontmatter fields.
// These are emitted into .github/skills/<name>/SKILL.md per the agentskills.io spec.
type CopilotConfig struct {
	// AllowedTools is a space-separated list of pre-approved tools (experimental).
	AllowedTools string `yaml:"allowed-tools,omitempty"`
	// Compatibility describes environment requirements for this skill.
	Compatibility string `yaml:"compatibility,omitempty"`
	// License specifies the license applied to this skill.
	License  string            `yaml:"license,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

// ClaudeConfig holds Claude-specific skill frontmatter fields.
// These are emitted into .claude/skills/<name>/SKILL.md when generating for the Claude platform.
// Field names use the same kebab-case / snake_case keys Claude expects in SKILL.md frontmatter.
type ClaudeConfig struct {
	// Discovery
	// when_to_use: additional context appended to description in the skill listing.
	WhenToUse string `yaml:"when_to_use,omitempty"`

	// Invocation control
	// disable-model-invocation: true prevents Claude from loading this skill automatically.
	DisableModelInvocation bool `yaml:"disable-model-invocation,omitempty"`
	// user-invocable: false hides the skill from the / menu (Claude can still invoke it).
	// Uses *bool because the meaningful state is false; omit to keep the default (true).
	UserInvocable *bool `yaml:"user-invocable,omitempty"`

	// Argument handling
	// argument-hint: hint shown in autocomplete (e.g. "[issue-number]").
	ArgumentHint string `yaml:"argument-hint,omitempty"`
	// arguments: named positional arguments for $name substitution. String or list.
	Arguments interface{} `yaml:"arguments,omitempty"`

	// Tool access
	// allowed-tools: tools Claude may use without per-use approval while this skill is active.
	AllowedTools string `yaml:"allowed-tools,omitempty"`
	// disallowed-tools: tools removed from Claude's pool while this skill is active.
	DisallowedTools string `yaml:"disallowed-tools,omitempty"`

	// Model and execution
	// model: model override for this skill's turn. Reverts on the next user message.
	Model string `yaml:"model,omitempty"`
	// effort: thinking effort level for this skill's turn.
	Effort string `yaml:"effort,omitempty"`
	// context: set to "fork" to run the skill in an isolated subagent context.
	Context string `yaml:"context,omitempty"`
	// agent: subagent type to use when context is "fork" (e.g. "Explore", "general-purpose").
	Agent string `yaml:"agent,omitempty"`
	// shell: shell for !`command` blocks. Options: bash (default) | powershell.
	Shell string `yaml:"shell,omitempty"`

	// File-based activation
	// paths: glob patterns that limit automatic activation to matching files.
	Paths interface{} `yaml:"paths,omitempty"`

	// Lifecycle
	// hooks: hooks scoped to this skill. Same format as agent hooks.
	Hooks interface{} `yaml:"hooks,omitempty"`
}
