// Package workflow defines the workflow data model and a registry over the
// bundled workflow YAMLs.
package workflow

// Workflow is the workflows/<name>.yaml shape.
// A workflow is a pure preset: it declares which agents and skills are active
// for a given task type. Orchestration order is defined in the orchestrator
// agent's own instructions.
type Workflow struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Agents      []string `yaml:"agents,omitempty"`
	Skills      []string `yaml:"skills,omitempty"`
}
