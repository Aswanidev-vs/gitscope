package ui

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	core "github.com/gitscope/internal/core"
	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/helpers"
	"github.com/gitscope/internal/state"
	"github.com/gitscope/utils"
)

var (
	edit      bool
	gitignore string
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
	output.SetPlaceHolder(`This area shows the output / responses for majority of the commands
that are triggered by the buttons. You can also add your gitignore
entries here.I meant .ext as a short form to refer
to any extension like .go, .html, .css, etc.

*.ext            → Ignore all .ext files everywhere
/*.ext           → Only ignore .ext files in repo root
folder/          → Ignore entire folder
**/name.ext      → Ignore file no matter where it appears
`)

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

	statusBtn_comp := StatusButton(output)
	statusBtn_comp.Resize(fyne.NewSize(100, 80))
	statusBtn_comp.Move(fyne.NewPos(220, 250))

	commitBtn_comp := CommitButton(w)
	commitBtn_comp.Resize(fyne.NewSize(100, 80))
	commitBtn_comp.Move(fyne.NewPos(329, 250))

	pushBtn := PushButton(w)
	pushBtn.Resize(fyne.NewSize(100, 40))
	pushBtn.Move(fyne.NewPos(439, 250))

	logBtn_comp := LogButton(output)
	logBtn_comp.Resize(fyne.NewSize(100, 80))
	logBtn_comp.Move(fyne.NewPos(1, 350))

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

	Reflogbtn := ReflogButton(w, output)
	Reflogbtn.Resize(fyne.NewSize(100, 40))
	Reflogbtn.Move(fyne.NewPos(1, 450))

	SwitchBranchBtn := SwitchBranchButton(w)
	SwitchBranchBtn.Resize(fyne.NewSize(100, 40))
	SwitchBranchBtn.Move(fyne.NewPos(110, 450))

	BranchRenameBtn := BranchRenameButton(w, output)
	BranchRenameBtn.Resize(fyne.NewSize(100, 40))
	BranchRenameBtn.Move(fyne.NewPos(220, 450))

	GitIgnoreBtn := GitIgnoreButton(output, w)
	GitIgnoreBtn.Resize(fyne.NewSize(110, 40))
	GitIgnoreBtn.Move(fyne.NewPos(510, 195))

	GitRemotebtn := RemoteButton(w, output)
	GitRemotebtn.Resize(fyne.NewSize(100, 40))
	GitRemotebtn.Move(fyne.NewPos(329, 450))

	Diffbtn_comp := DiffButton(output)
	Diffbtn_comp.Resize(fyne.NewSize(100, 80))
	Diffbtn_comp.Move(fyne.NewPos(439, 450))

	Resetbtn_comp := ResetButton(output, w)
	Resetbtn_comp.Resize(fyne.NewSize(100, 80))
	Resetbtn_comp.Move(fyne.NewPos(1, 550))

	magicSyncBtn := MagicSyncButton(output)
	magicSyncBtn.Resize(fyne.NewSize(100, 40))
	magicSyncBtn.Move(fyne.NewPos(110, 550))

	undoBtn := UndoButton(output, w)
	undoBtn.Resize(fyne.NewSize(100, 40))
	undoBtn.Move(fyne.NewPos(220, 550))

	cherryPickBtn := CherryPickButton(output, w)
	cherryPickBtn.Resize(fyne.NewSize(100, 40))
	cherryPickBtn.Move(fyne.NewPos(329, 550))

	conflictBtn := ConflictButton(output, w)
	conflictBtn.Resize(fyne.NewSize(100, 40))
	conflictBtn.Move(fyne.NewPos(439, 550))

	tagBtn_comp := TagButton(output, w)
	tagBtn_comp.Resize(fyne.NewSize(100, 80))
	tagBtn_comp.Move(fyne.NewPos(1, 650))

	fetchBtn_comp := FetchButton(output, w)
	fetchBtn_comp.Resize(fyne.NewSize(100, 80))
	fetchBtn_comp.Move(fyne.NewPos(110, 650))

	stashBtn_comp := StashButton(output, w)
	stashBtn_comp.Resize(fyne.NewSize(100, 80))
	stashBtn_comp.Move(fyne.NewPos(220, 650))

	mergeBtn_comp := MergeButton(output, w)
	mergeBtn_comp.Resize(fyne.NewSize(100, 80))
	mergeBtn_comp.Move(fyne.NewPos(329, 650))

	return container.NewWithoutLayout(initBtn, stageBtn, commitBtn_comp, statusBtn_comp, pushBtn, logBtn_comp, revertBtn, cloneBtn, Branchbtn, PullBtn, clearBtn, Reflogbtn, SwitchBranchBtn, BranchRenameBtn, GitIgnoreBtn, GitRemotebtn, Diffbtn_comp, Resetbtn_comp, magicSyncBtn, undoBtn, cherryPickBtn, conflictBtn, tagBtn_comp, fetchBtn_comp, stashBtn_comp, mergeBtn_comp, output)
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

func StatusButton(output *widget.Entry) fyne.CanvasObject {
	options := []string{"Standard", "Short (-s)", "Branch (-b)"}
	statusSelect := widget.NewSelect(options, func(value string) {})
	statusSelect.SetSelected("Standard")

	statusBtn := widget.NewButton("Status", func() {
		out, err := git.Status(statusSelect.Selected)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
	return container.NewVBox(statusBtn, statusSelect)
}

func StageButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Stage", func() {
		out, err := git.Stage("All (.)")
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}

func CommitButton(w fyne.Window) fyne.CanvasObject {
	options := []string{"Standard (-m)", "Stage All (-a)", "Amend (--amend)"}
	commitSelect := widget.NewSelect(options, func(value string) {})
	commitSelect.SetSelected("Standard (-m)")

	commitBtn := widget.NewButton("Commit", func() {
		repo := state.RepoPath
		checkdir, err := os.Stat(repo)
		if err != nil || !checkdir.IsDir() {
			dialog.ShowInformation("invalid repository path", "Please select a valid repository path before commit.", w)
			return
		}
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
				out, err := git.Commit(msg, commitSelect.Selected)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Commit Result", out, w)
				}
			}
		}, w)
	})
	return container.NewVBox(commitBtn, commitSelect)
}

