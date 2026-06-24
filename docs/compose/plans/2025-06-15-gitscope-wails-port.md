# GitScope Fyne-to-Wails Port Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port GitScope desktop GUI from Go+Fyne to Go+Wails with plain HTML/JS/CSS frontend, plus GitHub Actions CI/CD for cross-platform releases.

**Architecture:** Backend reuses existing `internal/git/` and `internal/state/` packages with Fyne dependencies removed. New Wails `app.go` exposes all git operations as methods callable from frontend. Frontend is plain HTML/JS/CSS with violet theme. CI/CD triggers on `v*` tags to build and release.

**Tech Stack:** Go 1.23+, Wails v2, HTML/CSS/JS, GitHub Actions

---

## File Structure

```
gitscope-wails/
├── main.go                    # Wails app entry
├── app.go                     # App struct with all git operations
├── go.mod                     # Dependencies
├── wails.json                 # Wails config
├── frontend/
│   ├── index.html             # Main HTML
│   ├── style.css              # Styling
│   ├── app.js                 # UI logic + Wails bindings
│   └── wailsjs/               # Auto-generated
├── internal/
│   ├── git/
│   │   └── git_ops.go         # Git operations (no Fyne deps)
│   └── state/
│       └── state.go           # Global RepoPath
.github/
└── workflows/
    └── release.yml            # CI/CD workflow
```

---

### Task 1: Create Wails App Backend

**Covers:** [S2, S4]

**Files:**
- Create: `gitscope-wails/main.go`
- Create: `gitscope-wails/app.go`

- [ ] **Step 1: Create main.go with Wails setup**

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "GitScope",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

- [ ] **Step 2: Create app.go with all git operation methods**

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gitscope/internal/git"
	"github.com/gitscope/internal/state"
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

func (a *App) OpenRepo(path string) error {
	state.RepoPath = path
	return nil
}

func (a *App) GetRepoPath() string {
	return state.RepoPath
}

func (a *App) Init() (string, error) {
	return git.Init()
}

func (a *App) Status(option string) (string, error) {
	return git.Status(option)
}

func (a *App) Stage(option string) (string, error) {
	return git.Stage(option)
}

func (a *App) Commit(msg, option string) (string, error) {
	return git.Commit(msg, option)
}

func (a *App) Push(branch string) (string, error) {
	return git.Push(state.RepoPath, branch)
}

func (a *App) Pull(branch string) (string, error) {
	return git.Pull(state.RepoPath, branch)
}

func (a *App) Log(option string) (string, error) {
	return git.Log(state.RepoPath, option)
}

func (a *App) Diff(option string) (string, error) {
	return git.Diff(option)
}

func (a *App) Reset(mode, target string) (string, error) {
	return git.Reset(mode, target)
}

func (a *App) Fetch(option string) (string, error) {
	return git.Fetch(state.RepoPath, option)
}

func (a *App) Stash(action string) (string, error) {
	return git.Stash(state.RepoPath, action)
}

func (a *App) Merge(branchname string) (string, error) {
	return git.Merge(state.RepoPath, branchname)
}

func (a *App) Tag(action, tagname string) (string, error) {
	return git.Tag(state.RepoPath, action, tagname)
}

func (a *App) BranchCreate(name string) (string, error) {
	return git.CreateBranch(state.RepoPath, name)
}

func (a *App) BranchDelete(name string) (string, error) {
	return git.DeleteBranch(state.RepoPath, name)
}

func (a *App) BranchRename(oldname, newname string) (string, error) {
	return git.BranchRename(oldname, newname)
}

func (a *App) SwitchBranch(branchname string) (string, error) {
	return git.SwitchBranch(state.RepoPath, branchname)
}

func (a *App) Revert(commitHash string) (string, error) {
	return git.Revert(commitHash, "--no-edit")
}

func (a *App) Remote(action, args string) (string, error) {
	return git.GitRemote(action, args)
}

func (a *App) Show(option, target string) (string, error) {
	return git.Show(state.RepoPath, option, target)
}

func (a *App) LsFiles(option string) (string, error) {
	return git.LsFiles(state.RepoPath, option)
}

func (a *App) Blame(file string) (string, error) {
	return git.Blame(state.RepoPath, file)
}

func (a *App) Worktree(action, args string) (string, error) {
	return git.Worktree(state.RepoPath, action, args)
}

func (a *App) Shortlog(option string) (string, error) {
	return git.Shortlog(state.RepoPath, option)
}

func (a *App) Reflog(option string) (string, error) {
	return git.Reflog(state.RepoPath, option)
}

func (a *App) Rebase(option, target string) (string, error) {
	return git.Rebase(state.RepoPath, option, target)
}

func (a *App) Clean(option string) (string, error) {
	return git.Clean(state.RepoPath, option)
}

func (a *App) CherryPick(hash string) (string, error) {
	return git.CherryPick(hash)
}

func (a *App) MagicSync() (string, error) {
	return git.MagicSync()
}

func (a *App) UndoLastCommit() (string, error) {
	return git.UndoLastCommit()
}

func (a *App) GetConflicts() ([]string, error) {
	return git.GetConflicts()
}

