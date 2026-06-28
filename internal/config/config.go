// Package config defines the project-level config model (.switchic/config.yaml)
// and the merging logic between bundled defaults, project overrides, and
// workspace overrides.
package config

import (
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// StringSlice is a []string that tolerates list items YAML parses as inline
// mappings (e.g. "frozen_string_literal: true" or "bunny: :mock") by
// reconstructing them as "key: value" strings. Items that are already plain
// scalars pass through unchanged.
type StringSlice []string

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *StringSlice) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("expected a sequence, got %s", value.Tag)
	}
	result := make([]string, 0, len(value.Content))
	for _, item := range value.Content {
		switch item.Kind {
		case yaml.ScalarNode:
			result = append(result, item.Value)
		case yaml.MappingNode:
			str, err := mappingNodeToString(item)
			if err != nil {
				return err
			}
			result = append(result, str)
		default:
			return fmt.Errorf("unsupported sequence item tag: %s", item.Tag)
		}
	}
	*s = result
	return nil
}

// mappingNodeToString reconstructs a plain "key: value" string from a YAML
// mapping node that was originally an unquoted string containing ": ".
func mappingNodeToString(node *yaml.Node) (string, error) {
	if len(node.Content)%2 != 0 {
		return "", fmt.Errorf("malformed mapping node")
	}
	var sb strings.Builder
	for i := 0; i+1 < len(node.Content); i += 2 {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(node.Content[i].Value)
		sb.WriteString(": ")
		sb.WriteString(node.Content[i+1].Value)
	}
	return sb.String(), nil
}

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
	Conventions  StringSlice       `yaml:"conventions,omitempty"`
	Dos          StringSlice       `yaml:"dos,omitempty"`
	Donts        StringSlice       `yaml:"donts,omitempty"`
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
