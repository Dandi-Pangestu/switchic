// Package workspace handles the multi-repo workspace manifest
// (switchic.workspace.yaml).
package workspace

import "github.com/Dandi-Pangestu/switchic/internal/config"

// Repo describes one repository inside a workspace.
type Repo struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Role  string `yaml:"role,omitempty"`
	Notes string `yaml:"notes,omitempty"`
}

// Manifest is the switchic.workspace.yaml shape. Active component lists may
// be inherited from a project config; when present here they take priority.
type Manifest struct {
	Name     string             `yaml:"name"`
	Platform string             `yaml:"platform"`
	Workflows config.ActiveList `yaml:"workflows,omitempty"`
	Repos    []Repo             `yaml:"repos"`
	Agents   config.ActiveList  `yaml:"agents,omitempty"`
	Skills   config.ActiveList  `yaml:"skills,omitempty"`
	Rules    config.ActiveList  `yaml:"rules,omitempty"`
	Notes    string             `yaml:"notes,omitempty"`
}
