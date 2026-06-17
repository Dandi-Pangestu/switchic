// Package rules holds rule definitions loaded from bundled YAMLs.
package rules

// Definition mirrors the rules/<name>.yaml shape.
type Definition struct {
	Name        string         `yaml:"name"`
	Description string         `yaml:"description"`
	Content     string         `yaml:"content"`
	Claude      *ClaudeConfig  `yaml:"claude,omitempty"`
	Copilot     *CopilotConfig `yaml:"copilot,omitempty"`

	// Dir is the subdirectory path relative to rules/ (e.g. "backend", "frontend/components").
	// Empty for top-level rules. Set during load, not from YAML.
	Dir string `yaml:"-"`
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
