package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gitscope/internal/state"
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

func TestParseRemotes(t *testing.T) {
	out := "origin\thttps://github.com/x/repo.git (fetch)\n" +
		"origin\thttps://github.com/x/repo.git (push)\n" +
		"upstream\thttps://github.com/y/repo.git (fetch)\n" +
		"upstream\thttps://github.com/y/repo.git (push)\n"
	remotes := parseRemotes(out)
	if len(remotes) != 2 {
		t.Fatalf("expected 2 remotes, got %d: %+v", len(remotes), remotes)
	}
	if remotes[0].Name != "origin" || remotes[0].FetchURL != "https://github.com/x/repo.git" {
		t.Errorf("unexpected origin: %+v", remotes[0])
	}
	if remotes[1].Name != "upstream" || remotes[1].PushURL != "https://github.com/y/repo.git" {
		t.Errorf("unexpected upstream: %+v", remotes[1])
	}
}

func TestParseTagList(t *testing.T) {
	tags := parseTagList("v1.0.0\nv1.1.0\n\n  v2.0.0  \n")
	want := []string{"v1.0.0", "v1.1.0", "v2.0.0"}
	if len(tags) != len(want) {
		t.Fatalf("expected %d tags, got %d: %v", len(want), len(tags), tags)
	}
	for i := range want {
		if tags[i] != want[i] {
			t.Errorf("tag %d: expected %q, got %q", i, want[i], tags[i])
		}
	}
}

func TestConfigGetRejectsBadKeys(t *testing.T) {
	// Initialize a real repo so validation passes and the key guard is exercised.
	dir := t.TempDir()
	if cmd := exec.Command("git", "-C", dir, "init"); cmd.Run() != nil {
		t.Skip("git not available")
	}
	prev := state.RepoPath
	state.RepoPath = dir
	defer func() { state.RepoPath = prev }()

	if _, err := ConfigGet(""); err == nil {
		t.Error("expected error for empty config key")
	}
	if _, err := ConfigGet("user name"); err == nil {
		t.Error("expected error for config key containing whitespace")
	}
}
