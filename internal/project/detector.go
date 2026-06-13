package project

import (
	"path/filepath"

	"github.com/Dandi-Pangestu/switchic/internal/util"
)

// Detect guesses the primary language of a repo by looking at well-known
// manifest files. Returns "" if it cannot tell.
func Detect(root string) string {
	checks := []struct {
		file string
		lang string
	}{
		{"go.mod", "go"},
		{"package.json", "typescript"},
		{"tsconfig.json", "typescript"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
		{"requirements.txt", "python"},
		{"pom.xml", "java"},
		{"build.gradle", "java"},
		{"Gemfile", "ruby"},
	}
	for _, c := range checks {
		if util.FileExists(filepath.Join(root, c.file)) {
			return c.lang
		}
	}
	return ""
}
