package util

import (
	"os"
	"path/filepath"
)

// EnsureDir creates dir (and any parents) with 0o755. No-op if it already exists.
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

// FileExists returns true iff path resolves to a regular file or any inode the
// caller treats as "present". It swallows non-existence errors and surfaces
// only stat failures via the second return.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WriteFile atomically replaces the file at path with the given bytes. The
// parent directory is created if needed. Mode defaults to 0o644.
func WriteFile(path string, data []byte) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return Wrap(err, "create parent directory")
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return Wrap(err, "write %s", tmp)
	}
	if err := os.Rename(tmp, path); err != nil {
		return Wrap(err, "rename %s -> %s", tmp, path)
	}
	return nil
}
