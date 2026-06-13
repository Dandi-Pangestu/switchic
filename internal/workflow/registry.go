package workflow

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// LoadAll reads every workflows/*.yaml from the bundled asset tree.
func LoadAll() (map[string]Workflow, error) {
	return LoadAllFrom(assets.FS())
}

// LoadAllFrom reads every workflows/*.yaml from the given filesystem.
// A missing "workflows" directory is treated as empty rather than an error.
func LoadAllFrom(fsys fs.FS) (map[string]Workflow, error) {
	out := map[string]Workflow{}
	entries, err := fs.ReadDir(fsys, "workflows")
	if err != nil {
		return out, nil
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, filepath.Join("workflows", e.Name()))
		if err != nil {
			return nil, util.Wrap(err, "read %s", e.Name())
		}
		var w Workflow
		if err := yaml.Unmarshal(data, &w); err != nil {
			return nil, util.Wrap(err, "parse %s", e.Name())
		}
		if w.Name == "" {
			w.Name = strings.TrimSuffix(e.Name(), ".yaml")
		}
		out[w.Name] = w
	}
	return out, nil
}

// Get returns the named workflow or ErrNotFound.
func Get(name string) (Workflow, error) {
	all, err := LoadAll()
	if err != nil {
		return Workflow{}, err
	}
	w, ok := all[name]
	if !ok {
		return Workflow{}, util.Wrap(util.ErrNotFound, "workflow %q", name)
	}
	return w, nil
}

// Names returns sorted workflow names in the bundle.
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
