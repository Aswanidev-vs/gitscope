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
	return container.NewWithoutLayout(initBtn, stageBtn, commitBtn, statusBtn, pushBtn, logBtn, revertBtn, cloneBtn, Branchbtn, PullBtn, clearBtn, Reflogbtn, SwitchBranchBtn, BranchRenameBtn, GitIgnoreBtn, GitRemotebtn, output)
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

					sha, err := helpers.GetPreviousCommit(state.RepoPath)
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

	f1 := widget.NewLabel("GitScope is a modern, lightweight, and visually intuitive Git client built with Go and Fyne. It simplifies essential")
	f2 := widget.NewLabel("version control operations making Git easier to use for both beginners and experienced developers.")
	f3 := widget.NewLabel("Version: 1.0.0")
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
		out, err := git.Reflog(state.RepoPath)
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

	items := []string{"Init", "Stage", "Status", "Commit", "Push", "Log", "Revert", "Clone", "Branch", "Pull", "Reflog", "GitIgnore", "Remote"}

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
			urlEntry := widget.NewEntry()
			urlEntry.SetPlaceHolder("Remote URL")

			formItems := []*widget.FormItem{

				widget.NewFormItem("", urlEntry),
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
