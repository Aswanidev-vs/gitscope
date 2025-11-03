package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gitscope/internal/state"
)

func Init() (string, error) {
	if state.RepoPath == "" {
		return "", errors.New("no repository path selected")
	}
	cmd := exec.Command("git", "-C", state.RepoPath, "init")
	out, err := cmd.Output()
	return string(out), err
}
func Status() (string, error) {
	cmd := exec.Command("git", "-C", state.RepoPath, "status")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Commit creates a new commit with the given message.
func Commit(msg string) (string, error) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "", errors.New("commit message cannot be empty")
	}

	// Validate repo path
	checkCmd := exec.Command("git", "-C", state.RepoPath, "rev-parse", "--is-inside-work-tree")
	if err := checkCmd.Run(); err != nil {
		return "", errors.New("invalid Git repository path")
	}

	// Check for staged changes
	statusCmd := exec.Command("git", "-C", state.RepoPath, "diff", "--cached", "--quiet")
	if err := statusCmd.Run(); err == nil {
		return "", errors.New("no staged changes to commit")
	}

	// Proceed with commit
	cmd := exec.Command("git", "-C", state.RepoPath, "commit", "-m", msg)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), fmt.Errorf("commit failed: %v\n%s", err, string(out))
	}

	if len(out) == 0 {
		return "Commit completed successfully.", nil
	}

	return string(out), nil
}

// Stage adds all modified/untracked files to the Git index (staging area)
func Stage() (string, error) {
	if state.RepoPath == "" {
		return "", errors.New("no repository path selected")
	}
	cmd := exec.Command("git", "-C", state.RepoPath, "add", ".")
	out, err := cmd.CombinedOutput()

	if err != nil {
		// Return full Git output for accurate debugging
		return string(out), fmt.Errorf("%v\n%s", err, string(out))
	}

	// Usually, git add . returns no output unless something goes wrong
	if len(out) == 0 {
		return "All changes staged successfully.", nil
	}

	return string(out), nil
}

func Push(repoPath, branch string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "push", "origin", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("push failed: %v\n%s", err, string(out))
	}
	return string(out), nil
}
func Log(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--oneline")
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), fmt.Errorf("log failed:%v\n%s", err, string(out))
	}
	return string(out), nil
}
func Revert(commitHash string) (string, error) {
	if strings.TrimSpace(commitHash) == "" {
		return "", errors.New("commit hash cannot be empty")
	}

	// 1️⃣ Validate repository
	checkRepo := exec.Command("git", "-C", state.RepoPath, "rev-parse", "--is-inside-work-tree")
	if err := checkRepo.Run(); err != nil {
		return "", errors.New("invalid Git repository path")
	}
	// Auto-stash uncommitted changes
	stashCmd := exec.Command("git", "-C", state.RepoPath, "stash", "--include-untracked")
	stashOut, _ := stashCmd.CombinedOutput()

	cmd := exec.Command("git", "-C", state.RepoPath, "revert", "--no-edit", commitHash)
	output, err := cmd.CombinedOutput()

	// Restore stash if it existed
	if !strings.Contains(string(stashOut), "No local changes") {
		exec.Command("git", "-C", state.RepoPath, "stash", "pop").Run()
	}

	// 2️⃣ Check for uncommitted changes before revert
	checkChanges := exec.Command("git", "-C", state.RepoPath, "status", "--porcelain")
	out, _ := checkChanges.Output()
	if strings.TrimSpace(string(out)) != "" {
		return "", errors.New("uncommitted changes present — please commit or stash before reverting")
	}

	return string(output), err
}
func Clone(repoPath, CloneUrl string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "clone", CloneUrl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("clone failed:%v\n%s", err, string(out))
	}

	return "successfully cloned the Repo", nil
}
