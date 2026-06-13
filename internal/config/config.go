// Package config defines the project-level config model (.switchic/config.yaml)
// and the merging logic between bundled defaults, project overrides, and
// workspace overrides.
package config

import "slices"

// ActiveList is the standard "active: [...]" container used for agents,
// skills, and rules. Keeping it a struct (rather than a bare slice) leaves
// room for future fields like "inactive" or "groups" without breaking YAML.
type ActiveList struct {
	Active []string `yaml:"active"`
}

// DocRef is a pointer to a documentation file with an optional "Read when" trigger.
type DocRef struct {
	Path string `yaml:"path"`
	When string `yaml:"when,omitempty"`
}

// Command pairs a shell command with a human-readable description.
type Command struct {
	Run         string `yaml:"run"`
	Description string `yaml:"description,omitempty"`
}

// Project is the .switchic/config.yaml shape.
type Project struct {
	Platform     string            `yaml:"platform"`
	Workflows    ActiveList        `yaml:"workflows"`
	Language     string            `yaml:"language,omitempty"`
	Name         string            `yaml:"name,omitempty"`
	Description  string            `yaml:"description,omitempty"`
	Stack        []string          `yaml:"stack,omitempty"`
	Commands     map[string]Command `yaml:"commands,omitempty"`
	Structure    map[string]string `yaml:"structure,omitempty"`
	Conventions  []string          `yaml:"conventions,omitempty"`
	Dos          []string          `yaml:"dos,omitempty"`
	Donts        []string          `yaml:"donts,omitempty"`
	Docs         []DocRef          `yaml:"docs,omitempty"`
	Agents       ActiveList        `yaml:"agents"`
	Skills       ActiveList        `yaml:"skills"`
	Rules        ActiveList        `yaml:"rules"`
}

// Clone returns a deep copy so callers can mutate without aliasing.
func (p Project) Clone() Project {
	cp := p
	cp.Workflows = ActiveList{Active: slices.Clone(p.Workflows.Active)}
	cp.Agents = ActiveList{Active: slices.Clone(p.Agents.Active)}
	cp.Skills = ActiveList{Active: slices.Clone(p.Skills.Active)}
	cp.Rules = ActiveList{Active: slices.Clone(p.Rules.Active)}
	cp.Stack = slices.Clone(p.Stack)
	cp.Docs = slices.Clone(p.Docs)
	cp.Conventions = slices.Clone(p.Conventions)
	cp.Dos = slices.Clone(p.Dos)
	cp.Donts = slices.Clone(p.Donts)
	if len(p.Commands) > 0 {
		cp.Commands = make(map[string]Command, len(p.Commands))
		for k, v := range p.Commands {
			cp.Commands[k] = v
		}
	}
	if len(p.Structure) > 0 {
		cp.Structure = make(map[string]string, len(p.Structure))
		for k, v := range p.Structure {
			cp.Structure[k] = v
		}
	}
	return cp
}
