# GitScope: Fyne to Wails Port & CI/CD Design

## [S1] Problem
Port the GitScope desktop GUI from Go+Fyne to Go+Wails, and implement GitHub Actions CI/CD for automatic cross-platform releases triggered by git tags.

## [S2] Solution Overview
### Backend (Go + Wails)
- Reuse existing `internal/git/` package (all git operations) - decouple from Fyne by removing `fyne.Window` params
- Reuse `internal/state/` for RepoPath
- New `internal/helpers/` - pure Go functions (no UI deps)
- New `app.go` - Wails app binding exposing all git operations to frontend

### Frontend (Plain HTML/JS/CSS)
- Sidebar with 4 tabs: Repository, Dashboard, Settings, Documentation
- Dashboard with 7 categories (Common, Branches, Remote, History, Changes, Advanced, Tools)
- Each button calls Wails bindings -> Go backend -> returns output
- Simple CSS styling matching the Fyne violet theme

### CI/CD (GitHub Actions)
- Trigger: push tag `v*`
- Build matrix: windows/amd64, linux/amd64, darwin/amd64, darwin/arm64
- Auto-create GitHub Release with binaries
- Version extracted from tag (e.g., v1.2.3 -> 1.2.3)

## [S3] Component Breakdown

### 3.1 Backend Files
| File | Purpose |
|------|---------|
| `main.go` | Wails app entry, embed assets, bind App |
| `app.go` | App struct with all git operation methods |
| `internal/state/state.go` | Global RepoPath (unchanged) |
| `internal/git/*.go` | Git operations (refactored to remove Fyne deps) |

### 3.2 Frontend Files
| File | Purpose |
|------|---------|
| `frontend/index.html` | Main HTML structure |
| `frontend/style.css` | Styling with violet theme |
| `frontend/app.js` | All UI logic and Wails bindings |

### 3.3 CI/CD
| File | Purpose |
|------|---------|
| `.github/workflows/release.yml` | Build + release on tag push |

## [S4] App Method Signatures (Go Backend)

```go
// Repository
func (a *App) OpenRepo(path string) error
func (a *App) GetRepoPath() string
func (a *App) Init() (string, error)
func (a *App) Clone(url string) (string, error)
func (a *App) NewRepoCmd(cmdText string) (string, error)
func (a *App) ExistingRepoCmd(cmdText string) (string, error)

// Git Operations
func (a *App) Status(option string) (string, error)
func (a *App) Stage(option string) (string, error)
func (a *App) Commit(msg, option string) (string, error)
func (a *App) Push(branch string) (string, error)
func (a *App) Pull(branch string) (string, error)
func (a *App) Log(option string) (string, error)
func (a *App) Diff(option string) (string, error)
func (a *App) Reset(mode, target string) (string, error)
func (a *App) Fetch(option string) (string, error)
func (a *App) Stash(action string) (string, error)
func (a *App) Merge(branchname string) (string, error)
func (a *App) Tag(action, tagname string) (string, error)
func (a *App) BranchCreate(name string) (string, error)
func (a *App) BranchDelete(name string) (string, error)
func (a *App) BranchRename(oldname, newname string) (string, error)
func (a *App) SwitchBranch(branchname string) (string, error)
func (a *App) Revert(commitHash string) (string, error)
func (a *App) Remote(action, args string) (string, error)
func (a *App) Show(option, target string) (string, error)
func (a *App) LsFiles(option string) (string, error)
func (a *App) Blame(file string) (string, error)
func (a *App) Worktree(action, args string) (string, error)
func (a *App) Shortlog(option string) (string, error)
func (a *App) Reflog(option string) (string, error)
func (a *App) Rebase(option, target string) (string, error)
func (a *App) Clean(option string) (string, error)
func (a *App) CherryPick(hash string) (string, error)
func (a *App) MagicSync() (string, error)
func (a *App) UndoLastCommit() (string, error)
func (a *App) GetConflicts() ([]string, error)
func (a *App) ResolveConflict(file, strategy string) (string, error)
func (a *App) GitIgnore(action string) (string, error)
func (a *App) ListBranches() ([]string, error)
```

## [S5] Frontend Structure

```
frontend/
├── index.html          # Main layout: sidebar + content area
├── style.css           # Violet theme, responsive grid
├── app.js              # Wails binding calls, DOM manipulation
└── wailsjs/            # Auto-generated Wails JS bindings
```

## [S6] CI/CD Workflow

```yaml
# Trigger on tag push v*
# Steps:
# 1. Checkout code
# 2. Set up Go
# 3. Install Wails
# 4. Build for each platform (matrix)
# 5. Upload artifacts
# 6. Create GitHub Release with all binaries
```

## [S7] Migration Checklist

- [ ] Remove Fyne imports from git package
- [ ] Decouple helpers from fyne.Window
- [ ] Create Wails App struct with all methods
- [ ] Build frontend HTML/CSS/JS
- [ ] Test all git operations
- [ ] Create GitHub Actions workflow
- [ ] Test release on tag push
