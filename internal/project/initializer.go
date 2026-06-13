package project

import (
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/util"
	"github.com/Dandi-Pangestu/switchic/internal/workspace"
)

// InitWorkspace writes a default workspace manifest at root if one does not
// exist. The manifest is seeded with the bundled defaults for platform,
// workflow, and active agents/skills/rules so that `switch claude` produces
// a non-empty CLAUDE.md out of the box.
func InitWorkspace(root, name string) (workspace.Manifest, error) {
	path := util.WorkspacePath(root)
	if util.FileExists(path) {
		return workspace.Manifest{}, util.ErrAlreadyExists
	}
	defaults, err := config.Defaults()
	if err != nil {
		return workspace.Manifest{}, err
	}
	m := workspace.Manifest{
		Name:     name,
		Platform: defaults.Platform,
		Workflows: defaults.Workflows,
		Agents:   defaults.Agents,
		Skills:   defaults.Skills,
		Rules:    defaults.Rules,
	}
	if err := workspace.Save(path, m); err != nil {
		return workspace.Manifest{}, err
	}
	return m, nil
}
