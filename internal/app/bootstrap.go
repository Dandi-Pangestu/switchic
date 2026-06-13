package app

import (
	"errors"

	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

// Context captures which mode the user is operating in and which configs are
// loaded. Commands call LoadContext at startup to centralize this discovery.
type Context struct {
	WorkingDir    string
	IsWorkspace   bool
	WorkspaceRoot string
	Project       config.Project
	Workspace     workspace.Manifest
}

// LoadContext discovers whether the user is in a workspace, a single-repo
// project, both, or neither. Missing files are not errors — the caller
// decides whether that is OK for the command at hand.
func LoadContext(cwd string) (Context, error) {
	c := Context{WorkingDir: cwd}

	if root := util.FindWorkspaceRoot(cwd); root != "" {
		m, err := workspace.Load(util.WorkspacePath(root))
		if err != nil && !errors.Is(err, util.ErrNotFound) {
			return c, err
		}
		if err == nil {
			c.IsWorkspace = true
			c.WorkspaceRoot = root
			c.Workspace = m
		}
	}

	p, err := config.Load(util.ProjectConfigPath(cwd))
	if err != nil && !errors.Is(err, util.ErrNotFound) {
		return c, err
	}
	if err == nil {
		c.Project = p
	}

	return c, nil
}

// PrimaryRoot returns the directory commands should treat as the generation
// target: the workspace root when in workspace mode, otherwise cwd.
func (c Context) PrimaryRoot() string {
	if c.IsWorkspace {
		return c.WorkspaceRoot
	}
	return c.WorkingDir
}

// HasProject is true when a .switchic/config.yaml was found.
func (c Context) HasProject() bool {
	return c.Project.Platform != "" || len(c.Project.Agents.Active) > 0
}
