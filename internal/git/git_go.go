package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gitscope/internal/state"
)

func Init() (string, error) {
	out, err := exec.Command("git", "init").Output()
	return string(out), err
}
func Status() (string, error) {
	cmd := exec.Command("git", "-C", state.RepoPath, "status")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
func Commit(msg string) (string, error) {
	if strings.TrimSpace(msg) == "" {
		return "", errors.New("commit message cannot be empty")
	}

	// Check if repo path is valid
	checkCmd := exec.Command("git", "-C", state.RepoPath, "rev-parse", "--is-inside-work-tree")
	if err := checkCmd.Run(); err != nil {
		return "", errors.New("invalid Git repository path")
	}

	// Check if there are staged changes
	statusCmd := exec.Command("git", "-C", state.RepoPath, "diff", "--cached", "--quiet")
	if err := statusCmd.Run(); err == nil {
		return "", errors.New("no staged changes to commit")
	}

	// Proceed with commit
	cmd := exec.Command("git", "-C", state.RepoPath, "commit", "-m", msg)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Stage adds all modified/untracked files to the Git index (staging area)
func Stage() (string, error) {
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

func Revert(repoPath, sha string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "revert", "--no-edit", sha)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Return the actual Git output directly for clarity
		return string(out), fmt.Errorf("%v\n%s", err, string(out))
	}
	stage := exec.Command("git", "-C", state.RepoPath, "add", ".")
	outs, err := stage.CombinedOutput()
	return string(outs), nil
}