func (a *App) ResolveConflict(file, strategy string) (string, error) {
	return git.ResolveConflict(file, strategy)
}

func (a *App) ListBranches() ([]string, error) {
	repo := state.RepoPath
	cmd := exec.Command("git", "-C", repo, "branch", "--list")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if after, ok := strings.CutPrefix(line, "* "); ok {
			line = after
		}
		branches = append(branches, line)
	}
	return branches, nil
}

func (a *App) NewRepoCmd(cmdText string) (string, error) {
	repo := state.RepoPath
	if repo == "" {
		return "", fmt.Errorf("no repository path selected")
	}
	lines := strings.Split(cmdText, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cmd := exec.Command("sh", "-c", line)
		cmd.Dir = repo
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if _, err := cmd.CombinedOutput(); err != nil {
			if !strings.Contains(line, "already exists") {
				return "", fmt.Errorf("command failed: %s", line)
			}
		}
	}
	return "All commands executed successfully", nil
}

func (a *App) ExistingRepoCmd(cmdText string) (string, error) {
	repo := state.RepoPath
	if repo == "" {
		return "", fmt.Errorf("no repository path selected")
	}
	lines := strings.Split(cmdText, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cmd := exec.Command("sh", "-c", line)
		cmd.Dir = repo
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if _, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("command failed: %s", line)
		}
	}
	return "All commands executed successfully", nil
}

func (a *App) GitIgnore(action string) (string, error) {
	repo := state.RepoPath
	filePath := filepath.Join(repo, ".gitignore")
	if action == "read" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", nil
			}
			return "", err
		}
		return string(content), nil
	}
	return "", nil
}

func (a *App) SaveGitIgnore(content string) error {
	repo := state.RepoPath
	filePath := filepath.Join(repo, ".gitignore")
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (a *App) Clone(cloneURL string) (string, error) {
	repo := state.RepoPath
	if repo == "" {
		return "", fmt.Errorf("no repository path selected")
	}
	cmd := exec.Command("git", "-C", repo, "clone", cloneURL)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("clone failed: %v", err)
	}
	return "Successfully cloned the repo", nil
}

