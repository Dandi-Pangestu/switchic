package workspace

import (
	"path/filepath"

	"github.com/Dandi-Pangestu/switchic/internal/util"
)

// MissingRepos returns repos whose path does not exist on disk. The check is
// resolved relative to workspaceRoot.
func MissingRepos(workspaceRoot string, m Manifest) []Repo {
	var missing []Repo
	for _, r := range m.Repos {
		p := r.Path
		if !filepath.IsAbs(p) {
			p = filepath.Join(workspaceRoot, p)
		}
		if !util.FileExists(p) {
			missing = append(missing, r)
		}
	}
	return missing
}
