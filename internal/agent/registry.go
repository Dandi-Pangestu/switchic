package agent

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// LoadAll reads every agents/*.yaml from the bundled asset tree.
func LoadAll() (map[string]Definition, error) {
	return LoadAllFrom(assets.FS())
}

// LoadAllFrom reads every agents/*.yaml from the given filesystem.
// A missing "agents" directory is treated as empty rather than an error,
// so callers can safely pass a user-local FS that may not define any agents.
func LoadAllFrom(fsys fs.FS) (map[string]Definition, error) {
	out := map[string]Definition{}
	entries, err := fs.ReadDir(fsys, "agents")
	if err != nil {
		return out, nil
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, filepath.Join("agents", e.Name()))
		if err != nil {
			return nil, util.Wrap(err, "read %s", e.Name())
		}
		var d Definition
		if err := yaml.Unmarshal(data, &d); err != nil {
			return nil, util.Wrap(err, "parse %s", e.Name())
		}
		if d.Name == "" {
			d.Name = strings.TrimSuffix(e.Name(), ".yaml")
		}
		out[d.Name] = d
	}
	return out, nil
}

// Names returns the sorted list of agent names available in the bundle.
func Names() ([]string, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(all))
	for k := range all {
		names = append(names, k)
	}
	sort.Strings(names)
	return names, nil
}
