package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRepoPath(t *testing.T) {
	dir := t.TempDir()
	if err := validateRepoPath(dir); err != nil {
		t.Errorf("expected nil for valid dir, got %v", err)
	}

	if err := validateRepoPath(filepath.Join(dir, "does-not-exist")); err == nil {
		t.Error("expected error for non-existent path")
	}

	file := filepath.Join(dir, "a-file")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := validateRepoPath(file); err == nil {
		t.Error("expected error when path is a file, not a directory")
	}
}

func TestValidateGitRepo(t *testing.T) {
	dir := t.TempDir()
	// A plain directory is not a Git work tree.
	if err := validateGitRepo(dir); err == nil {
		t.Error("expected error for non-Git directory")
	}
	if err := validateGitRepo(filepath.Join(dir, "missing")); err == nil {
		t.Error("expected error for non-existent path")
	}
}