func PushButton(w fyne.Window) fyne.CanvasObject {
	branchSelectorUI, getBranch := helpers.BranchSelector(state.RepoPath, w)

	pushBtn := widget.NewButton("Push", func() {
		repoPath := state.RepoPath
		if repoPath == "" {
			dialog.ShowError(errors.New("No repository selected."), w)
			return
		}

		// 1️⃣ Check if initialized
		if !helpers.IsInitialized(repoPath) {
			dialog.ShowInformation("Git Initialization", "Repository is not initialized.\nPlease run git init first.", w)
			return
		}

		branch := getBranch()
		if branch == "" {
			dialog.ShowError(errors.New("No branch selected."), w)
			return
		}

		// 2️⃣ Check if there are unstaged files
		statusCmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
		statusCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := statusCmd.Output()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Failed to check git status: %v", err), w)
			return
		}

		hasUnstaged := false
		hasStaged := false

		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}
			if line[0] == ' ' && line[1] != ' ' {
				hasUnstaged = true // modified but not staged
			}
			if line[0] != ' ' {
				hasStaged = true // staged file
			}
		}

		if hasUnstaged && !hasStaged {
			dialog.ShowInformation("Stage Required", "You have unstaged changes.\nPlease stage them before committing.", w)
			return
		}

		// 3️⃣ Check if staged but not committed
		diffCmd := exec.Command("git", "-C", repoPath, "diff", "--cached", "--name-only")
		diffCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		diffOut, _ := diffCmd.Output()
		if len(strings.TrimSpace(string(diffOut))) > 0 {
			dialog.ShowInformation("Commit Required", "You have staged files but no commit yet.\nPlease commit before pushing.", w)
			return
		}

		// // 4️⃣ Check if there are commits to push
		// cherryCmd := exec.Command("git", "-C", repoPath, "cherry", "-v")
		// cherryCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		// cherryOut, _ := cherryCmd.Output()
		// if len(strings.TrimSpace(string(cherryOut))) == 0 {
		// 	dialog.ShowInformation("No Commits to Push", "No new commits to push.\nMake a commit first.", w)
		// 	return
		// }

		// 5️⃣ Push if all good
		progress := dialog.NewProgressInfinite("Pushing Repository", "Please wait...", w)
		go func() {
			fyne.Do(func() { progress.Show() })
			output, err := git.Push(repoPath, branch)
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(fmt.Errorf("Push failed:\n%v\n\n%s", err, output), w)
					return
				}
				dialog.ShowInformation("Push Success", "Repository pushed successfully.", w)
			})
		}()
	})

	return container.NewVBox(pushBtn, branchSelectorUI)
}
func LogButton(output *widget.Entry) fyne.CanvasObject {
	options := []string{"Oneline", "Graph", "Pretty"}
	logSelect := widget.NewSelect(options, func(value string) {})
	logSelect.SetSelected("Oneline")

	logBtn := widget.NewButton("Log", func() {
		out, err := git.Log(state.RepoPath, logSelect.Selected)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
	return container.NewVBox(logBtn, logSelect)
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

			out, err := git.Revert(sha, "--no-edit")
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

		// files, err := os.ReadDir(state.RepoPath)
		_, err := os.ReadDir(state.RepoPath)

		if err != nil {
			dialog.ShowError(fmt.Errorf("Unable to read target folder: %w", err), w)
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

			out, err := git.Clone(state.RepoPath, cloneurl)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Clone failed: %w", err), w)
				return
			}

			dialog.ShowInformation("Clone Successful", out, w)
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

		dialog.ShowConfirm("Pull Options", "Do you want to reset the last commit before pulling?", func(reset bool) {
			progress := dialog.NewProgressInfinite("Running Commands", "Please wait while commands are executing...", w)
			go func() {
				progress.Show()
				defer progress.Hide()

				if reset {

					sha, err := helpers.GetPreviousCommit()
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

					git.Stage("All (.)")
				}

				output, err := git.Pull(state.RepoPath, branch)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Pull failed:\n%v\n\n%s", err, output), w)
					return
				}
				// msg := "Pull completed successfully."
				// if reset {
				// 	msg = "Last commit reset and pull completed successfully."
				// }
				// dialog.ShowInformation("Success", msg, w)
				// fyne.Do(func() {
				// 	dialog.ShowInformation("Success", msg, w)
				// })

				fyne.Do(func() {
					progress.Hide()

					if err != nil {
						dialog.ShowError(err, w)
						return
					}

					msg := "Pull completed successfully."
					if reset {
						msg = "Last commit reset and pull completed successfully."
					}

					dialog.ShowInformation("Success", msg, w)
				})

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

	f1 := widget.NewLabel("GitScope is a modern, lightweight, and visually intuitive Git client built with Go and Fyne. It simplifies essential")
	f2 := widget.NewLabel("version control operations making Git easier to use for both beginners and experienced developers.")
	f3 := widget.NewLabel("Version: 1.1.3")
	f4 := widget.NewLabel("Developer: Aswanidev VS")

	link := widget.NewHyperlink("🔗 View Project on GitHub", &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "Aswanidev-vs/GitScope",
	})

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
func ReflogButton(w fyne.Window, output *widget.Entry) *widget.Button {
	return widget.NewButton("Reflog", func() {

		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Cannot show the reflog because no Git repository has been selected. Please choose a repository and try again.", w)
			return
		}
		out, err := git.Reflog(state.RepoPath, "Standard")
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
}

func SwitchBranchButton(w fyne.Window) fyne.CanvasObject {
	branchSelectorUI, getBranch := helpers.BranchSelector(state.RepoPath, w)

	switchBtn := widget.NewButton("Switch Branch", func() {
		if state.RepoPath == "" {
			dialog.ShowError(errors.New("No repository selected"), w)
			return
		}
		branch := getBranch()
		if branch == "" {
			dialog.ShowError(errors.New("No branch selected"), w)
			return
		}

		output, err := git.SwitchBranch(state.RepoPath, branch)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Switch branch failed:\n%v\n\n%s", err, output), w)
			return
		}
		dialog.ShowInformation("Success", "Switched to branch: "+branch, w)
	})

	return container.NewVBox(switchBtn, branchSelectorUI)
}

func BranchRenameButton(w fyne.Window, output *widget.Entry) *widget.Button {
	return widget.NewButton(`Rename
Branch`, func() {
		oldInput := widget.NewEntry()
		oldInput.SetPlaceHolder("Current branch name")
		newInput := widget.NewEntry()
		newInput.SetPlaceHolder("New branch name")

		form := []*widget.FormItem{
			{Widget: oldInput},
			{Widget: newInput},
		}
		repo := state.RepoPath
		checkdir, err := os.Stat(repo)
		if err != nil || !checkdir.IsDir() {
			dialog.ShowInformation("Error", "invalid directory path", w)
			return
		}
		dialog.ShowForm("Rename Branch", "Rename", "Cancel", form, func(valid bool) {
			if valid {
				oldname := strings.TrimSpace(oldInput.Text)
				newname := strings.TrimSpace(newInput.Text)
				if oldname == "" || newname == "" {
					dialog.ShowInformation("Error", "Both old and new branch names are required", w)
					return
				}
				out, err := git.BranchRename(oldname, newname)
				if err != nil {
					output.SetText("error: " + err.Error())
				} else {
					output.SetText(out)
				}
			}
		}, w)
	})

}
func GitIgnoreButton(output *widget.Entry, w fyne.Window) *widget.Button {

	btn := widget.NewButton("gitignore(Edit)", nil)

	btn.OnTapped = func() {
		if state.RepoPath == "" {
			dialog.ShowError(errors.New("No repository selected"), w)
			return
		}
		if edit {
			err := os.WriteFile(gitignore, []byte(output.Text), 0644)
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("Saved", ".gitignore updated successfully", w)
			}
			edit = false
			btn.SetText("gitignore(Edit)")
			return
		}

		path, err := git.GitIgnore(state.RepoPath, output, w)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		edit = true
		gitignore = path
		btn.SetText(".gitignore (Save)")
	}

	return btn
}
func DocumentPage(w fyne.Window) fyne.CanvasObject {

	items := []string{"Init", "Stage", "Status", "Commit", "Push", "Log", "Revert", "Clone", "Branch", "Pull", "Reflog", "GitIgnore", "Remote", "Diff", "Reset", "Fetch", "Stash", "Merge", "Tag", "Cherry-pick"}

	masterContainer := container.NewStack()

	listContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(420, 60)))

	var listView fyne.CanvasObject

	for _, item := range items {
		name := item
		btn := widget.NewButton(name, func() {

			docPage := helpers.Decision(name)

			backBtn := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
				masterContainer.Objects = []fyne.CanvasObject{listView}
				masterContainer.Refresh()
			})

			detailPage := container.NewBorder(
				container.NewHBox(backBtn),
				nil, nil, nil,
				container.NewVScroll(docPage),
			)

			masterContainer.Objects = []fyne.CanvasObject{detailPage}
			masterContainer.Refresh()
		})

		listContainer.Add(btn)
	}
	listScroll := container.NewVScroll(listContainer)

	listView = container.NewBorder(
		widget.NewLabelWithStyle("Documentation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		listScroll,
	)

	masterContainer.Add(listView)

	return masterContainer

}
func RemoteButton(w fyne.Window, output *widget.Entry) fyne.CanvasObject {
	actions := []string{"list", "remove", "add"}
	actionSelect := widget.NewSelect(actions, func(value string) {})

	runBtn := widget.NewButton("Remote", func() {

		if state.RepoPath == "" {
			dialog.ShowInformation("repository", "invalid path select a repository", w)
			return
		}
		switch actionSelect.Selected {
		case "list":
			result, err := git.GitRemote("list", " ")
			if err != nil {
				output.SetText(fmt.Sprintf("Error: %v\n%s", err, result))
			} else {
				output.SetText(result)
			}

		case "remove":
			// remoteEntry := widget.NewEntry()
			// remoteEntry.SetPlaceHolder("Remote name ")
			// remoteurl := widget.NewEntry()
			// remoteurl.SetPlaceHolder("Remote url")
			// formItems := []*widget.FormItem{
			// 	widget.NewFormItem("", remoteEntry),
			// }

			// dialog.ShowForm("Remove Remote", "Remove", "Cancel", formItems, func(ok bool) {
			// 	if ok {
			// 		result, err := git.GitRemote("remove", "", remoteEntry.Text)
			// 		if err != nil {
			// 			output.SetText(fmt.Sprintf("Error: %v\n%s", err, result))
			// 		} else {
			// 			output.SetText(result)
			// 		}
			// 	}
			// }, w)
			cmd := exec.Command("git", "-C", state.RepoPath, "remote", "remove", "origin")
			out, err := cmd.CombinedOutput()
			msg := string(out)

			if err != nil {
				lower := strings.ToLower(msg)

				// remote does not exist
				if strings.Contains(lower, "no such remote") {
					output.SetText("No previous origin found ")
					// no return needed, gracefully continue
				} else {
					// real error
					output.SetText("Warning: Could not remove old origin:")
					output.SetText(msg)
				}
			} else {
				output.SetText("Old origin removed successfully")
			}

			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		case "add":
			nameEntry := widget.NewEntry()
			nameEntry.SetPlaceHolder("Remote Name")
			urlEntry := widget.NewEntry()
			urlEntry.SetPlaceHolder("Remote URL")

			formItems := []*widget.FormItem{
				widget.NewFormItem("Remote Name", nameEntry),
				widget.NewFormItem("Remote URL", urlEntry),
			}

			dialog.ShowForm("Add Remote", "Add", "Cancel", formItems, func(ok bool) {
				if ok {
					result, err := git.GitRemote("add", urlEntry.Text)
					if err != nil {
						output.SetText(fmt.Sprintf("Error: %v\n%s", err, result))
					} else {
						output.SetText(result)
					}
				}
			}, w)

		default:
			output.SetText("No action selected")
		}
	})

	return container.NewVBox(
		runBtn,
		actionSelect,
	)
}
func DiffButton(output *widget.Entry) fyne.CanvasObject {
	options := []string{"Unstaged", "Staged (--cached)", "Names (--name-only)", "Summary (--stat)"}
	diffSelect := widget.NewSelect(options, func(value string) {})
	diffSelect.SetSelected("Unstaged")

	diffBtn := widget.NewButton("Diff", func() {
		out, err := git.Diff(diffSelect.Selected)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})
	return container.NewVBox(diffBtn, diffSelect)
}
func ResetButton(output *widget.Entry, w fyne.Window) fyne.CanvasObject {
	options := []string{"--mixed", "--soft", "--hard"}
	resetSelect := widget.NewSelect(options, func(value string) {})
	resetSelect.SetSelected("--mixed")

	resetBtn := widget.NewButton("Reset", func() {
		input := widget.NewEntry()
		input.SetText("HEAD~1")
		input.SetPlaceHolder("e.g. HEAD~1 or hash")

		form := []*widget.FormItem{
			{Text: "Commit Hash", Widget: input},
		}

		dialog.ShowForm("Git Reset", "Reset", "Cancel", form, func(valid bool) {
			if !valid {
				return
			}

			target := strings.TrimSpace(input.Text)
			if target == "" {
				dialog.ShowInformation("Invalid Hash", "Please provide a valid commit hash or reference (like HEAD~1).", w)
				return
			}

			mode := resetSelect.Selected
			dialog.ShowConfirm("Confirm Reset",
				fmt.Sprintf("Are you sure you want to proceed with %s reset to %s?", mode, target),
				func(response bool) {
					if response {
						out, err := git.Reset(mode, target)
						if err != nil {
							output.SetText("error: " + err.Error())
						} else {
							output.SetText(out)
						}
					}
				}, w)
		}, w)
	})
	return container.NewVBox(resetBtn, resetSelect)
}

func MagicSyncButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Magic Sync", func() {
		if state.RepoPath == "" {
			output.SetText("Error: No repository selected.")
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			output.SetText("Error: Repository path does not exist.")
			return
		}

		output.SetText("Starting Magic Sync...")
		out, err := git.MagicSync()
		if err != nil {
			output.SetText(out + "\n\nError: " + err.Error())
		} else {
			output.SetText(out + "\n\nMagic Sync completed successfully!")
		}
	})
}

func UndoButton(output *widget.Entry, w fyne.Window) *widget.Button {
	return widget.NewButton("Undo", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		dialog.ShowConfirm("Undo Last Commit", "Are you sure you want to undo the last commit? Changes will remain staged.", func(confirm bool) {
			if confirm {
				out, err := git.UndoLastCommit()
				if err != nil {
					output.SetText("error: " + err.Error())
				} else {
					output.SetText(out)
				}
			}
		}, w)
	})
}

func CherryPickButton(output *widget.Entry, w fyne.Window) *widget.Button {
	return widget.NewButton("Cherry-pick", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		input := widget.NewEntry()
		input.SetPlaceHolder("Enter commit hash")
		form := []*widget.FormItem{
			{Text: "Hash", Widget: input},
		}
		dialog.ShowForm("Cherry-pick Commit", "Apply", "Cancel", form, func(valid bool) {
			if valid {
				hash := strings.TrimSpace(input.Text)
				if hash == "" {
					return
				}
				out, err := git.CherryPick(hash)
				if err != nil {
					output.SetText("error: " + err.Error())
				} else {
					output.SetText("Cherry-pick result:\n" + out)
				}
			}
		}, w)
	})
}