func (a *App) GetPreviousCommit() (string, error) {
	return git.GetPreviousCommit()
}
```

- [ ] **Step 3: Verify Go files compile**

Run: `cd gitscope-wails && go build ./...`
Expected: No errors (may need `go mod tidy` first)

- [ ] **Step 4: Commit**

```bash
git add gitscope-wails/main.go gitscope-wails/app.go
git commit -m "feat: add Wails backend with all git operation bindings"
```

---

### Task 2: Create Frontend HTML Structure

**Covers:** [S3, S5]

**Files:**
- Create: `gitscope-wails/frontend/index.html`

- [ ] **Step 1: Create index.html**

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitScope</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div class="app">
        <nav class="sidebar">
            <button class="nav-btn active" data-page="repository" title="Repository">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
                </svg>
            </button>
            <button class="nav-btn" data-page="dashboard" title="Git Operations">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                </svg>
            </button>
            <button class="nav-btn" data-page="settings" title="Settings">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="3"></circle>
                    <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
                </svg>
            </button>
            <button class="nav-btn" data-page="docs" title="Documentation">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <polyline points="14 2 14 8 20 8"></polyline>
                    <line x1="16" y1="13" x2="8" y2="13"></line>
                    <line x1="16" y1="17" x2="8" y2="17"></line>
                </svg>
            </button>
            <div class="spacer"></div>
        </nav>

        <main class="content">
            <!-- Repository Page -->
            <div id="page-repository" class="page active">
                <h2>Repository Setup</h2>
                <button id="btn-select-repo" class="btn">Select Repository</button>
                <div id="repo-path" class="repo-path">No repository selected</div>
                
                <div class="section">
                    <h3>Create New Repository</h3>
                    <textarea id="new-repo-cmds" placeholder="Paste GitHub commands here..."></textarea>
                    <button id="btn-run-new-repo" class="btn">Run Commands</button>
                </div>
                
                <div class="section">
                    <h3>Push Existing Repository</h3>
                    <textarea id="existing-repo-cmds" placeholder="Paste git commands here..."></textarea>
                    <button id="btn-run-existing-repo" class="btn">Run Commands</button>
                </div>
            </div>

            <!-- Dashboard Page -->
            <div id="page-dashboard" class="page">
                <div class="dashboard">
                    <div class="categories">
                        <h3>Categories</h3>
                        <button class="cat-btn active" data-cat="common">Common</button>
                        <button class="cat-btn" data-cat="branches">Branches</button>
                        <button class="cat-btn" data-cat="remote">Remote</button>
                        <button class="cat-btn" data-cat="history">History</button>
                        <button class="cat-btn" data-cat="changes">Changes</button>
                        <button class="cat-btn" data-cat="advanced">Advanced</button>
                        <button class="cat-btn" data-cat="tools">Tools</button>
                    </div>
                    
                    <div class="main-area">
                        <div class="top-bar">
                            <h2>GitScope Dashboard</h2>
                            <button id="btn-clear-console" class="btn btn-small">Clear Console</button>
                        </div>
                        
                        <div class="grid" id="cat-common">
                            <button class="grid-btn" data-op="init">Init</button>
                            <button class="grid-btn" data-op="stage">Stage</button>
                            <button class="grid-btn" data-op="status">Status</button>
                            <button class="grid-btn" data-op="commit">Commit</button>
                            <button class="grid-btn" data-op="push">Push</button>
                            <button class="grid-btn" data-op="log">Log</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-branches">
                            <button class="grid-btn" data-op="branch-create">Branch</button>
                            <button class="grid-btn" data-op="switch-branch">Switch Branch</button>
                            <button class="grid-btn" data-op="merge">Merge</button>
                            <button class="grid-btn" data-op="branch-rename">Rename</button>
                            <button class="grid-btn" data-op="tag">Tag</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-remote">
                            <button class="grid-btn" data-op="remote">Remote</button>
                            <button class="grid-btn" data-op="fetch">Fetch</button>
                            <button class="grid-btn" data-op="pull">Pull</button>
                            <button class="grid-btn" data-op="clone">Clone</button>
                            <button class="grid-btn" data-op="cherry-pick">Cherry-pick</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-history">
                            <button class="grid-btn" data-op="log">Log</button>
                            <button class="grid-btn" data-op="revert">Revert</button>
                            <button class="grid-btn" data-op="show">Show</button>
                            <button class="grid-btn" data-op="shortlog">Shortlog</button>
                            <button class="grid-btn" data-op="reflog">Reflog</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-changes">
                            <button class="grid-btn" data-op="diff">Diff</button>
                            <button class="grid-btn" data-op="stash">Stash</button>
                            <button class="grid-btn" data-op="clean">Clean</button>
                            <button class="grid-btn" data-op="ls-files">Ls-Files</button>
                            <button class="grid-btn" data-op="gitignore">GitIgnore</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-advanced">
                            <button class="grid-btn" data-op="reset">Reset</button>
                            <button class="grid-btn" data-op="rebase">Rebase</button>
                            <button class="grid-btn" data-op="undo">Undo</button>
                            <button class="grid-btn" data-op="worktree">Worktree</button>
                            <button class="grid-btn" data-op="conflicts">Conflicts</button>
                        </div>
                        
                        <div class="grid hidden" id="cat-tools">
                            <button class="grid-btn" data-op="blame">Blame</button>
                            <button class="grid-btn" data-op="magic-sync">Magic Sync</button>
                        </div>
                    </div>
                </div>
                
                <div class="console">
                    <div class="console-header">
                        <span>Console Output</span>
                    </div>
                    <pre id="console-output"></pre>
                </div>
            </div>

            <!-- Settings Page -->
            <div id="page-settings" class="page">
                <div class="settings-content">
                    <h2>About GitScope</h2>
                    <p>GitScope is a modern, lightweight, and visually intuitive Git client built with Go and Wails. It simplifies essential version control operations making Git easier to use for both beginners and experienced developers.</p>
                    <p>Version: 1.0.0</p>
                    <p>Developer: Aswanidev VS</p>
                    <a href="https://github.com/Aswanidev-vs/GitScope" target="_blank">View Project on GitHub</a>
                </div>
            </div>

            <!-- Documentation Page -->
            <div id="page-docs" class="page">
                <h2>Documentation</h2>
                <div class="doc-list" id="doc-list"></div>
                <div id="doc-content" class="doc-content hidden"></div>
            </div>
        </main>
    </div>

    <!-- Modal for operations -->
    <div id="modal" class="modal hidden">
        <div class="modal-content">
            <h3 id="modal-title"></h3>
            <div id="modal-body"></div>
            <div class="modal-actions">
                <button id="modal-cancel" class="btn">Cancel</button>
                <button id="modal-confirm" class="btn btn-primary">Confirm</button>
            </div>
        </div>
    </div>

    <script src="app.js"></script>
</body>
</html>
```

- [ ] **Step 2: Commit**

```bash
git add gitscope-wails/frontend/index.html
git commit -m "feat: add frontend HTML structure"
```

---

### Task 3: Create Frontend CSS Styling

**Covers:** [S5]

**Files:**
- Create: `gitscope-wails/frontend/style.css`

- [ ] **Step 1: Create style.css**

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #1a1b26;
    color: #c0caf5;
    height: 100vh;
    overflow: hidden;
}

.app {
    display: flex;
    height: 100vh;
}

.sidebar {
    width: 50px;
    background: #16161e;
    display: flex;
    flex-direction: column;
    padding: 10px 0;
    border-right: 1px solid #2f3347;
}

.nav-btn {
    width: 40px;
    height: 40px;
    margin: 5px auto;
    border: none;
    background: transparent;
    color: #565f89;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.nav-btn:hover {
    background: #2f3347;
    color: #c0caf5;
}

.nav-btn.active {
    background: #9400d3;
    color: white;
}

.spacer {
    flex: 1;
}

.content {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
}

.page {
    display: none;
}

.page.active {
    display: block;
}

h2 {
    color: #bb9af7;
    margin-bottom: 20px;
}

h3 {
    color: #7aa2f7;
    margin-bottom: 10px;
}

