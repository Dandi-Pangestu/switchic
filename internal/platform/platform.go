// Package platform defines the abstract platform adapter and the bundled
// adapters that ship with switchic. The MVP only implements Claude.
package platform

import (
	"github.com/Dandi-Pangestu/switchic/internal/agent"
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/rules"
	"github.com/Dandi-Pangestu/switchic/internal/skill"
	"github.com/Dandi-Pangestu/switchic/internal/workflow"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

// Context bundles everything an adapter needs to render its output. The
// adapter decides which fields apply to its platform — fields it doesn't use
// are ignored.
type Context struct {
	// Root is the directory the adapter should write to. For single-repo
	// mode this is the project root; for workspace mode it's the workspace
	// manifest's directory.
	Root string

	// Project is the merged project config; zero in pure-workspace mode.
	Project config.Project

	// Workspace is the manifest; zero in single-repo mode.
	Workspace workspace.Manifest

	// IsWorkspace is true when generating for a multi-repo workspace.
	IsWorkspace bool

	// Resolved component lists — the adapter should write only these.
	Workflows []workflow.Workflow
	Agents    []agent.Definition
	Skills    []skill.Definition
	Rules     []rules.Definition
}

// Adapter is implemented by every platform. Generate is the only required
// behavior: produce the platform-specific files in ctx.Root.
type Adapter interface {
	Name() string
	Generate(ctx Context) ([]string, error)
}
