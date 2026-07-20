package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLinkRepoAbsolutePath(t *testing.T) {
	root := t.TempDir()
	target := t.TempDir()
	r := Repo{Name: "svc", Path: target}

	if err := LinkRepo(root, r); err != nil {
		t.Fatalf("LinkRepo: %v", err)
	}

	linkPath := RepoLinkPath(root, "svc")
	got, err := os.Readlink(linkPath)
	if err != nil {
		t.Fatalf("Readlink: %v", err)
	}
	if got != target {
		t.Fatalf("link target = %q, want %q", got, target)
	}
}

func TestLinkRepoRelativePathIsNoop(t *testing.T) {
	root := t.TempDir()
	r := Repo{Name: "svc", Path: "../svc"}

	if err := LinkRepo(root, r); err != nil {
		t.Fatalf("LinkRepo: %v", err)
	}
	if _, err := os.Lstat(RepoLinkPath(root, "svc")); !os.IsNotExist(err) {
		t.Fatalf("expected no symlink created for relative path, lstat err = %v", err)
	}
}

func TestLinkRepoCollisionWithRealDir(t *testing.T) {
	root := t.TempDir()
	target := t.TempDir()
	r := Repo{Name: "svc", Path: target}

	linkPath := RepoLinkPath(root, "svc")
	if err := os.MkdirAll(linkPath, 0o755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := LinkRepo(root, r); err == nil {
		t.Fatal("expected error when a real directory occupies the link path")
	}

	fi, err := os.Lstat(linkPath)
	if err != nil {
		t.Fatalf("Lstat after failed link: %v", err)
	}
	if !fi.IsDir() {
		t.Fatal("real directory should not have been removed")
	}
}

func TestLinkRepoIdempotent(t *testing.T) {
	root := t.TempDir()
	target := t.TempDir()
	r := Repo{Name: "svc", Path: target}

	if err := LinkRepo(root, r); err != nil {
		t.Fatalf("first LinkRepo: %v", err)
	}
	if err := LinkRepo(root, r); err != nil {
		t.Fatalf("second LinkRepo: %v", err)
	}

	got, err := os.Readlink(RepoLinkPath(root, "svc"))
	if err != nil {
		t.Fatalf("Readlink: %v", err)
	}
	if got != target {
		t.Fatalf("link target = %q, want %q", got, target)
	}
}

func TestUnlinkRepo(t *testing.T) {
	root := t.TempDir()
	target := t.TempDir()
	r := Repo{Name: "svc", Path: target}

	if err := LinkRepo(root, r); err != nil {
		t.Fatalf("LinkRepo: %v", err)
	}
	if err := UnlinkRepo(root, "svc"); err != nil {
		t.Fatalf("UnlinkRepo: %v", err)
	}
	if _, err := os.Lstat(RepoLinkPath(root, "svc")); !os.IsNotExist(err) {
		t.Fatalf("expected symlink removed, lstat err = %v", err)
	}
}

func TestUnlinkRepoMissingIsNoop(t *testing.T) {
	root := t.TempDir()
	if err := UnlinkRepo(root, "does-not-exist"); err != nil {
		t.Fatalf("UnlinkRepo on missing link: %v", err)
	}
}

func TestUnlinkRepoLeavesNonSymlinkAlone(t *testing.T) {
	root := t.TempDir()
	linkPath := RepoLinkPath(root, "svc")
	if err := os.MkdirAll(linkPath, 0o755); err != nil {
		t.Fatalf("setup MkdirAll: %v", err)
	}

	if err := UnlinkRepo(root, "svc"); err != nil {
		t.Fatalf("UnlinkRepo: %v", err)
	}
	if _, err := os.Lstat(linkPath); err != nil {
		t.Fatalf("real directory should still exist: %v", err)
	}
}

func TestEnsureGitignoreEntryNoopWithoutFile(t *testing.T) {
	root := t.TempDir()
	if err := EnsureGitignoreEntry(root); err != nil {
		t.Fatalf("EnsureGitignoreEntry: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".gitignore")); !os.IsNotExist(err) {
		t.Fatal("expected no .gitignore to be created")
	}
}

func TestEnsureGitignoreEntryAppendsOnce(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules/\n"), 0o644); err != nil {
		t.Fatalf("setup WriteFile: %v", err)
	}

	if err := EnsureGitignoreEntry(root); err != nil {
		t.Fatalf("EnsureGitignoreEntry: %v", err)
	}
	if err := EnsureGitignoreEntry(root); err != nil {
		t.Fatalf("EnsureGitignoreEntry (second call): %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	if want := "node_modules/\n/repos/\n"; content != want {
		t.Fatalf("gitignore content = %q, want %q", content, want)
	}
}