.btn {
    padding: 10px 20px;
    border: none;
    background: #9400d3;
    color: white;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    transition: background 0.2s;
}

.btn:hover {
    background: #7b00b3;
}

.btn-small {
    padding: 6px 12px;
    font-size: 12px;
}

.btn-primary {
    background: #9ece6a;
}

.btn-primary:hover {
    background: #7aa2f7;
}

.repo-path {
    padding: 10px;
    background: #16161e;
    border-radius: 6px;
    margin: 10px 0;
    font-family: monospace;
    color: #9ece6a;
}

.section {
    margin-top: 20px;
    padding: 15px;
    background: #16161e;
    border-radius: 8px;
}

textarea {
    width: 100%;
    height: 100px;
    padding: 10px;
    background: #1a1b26;
    border: 1px solid #2f3347;
    border-radius: 6px;
    color: #c0caf5;
    font-family: monospace;
    font-size: 13px;
    resize: vertical;
    margin-bottom: 10px;
}

textarea:focus {
    outline: none;
    border-color: #9400d3;
}

.dashboard {
    display: flex;
    gap: 20px;
}

.categories {
    width: 120px;
}

.cat-btn {
    display: block;
    width: 100%;
    padding: 10px;
    margin: 5px 0;
    border: none;
    background: transparent;
    color: #c0caf5;
    text-align: left;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
}

.cat-btn:hover {
    background: #2f3347;
}

.cat-btn.active {
    background: #9400d3;
    color: white;
}

.main-area {
    flex: 1;
}

.top-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
}

.grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 10px;
}

.grid.hidden {
    display: none;
}

.grid-btn {
    padding: 15px 10px;
    border: 1px solid #2f3347;
    background: #16161e;
    color: #c0caf5;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 13px;
}

.grid-btn:hover {
    background: #2f3347;
    border-color: #9400d3;
}

.console {
    margin-top: 20px;
    background: #16161e;
    border-radius: 8px;
    overflow: hidden;
}

.console-header {
    padding: 10px 15px;
    background: #1a1b26;
    border-bottom: 1px solid #2f3347;
    font-weight: bold;
    color: #7aa2f7;
}

#console-output {
    padding: 15px;
    max-height: 200px;
    overflow-y: auto;
    font-family: monospace;
    font-size: 13px;
    white-space: pre-wrap;
    color: #9ece6a;
}

.settings-content {
    text-align: center;
    padding: 40px;
}

.settings-content p {
    margin: 10px 0;
    color: #a9b1d6;
}

.settings-content a {
    color: #9400d3;
    text-decoration: none;
}

.settings-content a:hover {
    text-decoration: underline;
}

.doc-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 10px;
}

.doc-btn {
    padding: 15px;
    border: 1px solid #2f3347;
    background: #16161e;
    color: #c0caf5;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
}

.doc-btn:hover {
    background: #2f3347;
    border-color: #9400d3;
}

.doc-content {
    padding: 20px;
    background: #16161e;
    border-radius: 8px;
    margin-top: 15px;
}

.doc-content.hidden {
    display: none;
}

.modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal.hidden {
    display: none;
}

.modal-content {
    background: #1a1b26;
    padding: 20px;
    border-radius: 12px;
    min-width: 400px;
    max-width: 500px;
}

.modal-content h3 {
    margin-bottom: 15px;
}

.modal-body input,
.modal-body textarea {
    width: 100%;
    margin-bottom: 10px;
}

.modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 15px;
}

select {
    padding: 10px;
    background: #16161e;
    border: 1px solid #2f3347;
    border-radius: 6px;
    color: #c0caf5;
    margin-bottom: 10px;
    width: 100%;
}

select:focus {
    outline: none;
    border-color: #9400d3;
}
```

- [ ] **Step 2: Commit**

```bash
git add gitscope-wails/frontend/style.css
git commit -m "feat: add frontend CSS styling with violet theme"
```

---

### Task 4: Create Frontend JavaScript Logic

**Covers:** [S5]

**Files:**
- Create: `gitscope-wails/frontend/app.js`

- [ ] **Step 1: Create app.js**

```javascript
// Document content for help pages
const docContent = {
    'Init': 'Initialize a new Git repository.\n\nCommand: git init',
    'Stage': 'Stage changes for commit.\n\nCommand: git add .',
    'Status': 'Show working tree status.\n\nCommand: git status',
    'Commit': 'Record changes to the repository.\n\nCommand: git commit -m "message"',
    'Push': 'Update remote references.\n\nCommand: git push origin branch',
    'Log': 'Show commit logs.\n\nCommand: git log --oneline',
    'Revert': 'Revert a commit.\n\nCommand: git revert <commit-hash>',
    'Clone': 'Clone a repository.\n\nCommand: git clone <url>',
    'Branch': 'Create/delete branches.\n\nCommand: git branch <name>',
    'Pull': 'Fetch and merge changes.\n\nCommand: git pull origin branch',
    'Reflog': 'Show reference logs.\n\nCommand: git reflog',
    'GitIgnore': 'Edit .gitignore file.',
    'Remote': 'Manage remote repositories.\n\nCommand: git remote -v',
    'Diff': 'Show changes.\n\nCommand: git diff',
    'Reset': 'Reset current HEAD.\n\nCommand: git reset --mixed HEAD~1',
    'Fetch': 'Download objects and refs.\n\nCommand: git fetch',
    'Stash': 'Stash changes.\n\nCommand: git stash save',
    'Merge': 'Merge branches.\n\nCommand: git merge <branch>',
    'Tag': 'Create/list/delete tags.\n\nCommand: git tag <name>',
    'Cherry-pick': 'Apply a commit.\n\nCommand: git cherry-pick <hash>',
    'Rebase': 'Reapply commits.\n\nCommand: git rebase <branch>',
    'Clean': 'Remove untracked files.\n\nCommand: git clean -n',
    'Show': 'Show commit details.\n\nCommand: git show <hash>',
    'Ls-files': 'List files in index.\n\nCommand: git ls-files',
    'Worktree': 'Manage working trees.\n\nCommand: git worktree list',
    'Shortlog': 'Summary of git log.\n\nCommand: git shortlog -s',
    'Blame': 'Show who modified each line.\n\nCommand: git blame <file>'
};