func ConflictButton(output *widget.Entry, w fyne.Window) *widget.Button {
	return widget.NewButton("Conflicts", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		conflicts, err := git.GetConflicts()
		if err != nil {
			output.SetText("error fetching conflicts: " + err.Error())
			return
		}
		if len(conflicts) == 0 {
			dialog.ShowInformation("No Conflicts", "No merge conflicts detected.", w)
			return
		}

		content := container.NewVBox()
		for _, file := range conflicts {
			f := file // capture for closure
			row := container.NewHBox(
				widget.NewLabel(f),
				layout.NewSpacer(),
				widget.NewButton("Keep Mine", func() {
					out, err := git.ResolveConflict(f, "ours")
					if err != nil {
						output.SetText("error: " + err.Error())
					} else {
						output.SetText(out)
						dialog.ShowInformation("Resolved", "Resolved "+f+" with 'ours'", w)
					}
				}),
				widget.NewButton("Take Theirs", func() {
					out, err := git.ResolveConflict(f, "theirs")
					if err != nil {
						output.SetText("error: " + err.Error())
					} else {
						output.SetText(out)
						dialog.ShowInformation("Resolved", "Resolved "+f+" with 'theirs'", w)
					}
				}),
			)
			content.Add(row)
		}

		scroll := container.NewVScroll(content)
		scroll.SetMinSize(fyne.NewSize(400, 300))
		dialog.ShowCustom("Resolve Conflicts", "Close", scroll, w)
	})
}

