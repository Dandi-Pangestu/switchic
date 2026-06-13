// Package project handles single-repo bootstrap: detecting language and
// writing the initial .switchic/config.yaml.
package project

import (
	"github.com/Dandi-Pangestu/switchic/internal/config"
	"github.com/Dandi-Pangestu/switchic/internal/util"
)

// Init writes a default project config at root if one does not already exist.
// Returns ErrAlreadyExists if the file is present.
func Init(root string) (config.Project, error) {
	path := util.ProjectConfigPath(root)
	if util.FileExists(path) {
		return config.Project{}, util.ErrAlreadyExists
	}
	p, err := config.Defaults()
	if err != nil {
		return config.Project{}, err
	}
	if lang := Detect(root); lang != "" {
		p.Language = lang
	}
	if err := config.WriteInitial(path, p); err != nil {
		return config.Project{}, err
	}
	return p, nil
}
