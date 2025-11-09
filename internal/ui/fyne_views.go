package ui

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	core "github.com/gitscope/internal/core"
	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/helpers"
	"github.com/gitscope/internal/state"
	"github.com/gitscope/utils"
)

func RepositoryPage(w fyne.Window) fyne.CanvasObject {
	output := widget.NewLabel("Current repo: " + state.RepoPath)

	browseBtn := widget.NewButton("Select Repository", func() {

		core.OpenRepo(w, func(msg string) {
			// output.SetText("Current repo: " + state.RepoPath)
			output.SetText("" + state.RepoPath)

		})
	})

	browseBtn.Resize(fyne.NewSize(100, 40))
	browseBtn.Move(fyne.NewPos(1, 50))

	NewRepo := core.CreateNewRepo(w, func(s string) {
		dialog.ShowInformation("Repo command ", s, w)
	})
	// NewRepo.Resize(fyne.NewSize(400, 500))
	// NewRepo.Move(fyne.NewPos(1, 500))

	ExistRepo := core.ExistingRepo(w, func(cmdText string) {
		helpers.ExistingRepoCmd(w, state.RepoPath, cmdText)
	})

	// ExistRepo.Resize(fyne.NewSize(250, 300))
	// ExistRepo.Move(fyne.NewPos(1, 600))
	return container.NewVBox(
		widget.NewLabel("Repository Setup"),
		browseBtn,
		output,
		NewRepo,
		ExistRepo,
	)
}

func dashBoardPage(w fyne.Window) fyne.CanvasObject {
	output := widget.NewMultiLineEntry()

	output.Resize(fyne.NewSize(500, 230))
	output.Move(fyne.NewPos(0, 1))
	clearBtn := widget.NewButton("Clear", func() {
		output.SetText("") // Clears the content
	})
	clearBtn.Resize(fyne.NewSize(100, 40))
	clearBtn.Move(fyne.NewPos(510, 1))

	initBtn := InitButton(output)
	initBtn.Resize(fyne.NewSize(100, 40))
	initBtn.Move(fyne.NewPos(1, 250))

	stageBtn := StageButton(output)
	stageBtn.Resize(fyne.NewSize(100, 40))
	stageBtn.Move(fyne.NewPos(110, 250))

	statusBtn := StatusButton(output)
	statusBtn.Resize(fyne.NewSize(100, 40))
	statusBtn.Move(fyne.NewPos(220, 250))

	commitBtn := CommitButton(w)
	commitBtn.Resize(fyne.NewSize(100, 40))
	commitBtn.Move(fyne.NewPos(329, 250))

	pushBtn := PushButton(w)
	pushBtn.Resize(fyne.NewSize(100, 40))
	pushBtn.Move(fyne.NewPos(439, 250))

	logBtn := LogButton(output)
	logBtn.Resize(fyne.NewSize(100, 40))
	logBtn.Move(fyne.NewPos(1, 350))

	revertBtn := RevertButton(w)
	revertBtn.Resize(fyne.NewSize(100, 40))
	revertBtn.Move(fyne.NewPos(110, 350))

	cloneBtn := CloneButton(w)
	cloneBtn.Resize(fyne.NewSize(100, 40))
	cloneBtn.Move(fyne.NewPos(220, 350))

	Branchbtn := BranchButton(w)
	Branchbtn.Resize(fyne.NewSize(100, 40))
	Branchbtn.Move(fyne.NewPos(329, 350))

	PullBtn := PullButton(w)
	PullBtn.Resize(fyne.NewSize(100, 40))
	PullBtn.Move(fyne.NewPos(439, 350))

	return container.NewWithoutLayout(initBtn, stageBtn, commitBtn, statusBtn, pushBtn, logBtn, revertBtn, cloneBtn, Branchbtn, PullBtn, clearBtn, output)
}
func InitButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Init", func() {
		out, err := git.Init()
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}

func StatusButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Status", func() {
		out, err := git.Status()
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}

func StageButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Stage", func() {
		out, err := git.Stage()
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}

func CommitButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Commit", func() {
		input := widget.NewEntry()
		form := []*widget.FormItem{
			{Text: "Message", Widget: input},
		}
		dialog.ShowForm("Enter your commit message", "Commit", "Cancel", form, func(valid bool) {
			if valid {
				msg := input.Text
				if msg == "" {
					dialog.ShowInformation("Empty Message", "Commit message cannot be empty", w)
					return
				}
				out, err := git.Commit(msg)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Commit Result", out, w)
				}
			}
		}, w)
	})
}
func PushButton(w fyne.Window) fyne.CanvasObject {

	branchSelectorUI, getBranch := helpers.BranchSelector(state.RepoPath, w)
	pushBtn := widget.NewButton("Push", func() {
		if state.RepoPath == "" {
			dialog.ShowError(errors.New("No repository selected"), w)
			return
		}
		branch := getBranch()
		if branch == "" {
			dialog.ShowError(errors.New("No branch selected"), w)
			return
		}

		progress := dialog.NewProgressInfinite("Running Commands", "Please wait while commands are executing...", w)

		go func() {
			progress.Show()
			output, err := git.Push(state.RepoPath, branch)
			progress.Hide()
			dialog.ShowInformation("Push Success", "Repository pushed successfully.", w)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Push failed:\n%v\n\n%s", err, output), w)
				return
			}

		}()
	})

	return container.NewVBox(
		pushBtn,
		branchSelectorUI,
	)
}

func LogButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Log", func() {
		out, err := git.Log(state.RepoPath)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}
func RevertButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Revert", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Missing Clone Destination", "Please select or create an empty folder before cloning the repository.", w)
			return
		}
		input := widget.NewEntry()
		input.SetPlaceHolder("e.g. a1s4fd6")
		input.Resize(fyne.NewSize(350, 40))

		form := []*widget.FormItem{
			{Widget: input},
		}

		dialog.ShowForm("Revert Commit hash", "Revert", "Cancel", form, func(valid bool) {
			if !valid {
				return
			}

			sha := strings.TrimSpace(input.Text)
			if sha == "" {
				dialog.ShowInformation("Empty SHA", "Commit hash cannot be empty.", w)
				return
			}

			out, err := git.Revert(sha)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			dialog.ShowInformation("Revert Successful", out, w)
		}, w)
	})
}
func CloneButton(w fyne.Window) *widget.Button {
	return widget.NewButton("Clone", func() {
		// Sync local repoPath with global state
		if state.RepoPath == "" {
			dialog.ShowInformation("Missing Clone Destination", "Please select or create an empty folder before cloning the repository.", w)
			return
		}

		// Optional: validate that folder is truly empty
		files, err := os.ReadDir(state.RepoPath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Unable to read target folder: %w", err), w)
			return
		}
		if len(files) > 0 {
			dialog.ShowInformation("Folder Not Empty", "Please choose an empty folder to clone into.", w)
			return
		}

		// Clone URL input
		input := widget.NewEntry()
		input.SetPlaceHolder("https://github.com/yourname/repositoryname.git")
		input.Resize(fyne.NewSize(350, 40))

		form := []*widget.FormItem{
			{Widget: input},
		}

		dialog.ShowForm("Clone Repository", "Clone", "Cancel", form, func(valid bool) {
			if !valid {
				return
			}

			cloneurl := strings.TrimSpace(input.Text)
			if cloneurl == "" {
				dialog.ShowInformation("Missing URL", "Repository URL cannot be empty.", w)
				return
			}
			if !strings.HasPrefix(cloneurl, "https://github.com/") {
				dialog.ShowInformation("Invalid URL", "Please enter a valid GitHub repository URL.", w)
				return
			}

			// out, err := git.Clone(state.RepoPath, cloneurl)
			// if err != nil {
			// 	dialog.ShowError(fmt.Errorf("Clone failed: %w", err), w)
			// 	return
			// }

			// dialog.ShowInformation("Clone Successful", out, w)
		}, w)
	})
}
func BranchButton(w fyne.Window) *widget.Button {
	btn := widget.NewButton("Branch", func() {
		// Check if repoPath is set
		if state.RepoPath == "" {
			dialog.ShowInformation("No Repository Found", "Please open or create a repository first.", w)
			return
		}

		// Check if repoPath folder exists
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowInformation("Invalid Repository", "The selected repository path does not exist. Please create or open a valid one.", w)
			return
		}

		// Create the branch button dialog
		dlg := utils.NewBranchButton(
			w,
			func(name string) {
				_, err := git.CreateBranch(state.RepoPath, name)
				if err != nil {
					dialog.ShowError(err, w)
				}
				// else {
				// 	dialog.ShowInformation("Branch Created", out, w)
				// }
			},
			func(name string) {
				_, err := git.DeleteBranch(state.RepoPath, name)
				if err != nil {
					dialog.ShowError(err, w)
				}
				// else {
				// 	dialog.ShowInformation("Branch Deleted", out, w)
				// }
			},
		)

		// Simulate clicking to open the dialog
		dlg.Tapped(nil)
	})

	return btn
}

