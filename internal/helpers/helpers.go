package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gitscope/internal/git"

	"github.com/gitscope/internal/state"
)

var loadingDialog dialog.Dialog

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
	cmd := exec.Command("git", "branch", "--list")
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

func BranchSelector(repoPath string) (fyne.CanvasObject, func() string) {
	selectEntry := widget.NewSelectEntry([]string{"Loading..."})
	selectEntry.SetPlaceHolder("Select a branch")

	go func() {
		branches, err := ListPush(repoPath)
		if err != nil {
			selectEntry.SetText("Error loading branches")
			fmt.Println("Branch load error:", err)
			return
		}

		selectEntry.SetOptions(branches)
		if len(branches) > 0 {
			selectEntry.SetText(branches[0])
		}
	}()

	getSelectedBranch := func() string {
		return strings.TrimSpace(selectEntry.Text)
	}

	ui := container.NewVBox(
		// widget.NewLabel("Available Branches"),
		selectEntry,
	)

	return ui, getSelectedBranch
}
