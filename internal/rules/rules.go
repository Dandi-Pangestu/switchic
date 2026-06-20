// Package rules holds rule definitions loaded from bundled YAMLs.
package rules

// Definition mirrors the rules/<name>.yaml shape.
type Definition struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Content     string         `yaml:"content"`
	Claude      *ClaudeConfig  `yaml:"claude,omitempty"`
	Copilot     *CopilotConfig `yaml:"copilot,omitempty"`
	Kiro        *KiroConfig    `yaml:"kiro,omitempty"`

	// Dir is the subdirectory path relative to rules/ (e.g. "backend", "frontend/components").
	// Empty for top-level rules. Set during load, not from YAML.
	Dir string `yaml:"-"`
}

// KiroConfig holds Kiro-specific rule metadata.
// These are emitted as frontmatter in .kiro/steering/<name>.md.
type KiroConfig struct {
	// Inclusion controls when Kiro loads this steering file.
	// Valid values: "always" (default), "fileMatch", "manual", "auto".
	Inclusion string `yaml:"inclusion,omitempty"`

	// FileMatchPattern is used when Inclusion is "fileMatch".
	// Accepts a single glob string or a list of glob strings.
	FileMatchPattern any `yaml:"file_match_pattern,omitempty"`

	// Name is used when Inclusion is "auto".
	// Serves as the identifier used in slash commands.
	Name string `yaml:"name,omitempty"`

	// Description is used when Inclusion is "auto".
	// Kiro matches this against user requests to decide when to include the file.
	Description string `yaml:"description,omitempty"`
}

// CopilotConfig holds GitHub Copilot-specific rule metadata.
// These are emitted as frontmatter in .github/instructions/<name>.instructions.md.
type CopilotConfig struct {
	// ApplyTo is a comma-separated list of glob patterns scoping this instruction file.
	// Falls back to claude.paths when absent.
	ApplyTo string `yaml:"apply-to,omitempty"`
	// ExcludeAgent prevents a specific agent ("code-review" or "cloud-agent") from
	// reading this instruction file.
	ExcludeAgent string `yaml:"exclude-agent,omitempty"`
}

// ClaudeConfig holds Claude-specific rule metadata.
// Paths, when set, scopes the rule to files matching those glob patterns
// (rendered as YAML frontmatter in .claude/rules/<name>.md).
type ClaudeConfig struct {
	Paths []string `yaml:"paths,omitempty"`
}
