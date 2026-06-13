package config

import (
	"errors"
	"os"

	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// Load reads the project config at path. Returns (zero, ErrNotFound) if missing.
func Load(path string) (Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Project{}, util.ErrNotFound
		}
		return Project{}, util.Wrap(err, "read %s", path)
	}
	var p Project
	if err := yaml.Unmarshal(data, &p); err != nil {
		return Project{}, util.Wrap(err, "parse %s", path)
	}
	return p, nil
}