// State
let currentPage = 'repository';
let currentCategory = 'common';

// DOM Elements
const pages = document.querySelectorAll('.page');
const navBtns = document.querySelectorAll('.nav-btn');
const catBtns = document.querySelectorAll('.cat-btn');
const grids = document.querySelectorAll('.grid');
const consoleOutput = document.getElementById('console-output');
const modal = document.getElementById('modal');
const modalTitle = document.getElementById('modal-title');
const modalBody = document.getElementById('modal-body');
const modalCancel = document.getElementById('modal-cancel');
const modalConfirm = document.getElementById('modal-confirm');

// Navigation
navBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const page = btn.dataset.page;
        switchPage(page);
    });
});

function switchPage(page) {
    currentPage = page;
    pages.forEach(p => p.classList.remove('active'));
    navBtns.forEach(b => b.classList.remove('active'));
    document.getElementById(`page-${page}`).classList.add('active');
    document.querySelector(`[data-page="${page}"]`).classList.add('active');
}

// Category switching
catBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const cat = btn.dataset.cat;
        currentCategory = cat;
        catBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        grids.forEach(g => g.classList.add('hidden'));
        document.getElementById(`cat-${cat}`).classList.remove('hidden');
    });
});

// Clear console
document.getElementById('btn-clear-console').addEventListener('click', () => {
    consoleOutput.textContent = '';
});

// Repository page
document.getElementById('btn-select-repo').addEventListener('click', async () => {
    const path = await prompt('Enter repository path:');
    if (path) {
        await window.go.main.App.OpenRepo(path);
        document.getElementById('repo-path').textContent = path;
    }
});

document.getElementById('btn-run-new-repo').addEventListener('click', async () => {
    const cmds = document.getElementById('new-repo-cmds').value;
    if (!cmds.trim()) return;
    try {
        const result = await window.go.main.App.NewRepoCmd(cmds);
        consoleOutput.textContent = result;
    } catch (err) {
        consoleOutput.textContent = 'Error: ' + err;
    }
});

document.getElementById('btn-run-existing-repo').addEventListener('click', async () => {
    const cmds = document.getElementById('existing-repo-cmds').value;
    if (!cmds.trim()) return;
    try {
        const result = await window.go.main.App.ExistingRepoCmd(cmds);
        consoleOutput.textContent = result;
    } catch (err) {
        consoleOutput.textContent = 'Error: ' + err;
    }
});

// Grid buttons - operations
document.querySelectorAll('.grid-btn').forEach(btn => {
    btn.addEventListener('click', async () => {
        const op = btn.dataset.op;
        await handleOperation(op);
    });
});

