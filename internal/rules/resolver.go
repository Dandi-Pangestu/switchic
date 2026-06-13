package rules

import (
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// LoadAll reads every rules/**/*.yaml in the bundled asset tree recursively.
// The registry key is the path relative to rules/ without the .yaml extension
// (e.g. "golang", "backend/api"). Definition.Dir holds the subdirectory portion.
func LoadAll() (map[string]Definition, error) {
	return LoadAllFrom(assets.FS())
}

// LoadAllFrom reads every rules/**/*.yaml in the given filesystem recursively.
// A missing "rules" directory is treated as empty rather than an error.
func LoadAllFrom(fsys fs.FS) (map[string]Definition, error) {
	out := map[string]Definition{}
	if _, err := fs.Stat(fsys, "rules"); err != nil {
		return out, nil
	}
	err := fs.WalkDir(fsys, "rules", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return util.Wrap(err, "read %s", p)
		}
		var def Definition
		if err := yaml.Unmarshal(data, &def); err != nil {
			return util.Wrap(err, "parse %s", p)
		}
		// p is like "rules/backend/api.yaml" — strip prefix and extension to get the key.
		rel := strings.TrimPrefix(p, "rules/")
		rel = strings.TrimSuffix(rel, ".yaml")
		if def.Name == "" {
			def.Name = path.Base(rel)
		}
		dir := path.Dir(rel)
		if dir == "." {
			dir = ""
		}
		def.Dir = dir
		out[rel] = def
		return nil
	})
	return out, err
}

// Names returns sorted rule names in the bundle.
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

// ResolveActive returns the active subset.
// Each entry in active is either an exact key ("backend/api") or a directory
// prefix ("backend"), which expands to all rules whose key starts with that
// prefix followed by "/". Expanded rules are returned in sorted key order.
func ResolveActive(all map[string]Definition, active []string) []Definition {
	out := make([]Definition, 0, len(active))
	for _, name := range active {
		if d, ok := all[name]; ok {
			out = append(out, d)
			continue
		}
		out = append(out, dirExpand(all, name)...)
	}
	return out
}

// Validate ensures every active name exists in the registry, either as an
// exact key or as a non-empty directory prefix.
func Validate(all map[string]Definition, active []string) []error {
	var errs []error
	for _, name := range active {
		if _, ok := all[name]; ok {
			continue
		}
		if len(dirExpand(all, name)) > 0 {
			continue
		}
		errs = append(errs, util.Wrap(util.ErrNotFound, "rule %q", name))
	}
	return errs
}

// dirExpand returns all definitions whose key starts with prefix+"/" sorted by key.
func dirExpand(all map[string]Definition, prefix string) []Definition {
	pfx := prefix + "/"
	keys := make([]string, 0)
	for k := range all {
		if strings.HasPrefix(k, pfx) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	out := make([]Definition, 0, len(keys))
	for _, k := range keys {
		out = append(out, all[k])
	}
	return out
}
