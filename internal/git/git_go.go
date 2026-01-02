package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
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
		// Check if push failed because branch is behind remote
		if strings.Contains(string(out), "non-fast-forward") || strings.Contains(string(out), "behind its remote") {
			// Switch to the branch
			checkoutCmd := exec.Command("git", "-C", repo, "checkout", branch)
			checkoutCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			checkoutOut, checkoutErr := checkoutCmd.CombinedOutput()
			if checkoutErr != nil {
				return string(checkoutOut), fmt.Errorf("checkout failed before pull: %v\n%s", checkoutErr, string(checkoutOut))
			}
			// Pull first to integrate remote changes
			pullCmd := exec.Command("git", "-C", repo, "pull", "origin", branch, "--no-edit")
			pullCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			pullOut, pullErr := pullCmd.CombinedOutput()
			if pullErr != nil {
				return string(pullOut), fmt.Errorf("pull failed before push: %v\n%s", pullErr, string(pullOut))
			}
			// Try push again
			pushCmd := exec.Command("git", "-C", repo, "push", "-u", "origin", branch)
			pushCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			out2, err2 := pushCmd.CombinedOutput()
			if err2 != nil {
				return string(out2), fmt.Errorf("push failed after pull: %v\n%s", err2, string(out2))
			}
			return string(out2), nil
		}
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
func Clone(repoPath, cloneURL string) (string, error) {
	parentDir := filepath.Dir(repoPath)

	// Ensure parent directory exists
	if info, err := os.Stat(parentDir); err != nil || !info.IsDir() {
		return "", errors.New("invalid parent directory path")
	}

	// Prepare git clone command
	cmd := exec.Command("git", "-C", repoPath, "clone", cloneURL)

	// Hide window on Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("clone failed: %v\n%s", err, string(out))
	}

	return "successfully cloned the repo", nil
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
	// Set upstream to origin/branchname
	pushCmd := exec.Command("git", "-C", repo, "push", "-u", "origin", branchname)
	pushCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	pushOut, pushErr := pushCmd.CombinedOutput()
	if pushErr != nil {
		return string(pushOut), fmt.Errorf("Creating New Branch succeeded, but setting upstream failed:%v\n%s", pushErr, string(pushOut))
	}
	return "successfully Created New Branch and set upstream", nil

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

func SwitchBranch(repoPath, branchname string) (string, error) {
	checkdir, err := os.Stat(repoPath)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repoPath, "switch", branchname)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("switch branch failed: %v\n%s", err, string(out))
	}
	return "Switched to branch " + branchname, nil
}

func BranchRename(oldname, newname string) (string, error) {
	repo := state.RepoPath
	cmd := exec.Command("git", "-C", repo, "branch", "-m", oldname, newname)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("branch rename failed: %v\n%s", err, string(out))
	}
	return "Branch renamed from " + oldname + " to " + newname, nil
}

func GitIgnore(repoPath string, output *widget.Entry, w fyne.Window) (string, error) {
	filePath := filepath.Join(repoPath, ".gitignore")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(output.Text), 0644); err != nil {
			dialog.ShowError(err, w)
			return "", err
		}
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		dialog.ShowError(err, w)
		return "", err
	}

	output.SetText(string(content))
	return filePath, nil
}

//	func GitRemote(repo string) (string, error) {
//		cmd := exec.Command("git", "-C", repo, "remote", "-v")
//		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
//		out, err := cmd.CombinedOutput()
//		if err != nil {
//			return string(out), fmt.Errorf("git remote failed: %v\n%s", err, string(out))
//		}
//		return "failed to list git remotes for repo,select a repository", err
//	}
//
// GitRemote performs git remote operations (list, add, remove) in the specified repository directory.
func GitRemote(action string, args string) (string, error) {
	repo := state.RepoPath
	var cmd *exec.Cmd
	Rmv := func(remoteName string) error {
		name := strings.TrimSpace(remoteName)
		if name == "" {
			return fmt.Errorf("missing remote name to remove")
		}
		// Assign to the outer cmd variable
		cmd = exec.Command("git", "-C", repo, "remote", "remove", name)
		return nil
	}
	switch action {

	case "list":
		cmd = exec.Command("git", "-C", repo, "remote", "-v")

	case "remove":
		if err := Rmv(args); err != nil {
			return "", err
		}
	case "add":

		cleaned := strings.TrimSpace(args)

		prefixes := []string{
			"git remote add ",
			"remote add ",
			"add ",
		}

		for _, p := range prefixes {
			if after, ok := strings.CutPrefix(cleaned, p); ok {
				cleaned = after
				break
			}
		}

		parts := strings.Fields(cleaned)
		if len(parts) < 2 {
			return "", fmt.Errorf("usage: add <name> <url>")
		}

		name := parts[0]
		url := parts[1]
		if name == "origin" {
			_ = Rmv("origin")
		}
		cmd = exec.Command("git", "-C", repo, "remote", "add", name, url)

	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git remote %s failed: %w", action, err)
	}

	return string(output), nil
}

func Diff() (string, error) {
	repo := state.RepoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "diff")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return string(out), err
}
