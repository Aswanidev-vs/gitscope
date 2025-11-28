package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gitscope/internal/doc"
	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/state"
)

func RullshellCommand(line string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/C", line)

	}
	return exec.Command("sh", "-c", line)
}
func NewRepoCmd(w fyne.Window, repoPath string, cmdText string) fyne.CanvasObject {
	if repoPath == "" {
		dialog.ShowError(errors.New("No repository path selected"), w)
		return nil
	}

	if err := os.MkdirAll(repoPath, 0755); err != nil {
		dialog.ShowError(fmt.Errorf("Failed to create directory: %v", err), w)
		return nil
	}

	state.RepoPath = repoPath

	progress := dialog.NewProgressInfinite("Running Commands", "Please wait while commands are executing...", w)

	go func() {
		progress.Show()

		lines := strings.Split(cmdText, "\n")
		var allErrors []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			fmt.Println("Running:", line)

			if strings.HasPrefix(line, "git commit") {
				msg := ""
				if parts := strings.SplitN(line, "-m", 2); len(parts) == 2 {
					msg = strings.Trim(parts[1], " \"'")
				}

				output, err := git.Commit(msg)
				outStr := string(output)

				// ✅ Ignore harmless commit warnings
				if strings.Contains(outStr, "no changes added to commit") ||
					strings.Contains(outStr, "no staged changes to commit") {
					fmt.Println("⚠️ Ignored: no new changes to commit.")
					continue
				}

				if err != nil {
					errMsg := fmt.Sprintf("❌ Commit failed: %v\n%s", err, output)
					fmt.Println(errMsg)
					allErrors = append(allErrors, errMsg)
				} else {
					fmt.Println("✅ Commit executed:", output)
				}
				continue
			}

			cmd := exec.Command("cmd", "/C", line)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			cmd.Dir = repoPath

			output, err := cmd.CombinedOutput()
			outStr := strings.TrimSpace(string(output))

			if err != nil {
				if strings.Contains(outStr, "already exists") ||
					strings.Contains(outStr, "nothing to commit") ||
					strings.Contains(outStr, "fatal: destination path") {
					fmt.Println("⚠️ Ignored non-fatal error:", outStr)
					continue
				}

				errMsg := fmt.Sprintf("❌ Error running: %s\nOutput:\n%s", line, outStr)
				fmt.Println(errMsg)
				allErrors = append(allErrors, errMsg)
				continue
			}

			fmt.Println("✅ Executed:", line)
			if outStr != "" {
				fmt.Println("Output:", outStr)
			}
		}

		progress.Hide()

		if len(allErrors) > 0 {
			dialog.ShowError(fmt.Errorf("Some commands failed:\n\n%s", strings.Join(allErrors, "\n\n")), w)

		} else {

			dialog.ShowInformation("Success", "All commands executed successfully.", w)

		}
	}()
	return nil
}
func ListPush(repoPath string) ([]string, error) {
	if repoPath == "" {
		return nil, errors.New("repository path cannot be empty")
	}

	cmd := exec.Command("git", "-C", repoPath, "branch", "--list")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove leading "* " from current branch
		if after, ok := strings.CutPrefix(line, "* "); ok {
			line = after
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func BranchSelector(repoPath string, w fyne.Window) (fyne.CanvasObject, func() string) {

	stop := make(chan struct{})

	selectEntry := widget.NewSelectEntry([]string{"Loading..."})
	selectEntry.SetPlaceHolder("Select a branch")

	var lastBranches []string
	var mu sync.Mutex

	updateList := func() {

		branches, err := ListPush(state.RepoPath)
		if err != nil {
			// fyne.Do(func() {
			// 	selectEntry.SetText("Error loading branches")
			// })
			fmt.Println("Branch load error:", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if !slicesEqual(branches, lastBranches) {
			lastBranches = branches
			fyne.Do(func() {
				selectEntry.SetOptions(branches)
				if len(branches) > 0 {
					selectEntry.SetText(branches[0])
				} else {
					selectEntry.SetText("")
				}
			})
		}
	}

	updateList()

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Branch watcher stopped for", repoPath)
				return
			case <-time.After(3 * time.Second):
				updateList()
			}
		}
	}()

	// stop watcher automatically when window closes
	w.SetOnClosed(func() {
		close(stop)
	})

	getSelectedBranch := func() string {
		return strings.TrimSpace(selectEntry.Text)
	}

	ui := container.NewVBox(selectEntry)
	return ui, getSelectedBranch
}

// Helper: compares two slices (order-sensitive)
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func ExistingRepoCmd(w fyne.Window, repoPath string, cmdText string) {
	if state.RepoPath == "" {
		dialog.ShowError(errors.New("No repository path selected"), w)
		return
	}

	progress := dialog.NewProgressInfinite("Executing Commands", "Please wait while running git commands...", w)
	progress.Show()

	go func() {
		defer progress.Hide()

		lines := strings.Split(cmdText, "\n")
		var allErrors []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) == 0 {
				continue
			}

			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Dir = repoPath
			out, err := cmd.CombinedOutput()
			fmt.Println("Running:", line)
			fmt.Println("Output:", string(out))

			if err != nil {
				allErrors = append(allErrors, fmt.Sprintf("❌ %s\n%s", line, string(out)))
			}
		}

		if len(allErrors) > 0 {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Git Command Errors",
				Content: strings.Join(allErrors, "\n"),
			})
			dialog.ShowError(errors.New(strings.Join(allErrors, "\n")), w)
		} else {
			dialog.ShowInformation("Success", "All commands executed successfully!", w)
		}
	}()
}
func GetPreviousCommit(repoPath string) (string, error) {
	repo := repoPath
	checkdir, err := os.Stat(repo)
	if err != nil || !checkdir.IsDir() {
		return "", errors.New("invalid directory path")
	}
	cmd := exec.Command("git", "-C", repo, "rev-parse", "HEAD~1")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("No previous commit to reset: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}
func IsInitialized(repoPath string) bool {
	gitDir := filepath.Join(repoPath, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
func Decision(title string) fyne.CanvasObject {
	switch title {
	case "Init":
		return doc.Init()
	case "Stage":
		return doc.Stage()
	case "Status":
		return doc.Status()
	case "Commit":
		return doc.Commit()
	case "Push":
		return doc.Push()
	case "Log":
		return doc.Log()
	case "Revert":
		return doc.Revert()
	case "Clone":
		return doc.Clone()
	case "Branch":
		return doc.Branch()
	case "Pull":
		return doc.Pull()
	case "Reflog":
		return doc.Reflog()
	case "GitIgnore":
		return doc.GitIgnore()
	default:
		return widget.NewLabel("Unknown Document")
	}
}
