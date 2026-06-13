package util

import (
	"os"
	"path/filepath"
)

// Project config locations relative to the working directory.
const (
	ProjectConfigDir  = ".switchic"
	ProjectConfigFile = ".switchic/config.yaml"
	RuntimeDir        = ".switchic/runtime"
	WorkspaceFile     = "switchic.workspace.yaml"
)

// ProjectConfigPath joins root with the standard project config location.
func ProjectConfigPath(root string) string {
	return filepath.Join(root, ProjectConfigFile)
}

// WorkspacePath joins root with the standard workspace manifest location.
func WorkspacePath(root string) string {
	return filepath.Join(root, WorkspaceFile)
}

// FindWorkspaceRoot walks upward from start looking for a workspace manifest.
// Returns the directory containing it, or "" if not found.
func FindWorkspaceRoot(start string) string {
	dir := start
	for {
		if FileExists(filepath.Join(dir, WorkspaceFile)) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// Cwd returns the current working directory or "." on failure.
func Cwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "."
	}
	return d
}
