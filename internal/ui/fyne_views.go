package ui

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	core "github.com/gitscope/internal/core"
	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/helpers"
	"github.com/gitscope/internal/state"
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
	NewRepo.Resize(fyne.NewSize(400, 500))
	NewRepo.Move(fyne.NewPos(1, 500))

	return container.NewVBox(
		widget.NewLabel("Repository Setup"),
		browseBtn,
		output,
		NewRepo,
	)
}

func dashBoardPage(w fyne.Window) fyne.CanvasObject {
	output := widget.NewMultiLineEntry()

	output.Resize(fyne.NewSize(500, 230))
	output.Move(fyne.NewPos(0, 1))

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

	revertBtn := RevertButton(output)
	revertBtn.Resize(fyne.NewSize(100, 40))
	revertBtn.Move(fyne.NewPos(110, 350))

	return container.NewWithoutLayout(initBtn, stageBtn, commitBtn, statusBtn, pushBtn, logBtn, revertBtn, output)
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

	branchSelectorUI, getBranch := helpers.BranchSelector(state.RepoPath)
	pushBtn := widget.NewButton("Push", func() {
		if state.RepoPath == "" {
			dialog.ShowError(errors.New("No repository selected"), w)
			return
		}
		progress := dialog.NewProgressInfinite("Running Commands", "Please wait while commands are executing...", w)

		go func() {
			progress.Show()

			branch := getBranch()
			if branch == "" {
				dialog.ShowError(errors.New("No branch selected"), w)
				return
			}

			output, err := git.Push(state.RepoPath, branch)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Push failed:\n%v\n\n%s", err, output), w)
				return
			}
			progress.Hide()
			dialog.ShowInformation("Push Success", "Repository pushed successfully.", w)
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
func RevertButton(output *widget.Entry) *widget.Button {
	return widget.NewButton("Revert", func() {
		out, err := git.Revert(state.RepoPath)
		if err != nil {
			output.SetText("error:" + err.Error())
		} else {
			output.SetText(out)
		}
	})
}