async function handleOperation(op) {
    switch (op) {
        case 'init':
            await simpleOp('Init', () => window.go.main.App.Init());
            break;
            
        case 'stage':
            await simpleOp('Stage', () => window.go.main.App.Stage('All (.)'));
            break;
            
        case 'status':
            await showOptions('Status', ['Standard', 'Short (-s)', 'Branch (-b)'], 
                (opt) => window.go.main.App.Status(opt));
            break;
            
        case 'commit':
            await showCommitDialog();
            break;
            
        case 'push':
            await showBranchDialog('Push', async (branch) => {
                return await window.go.main.App.Push(branch);
            });
            break;
            
        case 'log':
            await showOptions('Log', ['Oneline', 'Graph', 'Pretty'],
                (opt) => window.go.main.App.Log(opt));
            break;
            
        case 'branch-create':
            await showInputDialog('Create Branch', 'Branch name:', async (name) => {
                return await window.go.main.App.BranchCreate(name);
            });
            break;
            
        case 'switch-branch':
            await showBranchDialog('Switch Branch', async (branch) => {
                return await window.go.main.App.SwitchBranch(branch);
            });
            break;
            
        case 'merge':
            await showBranchDialog('Merge', async (branch) => {
                return await window.go.main.App.Merge(branch);
            });
            break;
            
        case 'branch-rename':
            await showRenameDialog();
            break;
            
        case 'tag':
            await showTagDialog();
            break;
            
        case 'remote':
            await showRemoteDialog();
            break;
            
        case 'fetch':
            await showOptions('Fetch', ['Default', 'All (--all)'],
                (opt) => window.go.main.App.Fetch(opt));
            break;
            
        case 'pull':
            await showBranchDialog('Pull', async (branch) => {
                return await window.go.main.App.Pull(branch);
            });
            break;
            
        case 'clone':
            await showInputDialog('Clone Repository', 'Repository URL:', async (url) => {
                return await window.go.main.App.Clone(url);
            });
            break;
            
        case 'cherry-pick':
            await showInputDialog('Cherry-pick', 'Commit hash:', async (hash) => {
                return await window.go.main.App.CherryPick(hash);
            });
            break;
            
        case 'revert':
            await showInputDialog('Revert', 'Commit hash:', async (hash) => {
                return await window.go.main.App.Revert(hash);
            });
            break;
            
        case 'show':
            await showOptions('Show', ['Head', 'Last 5'],
                (opt) => window.go.main.App.Show(opt, ''));
            break;
            
        case 'shortlog':
            await showOptions('Shortlog', ['Default', 'Summary (-s)', 'By Email (-e)'],
                (opt) => window.go.main.App.Shortlog(opt));
            break;
            
        case 'reflog':
            await simpleOp('Reflog', () => window.go.main.App.Reflog('Standard'));
            break;
            
        case 'diff':
            await showOptions('Diff', ['Unstaged', 'Staged (--cached)', 'Names (--name-only)'],
                (opt) => window.go.main.App.Diff(opt));
            break;
            
        case 'stash':
            await showOptions('Stash', ['Save', 'Pop', 'List', 'Drop'],
                (opt) => window.go.main.App.Stash(opt.toLowerCase()));
            break;
            
        case 'clean':
            await showOptions('Clean', ['Preview (-n)', 'Force (-f)', 'Full (-fdx)'],
                (opt) => window.go.main.App.Clean(opt));
            break;
            
        case 'ls-files':
            await showOptions('Ls-Files', ['Staged', 'Tracked', 'Untracked', 'Modified'],
                (opt) => window.go.main.App.LsFiles(opt));
            break;
            
        case 'gitignore':
            await showGitIgnoreDialog();
            break;
            
        case 'reset':
            await showResetDialog();
            break;
            
        case 'rebase':
            await showRebaseDialog();
            break;
            
        case 'undo':
            await showConfirm('Undo Last Commit', 'Are you sure?', async () => {
                return await window.go.main.App.UndoLastCommit();
            });
            break;
            
        case 'worktree':
            await showWorktreeDialog();
            break;
            
        case 'conflicts':
            await showConflictsDialog();
            break;
            
        case 'blame':
            await showInputDialog('Blame', 'File path:', async (file) => {
                return await window.go.main.App.Blame(file);
            });
            break;
            
        case 'magic-sync':
            await simpleOp('Magic Sync', () => window.go.main.App.MagicSync());
            break;
    }
}

async function simpleOp(title, fn) {
    try {
        const result = await fn();
        consoleOutput.textContent = result;
    } catch (err) {
        consoleOutput.textContent = 'Error: ' + err;
    }
}

async function showOptions(title, options, fn) {
    showModal(title, `
        <select id="modal-select">
            ${options.map(o => `<option value="${o}">${o}</option>`).join('')}
        </select>
    `, async () => {
        const opt = document.getElementById('modal-select').value;
        const result = await fn(opt);
        consoleOutput.textContent = result;
    });
}

async function showInputDialog(title, placeholder, fn) {
    showModal(title, `
        <input type="text" id="modal-input" placeholder="${placeholder}">
    `, async () => {
        const input = document.getElementById('modal-input').value;
        if (!input.trim()) return;
        const result = await fn(input);
        consoleOutput.textContent = result;
    });
}

async function showCommitDialog() {
    showModal('Commit', `
        <input type="text" id="commit-msg" placeholder="Commit message">
        <select id="commit-opt">
            <option value="Standard (-m)">Standard (-m)</option>
            <option value="Stage All (-a)">Stage All (-a)</option>
            <option value="Amend (--amend)">Amend (--amend)</option>
        </select>
    `, async () => {
        const msg = document.getElementById('commit-msg').value;
        const opt = document.getElementById('commit-opt').value;
        if (!msg.trim()) return;
        const result = await window.go.main.App.Commit(msg, opt);
        consoleOutput.textContent = result;
    });
}

async function showBranchDialog(title, fn) {
    const branches = await window.go.main.App.ListBranches();
    showModal(title, `
        <select id="modal-branch">
            ${branches.map(b => `<option value="${b}">${b}</option>`).join('')}
        </select>
    `, async () => {
        const branch = document.getElementById('modal-branch').value;
        const result = await fn(branch);
        consoleOutput.textContent = result;
    });
}

async function showRenameDialog() {
    showModal('Rename Branch', `
        <input type="text" id="old-branch" placeholder="Current branch name">
        <input type="text" id="new-branch" placeholder="New branch name">
    `, async () => {
        const old = document.getElementById('old-branch').value;
        const new_ = document.getElementById('new-branch').value;
        if (!old.trim() || !new_.trim()) return;
        const result = await window.go.main.App.BranchRename(old, new_);
        consoleOutput.textContent = result;
    });
}

