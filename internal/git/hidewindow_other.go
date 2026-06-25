//go:build !windows

package git

import "os/exec"

func hideWindow(cmd *exec.Cmd) {}
