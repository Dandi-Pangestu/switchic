package workspace

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Dandi-Pangestu/switchic/internal/util"
)

// LinkDir is the directory, relative to the workspace root, holding symlinks
// to out-of-tree repos (those registered with an absolute path).
const LinkDir = "repos"

// RepoLinkPath returns where name's symlink would live under workspaceRoot.
func RepoLinkPath(workspaceRoot, name string) string {
	return filepath.Join(workspaceRoot, LinkDir, name)
}

// LinkRepo creates or refreshes the symlink for r under workspaceRoot/repos.
// It is a no-op for repos registered with a relative path, since those are
// already reachable in-tree. It never removes a non-symlink that collides
// with the target path; that is reported as an error instead.
func LinkRepo(workspaceRoot string, r Repo) error {
	if !filepath.IsAbs(r.Path) {
		return nil
	}
	linkPath := RepoLinkPath(workspaceRoot, r.Name)

	fi, err := os.Lstat(linkPath)
	switch {
	case os.IsNotExist(err):
		if err := util.EnsureDir(filepath.Dir(linkPath)); err != nil {
			return util.Wrap(err, "create %s", filepath.Dir(linkPath))
		}
		return util.Wrap(os.Symlink(r.Path, linkPath), "link %s -> %s", linkPath, r.Path)
	case err != nil:
		return util.Wrap(err, "stat %s", linkPath)
	case fi.Mode()&os.ModeSymlink == 0:
		return util.Wrap(util.ErrAlreadyExists, "%s exists and is not a switchic-managed symlink", linkPath)
	}

	target, err := os.Readlink(linkPath)
	if err != nil {
		return util.Wrap(err, "read link %s", linkPath)
	}
	if target == r.Path {
		return nil
	}
	if err := os.Remove(linkPath); err != nil {
		return util.Wrap(err, "remove stale link %s", linkPath)
	}
	return util.Wrap(os.Symlink(r.Path, linkPath), "link %s -> %s", linkPath, r.Path)
}

// UnlinkRepo removes name's symlink under workspaceRoot/repos, if present. It
// only removes an actual symlink; anything else at that path is left alone.
func UnlinkRepo(workspaceRoot, name string) error {
	linkPath := RepoLinkPath(workspaceRoot, name)
	fi, err := os.Lstat(linkPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return util.Wrap(err, "stat %s", linkPath)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return nil
	}
	return util.Wrap(os.Remove(linkPath), "remove %s", linkPath)
}

// LinkAll regenerates symlinks for every absolute-path repo in m. Failures
// are collected per repo name rather than aborting the whole run, so one
// colliding or unlinkable repo doesn't block the rest.
func LinkAll(workspaceRoot string, m Manifest) map[string]error {
	errs := map[string]error{}
	for _, r := range m.Repos {
		if !filepath.IsAbs(r.Path) {
			continue
		}
		if err := LinkRepo(workspaceRoot, r); err != nil {
			errs[r.Name] = err
		}
	}
	return errs
}

// EnsureGitignoreEntry adds a "/repos/" line to workspaceRoot/.gitignore if
// that file already exists and doesn't already ignore it. It never creates a
// .gitignore that isn't already there, since the workspace root need not be a
// git repo itself.
func EnsureGitignoreEntry(workspaceRoot string) error {
	path := filepath.Join(workspaceRoot, ".gitignore")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return util.Wrap(err, "read %s", path)
	}

	content := string(data)
	for line := range strings.SplitSeq(content, "\n") {
		if strings.TrimSpace(line) == "/"+LinkDir+"/" {
			return nil
		}
	}

	if !strings.HasSuffix(content, "\n") && content != "" {
		content += "\n"
	}
	content += "/" + LinkDir + "/\n"
	return util.WriteFile(path, []byte(content))
}
