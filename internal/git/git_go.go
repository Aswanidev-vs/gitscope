package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gitscope/internal/state"
)

func Init() (string, error) {
	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}

	cmd := exec.Command("git", "-C", repo, "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return string(out), err
}
func Status() (string, error) {
	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}

	cmd := exec.Command("git", "-C", repo, "status")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Commit creates a new commit with the given message.
func Commit(msg string) (string, error) {
	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}

	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "", errors.New("commit message cannot be empty")
	}

	// Validate repo path
	checkCmd := exec.Command("git", "-C", repo, "rev-parse", "--is-inside-work-tree")
	checkCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := checkCmd.Run(); err != nil {
		return "", errors.New("invalid Git repository path")
	}

	// Check for staged changes
	statusCmd := exec.Command("git", "-C", repo, "diff", "--cached", "--quiet")
	statusCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := statusCmd.Run(); err == nil {
		return "", errors.New("no staged changes to commit")
	}

	// Proceed with commit
	cmd := exec.Command("git", "-C", repo, "commit", "-m", msg)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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

	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	// if state.RepoPath == "" {
	// 	return "", errors.New("no repository path selected")
	// }
	cmd := exec.Command("git", "-C", repo, "add", ".")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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
	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}

	cmd := exec.Command("git", "-C", repo, "push", "-u", "origin", branch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("push failed: %v\n%s", err, string(out))
	}
	return string(out), nil
}
func Log(repoPath string) (string, error) {
	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "log", "--oneline")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), fmt.Errorf("log failed:%v\n%s", err, string(out))
	}
	return string(out), nil
}
func Revert(commitHash string) (string, error) {

	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}

	if strings.TrimSpace(commitHash) == "" {
		return "", errors.New("commit hash cannot be empty")
	}

	// 1️⃣ Validate repository
	checkRepo := exec.Command("git", "-C", repo, "rev-parse", "--is-inside-work-tree")
	checkRepo.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := checkRepo.Run(); err != nil {
		return "", errors.New("invalid Git repository path")
	}
	// Auto-stash uncommitted changes
	stashCmd := exec.Command("git", "-C", repo, "stash", "--include-untracked")
	stashCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stashOut, _ := stashCmd.CombinedOutput()

	cmd := exec.Command("git", "-C", repo, "revert", "--no-edit", commitHash)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()

	// Restore stash if it existed
	if !strings.Contains(string(stashOut), "No local changes") {
		exec.Command("git", "-C", repo, "stash", "pop").Run()
	}

	// 2️⃣ Check for uncommitted changes before revert
	checkChanges := exec.Command("git", "-C", repo, "status", "--porcelain")
	checkChanges.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, _ := checkChanges.Output()
	if strings.TrimSpace(string(out)) != "" {
		return "", errors.New("uncommitted changes present — please commit or stash before reverting")
	}

	return string(output), err
}
func Clone(repoPath, CloneUrl string) (string, error) {
	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "clone", CloneUrl)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()

	if err != nil {
		return string(out), fmt.Errorf("clone failed:%v\n%s", err, string(out))
	}

	return "successfully cloned the Repo", nil
}
func CreateBranch(repoPath, branchname string) (string, error) {

	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "branch", branchname)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Creating New Branch failed:%v\n%s", err, string(out))
	}
	return "successfully Created New Branch", nil

}
func DeleteBranch(repoPath, branchname string) (string, error) {
	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "branch", "-d", branchname)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Creating New Branch failed:%v\n%s", err, string(out))
	}
	return "successfully Deleted New Branch", nil
}
func Pull(repoPath, branch string) (string, error) {

	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "pull", "origin", branch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("An issue occurred while pulling: %v\n%s", err, string(out))
	}
	return "Successfully pulled branch: " + branch, nil
}

func Reflog(repoPath string) (string, error) {
	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "reflog")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("An issue occurred while reflog: %v\n%s", err, string(out))
	}
	return string(out), nil
}
