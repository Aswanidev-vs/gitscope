package git

import (
	"errors"
	"os"
	"os/exec"
)

// validateRepoPath returns an error if repo is not an existing directory.
func validateRepoPath(repo string) error {
	info, err := os.Stat(repo)
	if err != nil || !info.IsDir() {
		return errors.New("invalid directory path")
	}
	return nil
}

// validateGitRepo verifies repo is an existing directory inside a Git work tree.
func validateGitRepo(repo string) error {
	if err := validateRepoPath(repo); err != nil {
		return err
	}
	cmd := exec.Command("git", "-C", repo, "rev-parse", "--is-inside-work-tree")
	hideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return errors.New("invalid Git repository path")
	}
	return nil
}
