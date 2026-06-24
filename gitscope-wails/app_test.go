package main

import (
	"testing"

	"github.com/gitscope/internal/state"
)

func TestIsGitAvailable(t *testing.T) {
	app := NewApp()
	if !app.IsGitAvailable() {
		t.Skip("git not available in PATH")
	}
}

func TestGetRepoPathEmpty(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	if got := app.GetRepoPath(); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestGetRepoPathSet(t *testing.T) {
	state.RepoPath = "/tmp/test-repo"
	app := NewApp()
	if got := app.GetRepoPath(); got != "/tmp/test-repo" {
		t.Errorf("expected /tmp/test-repo, got %q", got)
	}
	state.RepoPath = ""
}

func TestIsRepoInitialized(t *testing.T) {
	app := NewApp()
	state.RepoPath = ""
	if app.IsRepoInitialized() {
		t.Error("expected false for empty path")
	}
}

func TestStatusNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Status("Standard")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestInitNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Init()
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestStageNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Stage("All (.)")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestGetBranchesEmpty(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.GetBranches()
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestDiffNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Diff("Unstaged")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestReadGitIgnoreNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.ReadGitIgnore()
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestWriteGitIgnoreNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	err := app.WriteGitIgnore("*.exe")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestLogNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Log("Oneline")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestResetNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Reset("--mixed", "HEAD~1")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestStashNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Stash("Save")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestMergeNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Merge("main")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestFetchNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.Fetch("Default")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}

func TestRunCommandsNoRepo(t *testing.T) {
	state.RepoPath = ""
	app := NewApp()
	_, err := app.RunCommands("git status")
	if err == nil {
		t.Error("expected error for empty repo path")
	}
}
