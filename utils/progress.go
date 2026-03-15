package utils

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowCustomProgress creates and shows a custom non-blocking progress dialog.
// It returns the dialog so the caller can Hide() it when finished.
func ShowCustomProgress(title, message string, w fyne.Window) dialog.Dialog {
	prog := widget.NewProgressBarInfinite()
	label := widget.NewLabel(message)
	label.Alignment = fyne.TextAlignCenter
	content := container.NewPadded(container.NewVBox(label, prog))
	dlg := dialog.NewCustomWithoutButtons(title, content, w)
	dlg.Show()
	return dlg
}
