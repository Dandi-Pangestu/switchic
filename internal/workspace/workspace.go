package workspace

import (
	"path/filepath"
	"slices"

	"github.com/Dandi-Pangestu/switchic/internal/util"
)

// AddRepo registers a repo by absolute or relative path. Name defaults to the
// basename of the path if empty. Returns ErrAlreadyExists on duplicate name.
// notes is a short description shown in the generated context file.
// contextFile overrides the platform-default context file path for this repo
// (e.g. "docs/CLAUDE.md"); leave empty to use the platform default.
func (m *Manifest) AddRepo(path, role, notes, contextFile string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return util.Wrap(err, "resolve %s", path)
	}
	name := filepath.Base(abs)
	for _, r := range m.Repos {
		if r.Name == name {
			return util.ErrAlreadyExists
		}
	}
	m.Repos = append(m.Repos, Repo{
		Name:        name,
		Path:        path,
		Role:        role,
		Notes:       notes,
		ContextFile: contextFile,
	})
	return nil
}

// RemoveRepo deletes a repo by name. Returns ErrNotFound if absent.
func (m *Manifest) RemoveRepo(name string) error {
	for i, r := range m.Repos {
		if r.Name == name {
			m.Repos = slices.Delete(m.Repos, i, i+1)
			return nil
		}
	}
	return util.ErrNotFound
}
