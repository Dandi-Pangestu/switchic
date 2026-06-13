// Package rules holds rule definitions loaded from bundled YAMLs.
package rules

// Definition mirrors the rules/<name>.yaml shape.
type Definition struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Content     string        `yaml:"content"`
	Claude      *ClaudeConfig `yaml:"claude,omitempty"`

	// Dir is the subdirectory path relative to rules/ (e.g. "backend", "frontend/components").
	// Empty for top-level rules. Set during load, not from YAML.
	Dir string `yaml:"-"`
}

// ClaudeConfig holds Claude-specific rule metadata.
// Paths, when set, scopes the rule to files matching those glob patterns
// (rendered as YAML frontmatter in .claude/rules/<name>.md).
type ClaudeConfig struct {
	Paths []string `yaml:"paths,omitempty"`
}
