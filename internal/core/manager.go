package core

import (
	"fmt"

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
	multiline := widget.NewMultiLineEntry()
	multiline.SetPlaceHolder("Paste the GitHub commands to create a new repository")
	multiline.Resize(fyne.NewSize(500, 400))

	submit := widget.NewButton("Submit", func() {
		text := multiline.Text
		fmt.Println("Submitted text:", text)
		onSubmit(text)
		// helpers.NewRepoCmd(w, state.RepoPath, text)
		helpers.NewRepoCmd(w, state.RepoPath, text)
	})

	content := container.NewVBox(
		widget.NewLabel("Create New Repository"),
		multiline,
		submit,
	)
	content.Resize(fyne.NewSize(500, 400))
	return content
}