async function showTagDialog() {
    showModal('Tag', `
        <select id="tag-action">
            <option value="list">List</option>
            <option value="create">Create</option>
            <option value="delete">Delete</option>
        </select>
        <input type="text" id="tag-name" placeholder="Tag name" style="display:none">
    `, async () => {
        const action = document.getElementById('tag-action').value;
        const name = document.getElementById('tag-name').value;
        if (action !== 'list' && !name.trim()) return;
        const result = await window.go.main.App.Tag(action, name);
        consoleOutput.textContent = result;
    });
    
    document.getElementById('tag-action').addEventListener('change', (e) => {
        document.getElementById('tag-name').style.display = 
            e.target.value === 'list' ? 'none' : 'block';
    });
}

async function showRemoteDialog() {
    showModal('Remote', `
        <select id="remote-action">
            <option value="list">List</option>
            <option value="remove">Remove</option>
            <option value="add">Add</option>
        </select>
        <input type="text" id="remote-args" placeholder="Remote name/URL" style="display:none">
    `, async () => {
        const action = document.getElementById('remote-action').value;
        const args = document.getElementById('remote-args').value;
        const result = await window.go.main.App.Remote(action, args);
        consoleOutput.textContent = result;
    });
    
    document.getElementById('remote-action').addEventListener('change', (e) => {
        document.getElementById('remote-args').style.display = 
            e.target.value === 'list' ? 'none' : 'block';
    });
}

async function showGitIgnoreDialog() {
    const content = await window.go.main.App.GitIgnore('read');
    showModal('Edit .gitignore', `
        <textarea id="gitignore-content" style="height:200px">${content || ''}</textarea>
    `, async () => {
        const newContent = document.getElementById('gitignore-content').value;
        await window.go.main.App.SaveGitIgnore(newContent);
        consoleOutput.textContent = 'Gitignore saved successfully';
    });
}

async function showResetDialog() {
    showModal('Reset', `
        <input type="text" id="reset-target" value="HEAD~1" placeholder="Commit hash">
        <select id="reset-mode">
            <option value="--mixed">--mixed</option>
            <option value="--soft">--soft</option>
            <option value="--hard">--hard</option>
        </select>
    `, async () => {
        const target = document.getElementById('reset-target').value;
        const mode = document.getElementById('reset-mode').value;
        const result = await window.go.main.App.Reset(mode, target);
        consoleOutput.textContent = result;
    });
}

async function showRebaseDialog() {
    showModal('Rebase', `
        <select id="rebase-opt">
            <option value="Interactive (-i)">Interactive (-i)</option>
            <option value="Continue">Continue</option>
            <option value="Abort">Abort</option>
            <option value="Skip">Skip</option>
        </select>
        <input type="text" id="rebase-target" placeholder="Target branch" style="display:none">
    `, async () => {
        const opt = document.getElementById('rebase-opt').value;
        const target = document.getElementById('rebase-target').value;
        const result = await window.go.main.App.Rebase(opt, target);
        consoleOutput.textContent = result;
    });
    
    document.getElementById('rebase-opt').addEventListener('change', (e) => {
        const needTarget = !['Continue', 'Abort', 'Skip'].includes(e.target.value);
        document.getElementById('rebase-target').style.display = 
            needTarget ? 'block' : 'none';
    });
}

async function showWorktreeDialog() {
    showModal('Worktree', `
        <select id="wt-action">
            <option value="List">List</option>
            <option value="Add">Add</option>
            <option value="Remove">Remove</option>
            <option value="Prune">Prune</option>
        </select>
        <input type="text" id="wt-args" placeholder="Path/branch" style="display:none">
    `, async () => {
        const action = document.getElementById('wt-action').value;
        const args = document.getElementById('wt-args').value;
        const result = await window.go.main.App.Worktree(action, args);
        consoleOutput.textContent = result;
    });
    
    document.getElementById('wt-action').addEventListener('change', (e) => {
        const needArgs = ['Add', 'Remove'].includes(e.target.value);
        document.getElementById('wt-args').style.display = needArgs ? 'block' : 'none';
    });
}

async function showConflictsDialog() {
    const conflicts = await window.go.main.App.GetConflicts();
    if (conflicts.length === 0) {
        consoleOutput.textContent = 'No merge conflicts detected.';
        return;
    }
    
    const html = conflicts.map(f => `
        <div style="display:flex;justify-content:space-between;margin:5px 0;padding:10px;background:#1a1b26;border-radius:4px">
            <span>${f}</span>
            <div>
                <button class="btn btn-small" onclick="resolveConflict('${f}','ours')">Keep Mine</button>
                <button class="btn btn-small" onclick="resolveConflict('${f}','theirs')">Take Theirs</button>
            </div>
        </div>
    `).join('');
    
    showModal('Resolve Conflicts', html, null);
}

async function resolveConflict(file, strategy) {
    const result = await window.go.main.App.ResolveConflict(file, strategy);
    consoleOutput.textContent = result;
    hideModal();
}

