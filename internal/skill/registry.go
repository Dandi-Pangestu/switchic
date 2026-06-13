package skill

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/assets"
	"github.com/Dandi-Pangestu/switchic/internal/util"

	"gopkg.in/yaml.v3"
)

// LoadAll reads skills from the bundled asset tree. Two formats are supported:
//
//   - Flat:   skills/<name>.yaml
//   - Folder: skills/<name>/skill.yaml  (+ any extra files in sub-folders)
func LoadAll() (map[string]Definition, error) {
	return LoadAllFrom(assets.FS())
}

// LoadAllFrom reads skills from the given filesystem. A missing "skills"
// directory is treated as empty rather than an error.
func LoadAllFrom(fsys fs.FS) (map[string]Definition, error) {
	out := map[string]Definition{}
	entries, err := fs.ReadDir(fsys, "skills")
	if err != nil {
		return out, nil
	}
	for _, e := range entries {
		if e.IsDir() {
			d, err := loadFolderFrom(fsys, e.Name())
			if err != nil {
				return nil, err
			}
			out[d.Name] = d
			continue
		}
		if !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, filepath.Join("skills", e.Name()))
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

// loadFolderFrom loads a folder-based skill from skills/<dirName>/skill.yaml
// and collects all other files under that directory as EmbeddedFiles.
func loadFolderFrom(fsys fs.FS, dirName string) (Definition, error) {
	skillYAML := filepath.Join("skills", dirName, "skill.yaml")
	data, err := fs.ReadFile(fsys, skillYAML)
	if err != nil {
		return Definition{}, util.Wrap(err, "read %s", skillYAML)
	}
	var d Definition
	if err := yaml.Unmarshal(data, &d); err != nil {
		return Definition{}, util.Wrap(err, "parse %s", skillYAML)
	}
	if d.Name == "" {
		d.Name = dirName
	}

	root := filepath.Join("skills", dirName)
	err = fs.WalkDir(fsys, root, func(path string, de fs.DirEntry, werr error) error {
		if werr != nil {
			return werr
		}
		if de.IsDir() || path == skillYAML {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return util.Wrap(err, "read embedded file %s", path)
		}
		d.Files = append(d.Files, EmbeddedFile{RelPath: rel, Content: content})
		return nil
	})
	if err != nil {
		return Definition{}, err
	}
	return d, nil
}

// Names returns the sorted list of skill names in the bundle.
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
