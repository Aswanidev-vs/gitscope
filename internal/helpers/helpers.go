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

				output, err := git.Commit(msg, "Standard")
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
			cmd.Stdout = nil
			cmd.Stderr = nil
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

		fyne.Do(func() { progress.Hide() })

		if len(allErrors) > 0 {
			fyne.Do(func() {
				dialog.ShowError(fmt.Errorf("Some commands failed:\n\n%s", strings.Join(allErrors, "\n\n")), w)
			})

		} else {
			fyne.Do(func() {
				dialog.ShowInformation("Success", "All commands executed successfully.", w)
			})

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
		if state.RepoPath == "" {
			return
		}

		branches, err := ListPush(state.RepoPath)
		if err != nil {
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

// OptionSelector creates a UI component for selecting from a list of static options.
func OptionSelector(options []string) (fyne.CanvasObject, func() string) {
	selected := options[0] // Default to the first option
	selectWidget := widget.NewSelect(options, func(s string) {
		selected = s
	})
	selectWidget.SetSelected(selected)

	getSelectedOption := func() string {
		return selected
	}

	return selectWidget, getSelectedOption
}

func ExistingRepoCmd(w fyne.Window, repoPath string, cmdText string) {
	if repoPath == "" {
		dialog.ShowError(errors.New("no repository path selected"), w)
		return
	}

	// Validate repository path exists
	if _, err := os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			dialog.ShowError(fmt.Errorf("repository path does not exist: %s", repoPath), w)
		} else {
			dialog.ShowError(fmt.Errorf("could not access path: %v", err), w)
		}
		return
	}

	// Validate that the path is an initialized Git repository
	if !IsInitialized(repoPath) {
		dialog.ShowError(fmt.Errorf("path is not a git repository: %s", repoPath), w)
		return
	}

	// Find Git executable path
	gitPath := getGitPath()
	if gitPath == "" {
		dialog.ShowError(errors.New("Git is not installed or not found in PATH. Please install Git and try again."), w)
		return
	}

	// Initialize progress bar
	progress := dialog.NewProgressInfinite("Running Commands", "Please wait while git commands execute...", w)

	go func() {
		fyne.Do(func() { progress.Show() })

		lines := strings.Split(cmdText, "\n")
		var allErrors []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			fmt.Println("Running command:", line)

			var cmd *exec.Cmd
			if strings.HasPrefix(line, "git ") {
				// Parse Git command
				args := strings.Fields(line)
				if len(args) > 1 {
					cmd = exec.Command(gitPath, args[1:]...)
				} else {
					cmd = exec.Command(gitPath)
				}
			} else {
				// For non-Git commands, use cmd /C
				cmd = exec.Command("cmd", "/C", line)
			}

			// CRITICAL: Set the working directory to the repository path
			cmd.Dir = repoPath
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

			// Run the command and capture combined output (stdout + stderr)
			output, err := cmd.CombinedOutput()
			outStr := strings.TrimSpace(string(output))
			fmt.Printf("Command output: %q, Error: %v\n", outStr, err)

			if err != nil {
				// Ignore non-fatal Git errors
				if strings.Contains(outStr, "already exists") ||
					strings.Contains(outStr, "nothing to commit") ||
					strings.Contains(outStr, "fatal: destination path") ||
					strings.Contains(outStr, "has no upstream branch") ||
					strings.Contains(outStr, "Everything up-to-date") ||
					strings.Contains(outStr, "src refspec main does not match any") {
					fmt.Println("⚠️ Ignored non-fatal error:", outStr)
					continue
				}

				errorMsg := fmt.Sprintf("Command: %s\nError: %v\nOutput: %s", line, err, outStr)
				allErrors = append(allErrors, errorMsg)
				// Optional: break if a command fails
				// break
			} else {
				fmt.Println("✅ Executed:", line)
				if outStr != "" {
					fmt.Println("Output:", outStr)
				}
			}
		}

		// After the command loop finishes...
		fyne.Do(func() {
			progress.Hide() // Hide progress immediately on UI thread

			if len(allErrors) > 0 {
				fmt.Println("❌ Some commands failed.")
				errorDetail := strings.Join(allErrors, "\n\n")
				dialog.ShowError(
					fmt.Errorf("The following commands failed:\n\n%s", errorDetail),
					w,
				)
			} else {
				fmt.Println("✅ All commands executed successfully.")
				dialog.ShowInformation(
					"Success",
					"All commands executed successfully.",
					w,
				)
			}
		})
	}()
}
func GetPreviousCommit() (string, error) {
	repo := state.RepoPath
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

func getGitPath() string {
	// Try common paths on Windows
	paths := []string{
		"C:\\Program Files\\Git\\bin\\git.exe",
		"C:\\Program Files (x86)\\Git\\bin\\git.exe",
		"C:\\Git\\bin\\git.exe",
		"git", // if in PATH
	}
	for _, p := range paths {
		if p == "git" {
			// Check if git is in PATH
			cmd := exec.Command("where", "git")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			if err := cmd.Run(); err == nil {
				return "git"
			}
		} else {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	}
	return ""
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
	case "Remote":
		return doc.Remote()
	case "Diff":
		return doc.Diff()
	case "Reset":
		return doc.Reset()
	case "Fetch":
		return doc.Fetch()
	case "Stash":
		return doc.Stash()
	case "Merge":
		return doc.Merge()
	case "Tag":
		return doc.Tag()
	case "Cherry-pick":
		return doc.CherryPick()
	default:
		return widget.NewLabel("Unknown Document")
	}
}