func TagButton(output *widget.Entry, w fyne.Window) fyne.CanvasObject {
	actions := []string{"List", "Create", "Push", "Delete"}
	tagSelect := widget.NewSelect(actions, func(value string) {})
	tagSelect.SetSelected("List")

	tagBtn := widget.NewButton("Tag", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		action := strings.ToLower(tagSelect.Selected)

		if action == "list" {
			out, err := git.Tag(state.RepoPath, "list", "")
			if err != nil {
				output.SetText("error: " + err.Error())
			} else {
				output.SetText(out)
			}
			return
		}

		// Create, Push, Delete need a tag name
		tagNameEntry := widget.NewEntry()
		tagNameEntry.SetPlaceHolder("Enter tag name")
		form := []*widget.FormItem{
			{Text: "Tag Name", Widget: tagNameEntry},
		}

		dialogTitle := fmt.Sprintf("Git Tag: %s", tagSelect.Selected)
		dialog.ShowForm(dialogTitle, tagSelect.Selected, "Cancel", form, func(valid bool) {
			if valid {
				name := strings.TrimSpace(tagNameEntry.Text)
				if name == "" {
					return
				}
				out, err := git.Tag(state.RepoPath, action, name)
				if err != nil {
					output.SetText("error: " + err.Error())
				} else {
					output.SetText(out)
				}
			}
		}, w)
	})

	return container.NewVBox(tagBtn, tagSelect)
}