async function showConfirm(title, message, fn) {
    showModal(title, `<p>${message}</p>`, async () => {
        const result = await fn();
        consoleOutput.textContent = result;
    });
}

function showModal(title, bodyHtml, onConfirm) {
    modalTitle.textContent = title;
    modalBody.innerHTML = bodyHtml;
    modal.classList.remove('hidden');
    
    modalCancel.onclick = hideModal;
    
    if (onConfirm) {
        modalConfirm.onclick = async () => {
            await onConfirm();
            hideModal();
        };
        modalConfirm.style.display = 'block';
    } else {
        modalConfirm.style.display = 'none';
    }
}

function hideModal() {
    modal.classList.add('hidden');
}

// Documentation
function initDocs() {
    const docList = document.getElementById('doc-list');
    const docKeys = Object.keys(docContent);
    
    docKeys.forEach(key => {
        const btn = document.createElement('button');
        btn.className = 'doc-btn';
        btn.textContent = key;
        btn.addEventListener('click', () => {
            document.getElementById('doc-content').textContent = docContent[key];
            document.getElementById('doc-content').classList.remove('hidden');
        });
        docList.appendChild(btn);
    });
}

// Init
initDocs();
```

- [ ] **Step 2: Commit**

```bash
git add gitscope-wails/frontend/app.js
git commit -m "feat: add frontend JavaScript with all operation handlers"
```

---

### Task 5: Remove Fyne Dependencies from Git Package

**Covers:** [S2]

**Files:**
- Modify: `gitscope-wails/internal/git/git_ops.go`
- Create: `gitscope-wails/internal/state/state.go`

- [ ] **Step 1: Create state package**

```go
package state

var RepoPath string
```

- [ ] **Step 2: Create git_ops.go (copy from original, remove Fyne deps)**

Copy `internal/git/git_go.go` to `gitscope-wails/internal/git/git_ops.go` and remove:
- All `fyne.io/fyne/v2` imports
- `fyne.Window` parameters
- `widget.Entry` parameters
- Dialog calls (these will be handled by frontend)

Remove the `GitIgnore` function that uses Fyne widgets - replace with simpler version:

```go
func GitIgnoreRead(repoPath string) (string, error) {
	filePath := filepath.Join(repoPath, ".gitignore")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
```

- [ ] **Step 3: Verify compilation**

Run: `cd gitscope-wails && go build ./...`
Expected: No errors

- [ ] **Step 4: Commit**

```bash
git add gitscope-wails/internal/git/git_ops.go gitscope-wails/internal/state/state.go
git commit -m "refactor: remove Fyne dependencies from git package"
```

---

### Task 6: Create GitHub Actions CI/CD Workflow

**Covers:** [S6]

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create release.yml**

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: windows
            goarch: amd64
            ext: .exe
          - goos: linux
            goarch: amd64
            ext: ''
          - goos: darwin
            goarch: amd64
            ext: ''
          - goos: darwin
            goarch: arm64
            ext: ''
    
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

      - name: Extract version from tag
        id: version
        run: echo "VERSION=${GITHUB_REF#refs/tags/v}" >> $GITHUB_OUTPUT

      - name: Build
        working-directory: gitscope-wails
        run: wails build -platform ${{ matrix.goos }}/${{ matrix.goarch }} -o gitscope-wails${{ matrix.ext }}

      - name: Create artifact
        uses: actions/upload-artifact@v4
        with:
          name: gitscope-wails-${{ matrix.goos }}-${{ matrix.goarch }}
          path: gitscope-wails/build/bin/gitscope-wails${{ matrix.ext }}

  publish:
    needs: release
    runs-on: ubuntu-latest
    steps:
      - name: Download all artifacts
        uses: actions/download-artifact@v4

      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          generate_release_notes: true
          files: |
            gitscope-wails-windows-amd64/*
            gitscope-wails-linux-amd64/*
            gitscope-wails-darwin-amd64/*
            gitscope-wails-darwin-arm64/*
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add GitHub Actions workflow for cross-platform releases"
```

---

### Task 7: Final Integration and Testing

**Covers:** [S7]

**Files:**
- Modify: `gitscope-wails/wails.json`

- [ ] **Step 1: Update wails.json**

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "gitscope-wails",
  "outputfilename": "gitscope-wails",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "Aswanidev VS"
  }
}
```

- [ ] **Step 2: Initialize npm and install dependencies**

```bash
cd gitscope-wails/frontend
npm init -y
npm install
```

- [ ] **Step 3: Run Wails dev to test**

```bash
cd gitscope-wails
wails dev
```

- [ ] **Step 4: Test all operations**

1. Select a repository
2. Test Init, Status, Stage, Commit
3. Test Branch create/delete
4. Test Push, Pull
5. Test all dashboard operations

- [ ] **Step 5: Create release tag**

```bash
git add -A
git commit -m "feat: complete Wails port with CI/CD"
git tag v1.0.0
git push origin v1.0.0
```

- [ ] **Step 6: Verify GitHub Actions triggered**

Go to GitHub repository -> Actions tab -> verify release workflow is running
