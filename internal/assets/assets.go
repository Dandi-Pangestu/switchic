// Package assets exposes the bundled default YAML files (platforms, workflows,
// agents, skills, rules) as an embedded filesystem. Users do not need to
// install these separately — they ship inside the switchic binary.
package assets

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:bundled
var raw embed.FS

// FS returns the bundled filesystem rooted at "bundled/". Callers see
// "configs/...", "agents/...", etc. — the "bundled/" prefix is stripped.
func FS() fs.FS {
	sub, err := fs.Sub(raw, "bundled")
	if err != nil {
		// Should never happen — the directory is statically embedded.
		panic(err)
	}
	return sub
}

// LocalFS returns an fs.FS rooted at root/.switchic, or nil if that directory
// does not exist. Registries use it to load user-defined assets that take
// precedence over the bundled defaults — same subdirectory layout applies
// (agents/, skills/, rules/, workflows/).
func LocalFS(root string) fs.FS {
	dir := filepath.Join(root, ".switchic")
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	return os.DirFS(dir)
}
