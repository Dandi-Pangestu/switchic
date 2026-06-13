package workspace

import (
	"errors"
	"os"

	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// Load reads a workspace manifest. Returns (zero, ErrNotFound) if missing.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Manifest{}, util.ErrNotFound
		}
		return Manifest{}, util.Wrap(err, "read %s", path)
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, util.Wrap(err, "parse %s", path)
	}
	return m, nil
}

// Save writes the manifest as YAML to path.
func Save(path string, m Manifest) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return util.Wrap(err, "marshal workspace manifest")
	}
	return util.WriteFile(path, data)
}
