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

func Stage() (string, error) {
	cmd := exec.Command("git", "-C", state.RepoPath, "add", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return "Files staged successfully.", nil
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
func Revert(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "revert", "--no-edit")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Revert failed:%v\n%s", err, string(out))
	}
	return string(out), nil
}
