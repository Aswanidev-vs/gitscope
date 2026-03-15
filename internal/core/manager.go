package core

import (
	"errors"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gitscope/internal/helpers"
	"github.com/gitscope/internal/state"
)

func OpenRepo(w fyne.Window, onResult func(string)) {
	dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil || uri == nil {
			onResult("error: unable to open folder")
			return
		}

		state.RepoPath = uri.Path() // Assign to global variable (no :=)
		onResult("Repository set to: " + state.RepoPath)
	}, w).Show()
}
func CreateNewRepo(w fyne.Window, onSubmit func(string)) fyne.CanvasObject {
	label := widget.NewLabel("Create New Repository")
	multiline := widget.NewMultiLineEntry()
	multiline.SetPlaceHolder("Paste the GitHub commands to create a new repository")
	multiline.SetMinRowsVisible(5)

	submit := widget.NewButton("Run Commands", func() {
		text := multiline.Text
		fmt.Println("Submitted text:", text)
		onSubmit(text)
		helpers.NewRepoCmd(w, state.RepoPath, text)
	})

	return container.NewVBox(
		label,
		multiline,
		submit,
	)
}

func ExistingRepo(w fyne.Window, onSubmit func(string)) fyne.CanvasObject {
	label := widget.NewLabel("Push Existing Repository")
	multiline := widget.NewMultiLineEntry()
	multiline.SetPlaceHolder("Paste git commands here...\nExample:\ngit remote add origin https://github.com/yourname/t2.git\ngit branch -M main\ngit push -u origin main")
	multiline.SetMinRowsVisible(5)

	submit := widget.NewButton("Run Commands", func() {
		cmdText := strings.TrimSpace(multiline.Text)
		if cmdText == "" {
			dialog.ShowError(errors.New("No commands entered"), w)
			return
		}
		onSubmit(cmdText)
	})

	return container.NewVBox(
		label,
		multiline,
		submit,
	)
}
