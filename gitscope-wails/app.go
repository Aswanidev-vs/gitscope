package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/state"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) IsGitAvailable() bool {
	cmd := exec.Command("git", "--version")
	hideWindow(cmd)
	return cmd.Run() == nil
}

func (a *App) GetRepoPath() string {
	return state.RepoPath
}

func (a *App) SelectRepo() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Repository",
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("no directory selected")
	}
	state.RepoPath = dir
	return dir, nil
}

func (a *App) OpenFolder() (string, error) {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Open Folder",
	})
	if err != nil {
		return "", err
	}
	if dir == "" {
		return "", fmt.Errorf("no folder selected")
	}
	state.RepoPath = dir
	return dir, nil
}

func (a *App) Init() (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Init()
}

func (a *App) Status(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Status(option)
}

func (a *App) Stage(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Stage(option)
}

func (a *App) Commit(msg, option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Commit(msg, option)
}

func (a *App) Push(branch string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Push(state.RepoPath, branch)
}

func (a *App) Pull(branch string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Pull(state.RepoPath, branch)
}

func (a *App) Log(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Log(state.RepoPath, option)
}

func (a *App) Revert(hash string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Revert(hash, "--no-edit")
}

func (a *App) Clone(url string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no destination path selected")
	}
	return git.Clone(state.RepoPath, url)
}

func (a *App) CreateBranch(name string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.CreateBranch(state.RepoPath, name)
}

func (a *App) DeleteBranch(name string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.DeleteBranch(state.RepoPath, name)
}

func (a *App) SwitchBranch(name string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.SwitchBranch(state.RepoPath, name)
}

func (a *App) BranchRename(oldName, newName string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.BranchRename(oldName, newName)
}

func (a *App) Diff(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Diff(option)
}

func (a *App) Reset(mode, target string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Reset(mode, target)
}

func (a *App) Fetch(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Fetch(state.RepoPath, option)
}

func (a *App) Stash(action string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Stash(state.RepoPath, action)
}

func (a *App) Merge(branch string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Merge(state.RepoPath, branch)
}

func (a *App) Tag(action, name string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Tag(state.RepoPath, action, name)
}

func (a *App) Remote(action, args string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.GitRemote(action, args)
}

func (a *App) Reflog(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Reflog(state.RepoPath, option)
}

func (a *App) Show(option, target string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Show(state.RepoPath, option, target)
}

func (a *App) LsFiles(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.LsFiles(state.RepoPath, option)
}

func (a *App) Clean(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Clean(state.RepoPath, option)
}

func (a *App) Shortlog(option string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Shortlog(state.RepoPath, option)
}

func (a *App) Blame(file string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Blame(state.RepoPath, file)
}

func (a *App) Worktree(action, args string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Worktree(state.RepoPath, action, args)
}

func (a *App) Rebase(option, target string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.Rebase(state.RepoPath, option, target)
}

func (a *App) CherryPick(hash string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.CherryPick(hash)
}

func (a *App) UndoLastCommit() (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.UndoLastCommit()
}

func (a *App) MagicSync() (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.MagicSync()
}

func (a *App) GetConflicts() ([]string, error) {
	if state.RepoPath == "" {
		return nil, fmt.Errorf("no repository selected")
	}
	return git.GetConflicts()
}

func (a *App) ResolveConflict(file, strategy string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	return git.ResolveConflict(file, strategy)
}

func (a *App) GetBranches() ([]string, error) {
	if state.RepoPath == "" {
		return nil, fmt.Errorf("no repository selected")
	}
	cmd := exec.Command("git", "-C", state.RepoPath, "branch", "--list")
	hideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		branches = append(branches, strings.TrimPrefix(line, "* "))
	}
	return branches, nil
}

func (a *App) GetCurrentBranch() (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	cmd := exec.Command("git", "-C", state.RepoPath, "branch", "--show-current")
	hideWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (a *App) ReadGitIgnore() (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	path := filepath.Join(state.RepoPath, ".gitignore")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}

func (a *App) WriteGitIgnore(content string) error {
	if state.RepoPath == "" {
		return fmt.Errorf("no repository selected")
	}
	path := filepath.Join(state.RepoPath, ".gitignore")
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *App) RunCommands(cmdText string) (string, error) {
	if state.RepoPath == "" {
		return "", fmt.Errorf("no repository selected")
	}
	lines := strings.Split(cmdText, "\n")
	var log strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var cmd *exec.Cmd
		if strings.HasPrefix(line, "git ") {
			args := strings.Fields(line)
			cmd = exec.Command("git", args[1:]...)
		} else {
			cmd = exec.Command("cmd", "/C", line)
		}
		cmd.Dir = state.RepoPath
		hideWindow(cmd)
		out, err := cmd.CombinedOutput()
		log.WriteString(fmt.Sprintf("> %s\n%s\n", line, string(out)))
		if err != nil {
			log.WriteString(fmt.Sprintf("Error: %v\n", err))
		}
	}
	return log.String(), nil
}

func (a *App) IsRepoInitialized() bool {
	if state.RepoPath == "" {
		return false
	}
	gitDir := filepath.Join(state.RepoPath, ".git")
	info, err := os.Stat(gitDir)
	return err == nil && info.IsDir()
}
