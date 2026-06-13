package config

import (
	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
	"io/fs"
)

// Defaults loads the bundled default project config from the embedded assets.
func Defaults() (Project, error) {
	data, err := fs.ReadFile(assets.FS(), "configs/default.yaml")
	if err != nil {
		return Project{}, util.Wrap(err, "read bundled default config")
	}
	var p Project
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Project{}, util.Wrap(err, "parse bundled default config")
	}
	return p, nil
}