func PullButton(w fyne.Window) fyne.CanvasObject {
	branchSelectorUI, getBranch := helpers.BranchSelector(state.RepoPath, w)

	pullBtn := widget.NewButton("Pull", func() {
		if state.RepoPath == "" {
			dialog.ShowError(errors.New("No repository selected"), w)
			return
		}
		branch := getBranch()
		if branch == "" {
			dialog.ShowError(errors.New("No branch selected"), w)
			return
		}

		// Ask developer if they also want to reset before pulling
		dialog.ShowConfirm("Pull Options", "Do you want to reset the last commit before pulling?", func(reset bool) {
			progress := dialog.NewProgressInfinite("Running Commands", "Please wait while commands are executing...", w)
			go func() {
				progress.Show()
				defer progress.Hide()

				if reset {
					// Perform optional reset
					sha, err := git.GetPreviousCommit(state.RepoPath)
					if err != nil {
						dialog.ShowError(err, w)
						return
					}

					resetCmd := exec.Command("git", "-C", state.RepoPath, "reset", "--soft", sha)
					out, err := resetCmd.CombinedOutput()
					if err != nil {
						dialog.ShowError(fmt.Errorf("Reset failed:\n%v\n\n%s", err, string(out)), w)
						return
					}

				}

				// Always perform pull
				output, err := git.Pull(state.RepoPath, branch)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Pull failed:\n%v\n\n%s", err, output), w)
					return
				}

				msg := "Pull completed successfully."
				if reset {
					msg = "Last commit reset and pull completed successfully."
				}
				dialog.ShowInformation("Success", msg, w)
			}()
		}, w)
	})

	return container.NewVBox(
		pullBtn,
		branchSelectorUI,
	)
}

func SettingPage(w fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("assets/icons/gitscope_logo_v6.png")
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(120, 120))

	// About text (professional tone and cleaner formatting)
	f1 := widget.NewLabel("GitScope is a modern, lightweight, and visually intuitive Git client built with Go and Fyne. It simplifies essential")
	f2 := widget.NewLabel("version control operations making Git easier to use for both beginners and experienced developers.")
	f3 := widget.NewLabel("Version: 1.0.0")
	f4 := widget.NewLabel("Developer: Aswanidev VS")

	// GitHub link
	link := widget.NewHyperlink("🔗 View Project on GitHub", &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "Aswanidev-vs/GitScope",
	})

	// Foreground layout (content)
	content := container.NewVBox(
		widget.NewLabelWithStyle("About GitScope", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		f1,
		f2,
		f3,
		f4,
		link,
	)
	centeredContent := container.NewCenter(content)

	return container.NewStack(centeredContent)
}