func FetchButton(output *widget.Entry, w fyne.Window) fyne.CanvasObject {
	options := []string{"Default", "All (--all)"}
	fetchSelect := widget.NewSelect(options, func(value string) {})
	fetchSelect.SetSelected("Default")

	fetchBtn := widget.NewButton("Fetch", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		out, err := git.Fetch(state.RepoPath, fetchSelect.Selected)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})

	return container.NewVBox(fetchBtn, fetchSelect)
}

func StashButton(output *widget.Entry, w fyne.Window) fyne.CanvasObject {
	options := []string{"Save", "Pop", "List", "Drop", "Apply"}
	stashSelect := widget.NewSelect(options, func(value string) {})
	stashSelect.SetSelected("Save")

	stashBtn := widget.NewButton("Stash", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		action := strings.ToLower(stashSelect.Selected)
		out, err := git.Stash(state.RepoPath, action)
		if err != nil {
			output.SetText("error: " + err.Error())
		} else {
			output.SetText(out)
		}
	})

	return container.NewVBox(stashBtn, stashSelect)
}

func MergeButton(output *widget.Entry, w fyne.Window) fyne.CanvasObject {
	items, getSelected := helpers.BranchSelector(state.RepoPath, w)

	mergeBtn := widget.NewButton("Merge", func() {
		if state.RepoPath == "" {
			dialog.ShowInformation("Repository Not Selected", "Please select a repository first.", w)
			return
		}
		if _, err := os.Stat(state.RepoPath); os.IsNotExist(err) {
			dialog.ShowError(errors.New("Repository path does not exist"), w)
			return
		}

		branch := getSelected()
		if branch == "" {
			dialog.ShowError(errors.New("Please select a branch to merge"), w)
			return
		}

		dialog.ShowConfirm("Confirm Merge", fmt.Sprintf("Are you sure you want to merge '%s' into current branch?", branch), func(confirm bool) {
			if confirm {
				out, err := git.Merge(state.RepoPath, branch)
				if err != nil {
					output.SetText("error: " + err.Error())
				} else {
					output.SetText(out)
				}
			}
		}, w)
	})

	return container.NewVBox(mergeBtn, items)
}
