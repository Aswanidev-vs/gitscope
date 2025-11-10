package main

import (
	"fmt"
	"os/exec"
	"syscall"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/gitscope/internal/ui"
)

func main() {
	// Check if Git is available before starting the app
	if !isGitAvailable() {
		// Create a temporary app just to show the error dialog
		tempApp := app.New()
		tempWindow := tempApp.NewWindow("GitScope - Git Not Found")
		tempWindow.SetContent(nil) // Empty content

		dialog.ShowError(fmt.Errorf("Git is not installed or not available in PATH.\n\nPlease install Git and ensure it's added to your system PATH.\n\nDownload from: https://git-scm.com/downloads"), tempWindow)
		tempWindow.ShowAndRun()
		return
	}

	ui.App()
}

func isGitAvailable() bool {
	cmd := exec.Command("git", "--version")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	return err == nil
}
