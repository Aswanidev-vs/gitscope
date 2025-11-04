package utils

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gitscope/internal/state"
)

func NewBranchButton(win fyne.Window, onCreate, onDelete func(string)) *widget.Button {
	btn := widget.NewButton("Branch", func() {
		repoPath := state.RepoPath // ✅ directly use the global state

		// --- Check if repository is valid ---
		if repoPath == "" {
			dialog.ShowInformation("No Repository", "No repository path set. Please initialize or open a repository first.", win)
			return
		}

		info, err := os.Stat(repoPath)
		if err != nil || !info.IsDir() {
			dialog.ShowInformation("Invalid Repository", "The specified repository path is invalid or missing.", win)
			return
		}

		if _, err := os.Stat(repoPath + "/.git"); os.IsNotExist(err) {
			dialog.ShowInformation("Not a Git Repository", "This folder is not a Git repository. Please run Init first.", win)
			return
		}

		// --- Show dialog for branch creation/deletion ---
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter branch name")

		createBtn := widget.NewButton("Create", func() {
			name := entry.Text
			if name == "" {
				dialog.ShowInformation("Invalid Input", "Please enter a branch name.", win)
				return
			}
			if onCreate != nil {
				onCreate(name)
			}
			dialog.ShowInformation("Branch Created", "Branch '"+name+"' created successfully.", win)
		})

		deleteBtn := widget.NewButton("Delete", func() {
			name := entry.Text
			if name == "" {
				dialog.ShowInformation("Invalid Input", "Please enter a branch name.", win)
				return
			}
			if onDelete != nil {
				onDelete(name)
			}
			dialog.ShowInformation("Branch Deleted", "Branch '"+name+"' deleted successfully.", win)
		})

		buttons := container.NewHBox(
			createBtn,
			layout.NewSpacer(),
			deleteBtn,
		)

		content := container.NewVBox(
			entry,
			layout.NewSpacer(),
			buttons,
		)

		d := dialog.NewCustom("Branch Actions", "Close", content, win)
		d.Resize(fyne.NewSize(320, 160))
		d.Show()
	})

	return btn
}
