<p align="center">
  <img width="150" height="150" src="https://github.com/user-attachments/assets/992700cd-57ce-4866-83c0-8d71cf27bbec" />
</p>


# **GitScope**

GitScope is a modern Git desktop client built with **Go** and **[Wails v2](https://wails.io)**. It provides a fast, responsive web-based UI for essential Git operations without needing the command line, while using Git's full power internally.

GitScope is ideal for developers who want a simple, cross-platform Git companion for everyday tasks such as committing, branching, pushing, pulling, and browsing repository history.

## **Features**

* **Repository Setup**
  Select, create, or initialize repositories using a native folder picker.

* **Initialization and Status**
  Initialize new Git repositories and view repository status in multiple formats (standard, short, branch).

* **Staging and Committing**
  Stage all changes and create commits with messages and options (stage-all, amend).

* **Branch Management**
  Create, delete, switch, and rename branches with remote upstream support.

* **Push and Pull**
  Push commits to remotes with upstream tracking and pull changes with branch selection.

* **Clone Repositories**
  Clone any remote repository into a selected destination folder.

* **Logs and History**
  View commit history in oneline, graph, or pretty format. Browse reflog entries and revert specific commits.

* **Diff**
  View unstaged, staged, named-only, or stat diffs from the dashboard.

* **Stash**
  Save, pop, list, drop, or apply stashes.

* **Tag**
  List, create, delete, and push tags.

* **Remote Management**
  List, add, and remove remote connections.

* **Merge and Rebase**
  Merge branches and perform rebase operations (continue, abort, skip, onto).

* **Reset**
  Soft, mixed, or hard resets to any target.

* **Clean**
  Preview or force-remove untracked files and directories.

* **Cherry-pick**
  Apply specific commits to the current branch.

* **Worktree**
  Manage multiple working trees.

* **Blame**
  Show per-line blame annotations for any file.

* **Shortlog**
  Summarize commit history grouped by author.

* **Show**
  Display full details, stats, or patches for any commit.

* **Magic Sync**
  One-click stash, fetch, pull --rebase, and stash pop workflow.

* **Conflict Resolution**
  Detect merge conflicts and resolve with "keep mine" or "take theirs" strategies.

* **GitIgnore Editor**
  Create or edit the `.gitignore` file directly inside the app.

* **Custom Git Commands**
  Run user-defined Git command sequences for advanced workflows.

* **Console Output**
  Color-coded console panel showing command results, errors, and progress in real time.

* **Cross-Platform**
  Runs on Windows, macOS, and Linux.

## **Prerequisites**

* **Go**: Version 1.21 or later
* **Node.js**: Version 18 or later (for frontend build)
* **Git**: Installed and accessible from PATH
* **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

---

## **Installation**

### **Build from Source**

```bash
git clone https://github.com/Aswanidev-vs/GitScope.git
cd GitScope
go mod tidy
cd gitscope-wails
wails build
```

### **Development Mode**

```bash
cd gitscope-wails
wails dev
```

This starts the Vite dev server with hot reload for the frontend and live-reload for the Go backend.

### **Build Production Binary**

```bash
cd gitscope-wails
wails build
```

The binary will be at `build/bin/gitscope-wails` (or `gitscope-wails.exe` on Windows).

---

## **Usage**

1. **Start the App**
   Launch using `go run ./gitscope-wails/` or the built executable.

2. **Select a Repository**
   * Pick an existing folder
   * Or create a new folder and initialize it using the "Init" button

3. **Perform Git Operations**
   Use the dashboard to run all supported operations:

   | Category | Commands |
   |----------|----------|
   | Common | Init, Stage All, Status, Commit, Push, Log |
   | Branches | Branch Create/Delete, Switch, Merge, Rename, Tag |
   | Remote | Remote, Fetch, Pull, Clone, Cherry-pick |
   | History | Log, Revert, Show, Shortlog, Reflog |
   | Changes | Diff, Stash, Clean, Ls-Files, .gitignore |
   | Advanced | Reset, Rebase, Undo, Worktree, Conflicts |
   | Tools | Blame, Magic Sync |

   You can also run custom Git commands from the Repository Setup page.

The app uses `os/exec` to run Git commands. Ensure Git is installed and on your PATH.

---

## **Project Structure**

```
GitScope/
│
├── main.go                     # Fyne-based GUI entry point (legacy)
├── go.mod / go.sum             # Shared Go module
│
├── gitscope-wails/             # Wails-based GUI (primary)
│   ├── main.go                 # Wails app entry point
│   ├── app.go                  # Bound Go methods exposed to frontend
│   ├── app_test.go             # Unit tests
│   ├── wails.json              # Wails project config
│   ├── hidewindow_windows.go   # Windows-specific process flags
│   ├── hidewindow_other.go     # Cross-platform no-op
│   ├── frontend/               # Vue/JS frontend
│   │   ├── index.html
│   │   ├── vite.config.js
│   │   ├── package.json
│   │   ├── src/
│   │   │   ├── main.js         # UI logic, dialogs, console
│   │   │   ├── style.css       # Theme and layout
│   │   │   └── icons.js        # SVG icon library
│   │   └── dist/               # Built frontend assets
│   └── build/                  # Platform-specific build config
│       ├── appicon.png
│       └── windows/icon.ico
│
├── internal/
│   ├── git/
│   │   ├── git_go.go           # Git command wrappers (all platforms)
│   │   ├── hidewindow_windows.go
│   │   └── hidewindow_other.go
│   ├── state/
│   │   └── state.go            # Runtime repo state (RepoPath)
│   ├── ui/                     # Fyne UI (legacy)
│   │   ├── fyne_app.go
│   │   └── fyne_views.go
│   ├── helpers/                # Shell utilities (legacy Fyne)
│   ├── doc/                    # Documentation views (legacy Fyne)
│   └── core/                   # Form manager (legacy Fyne)
│
├── utils/                      # Legacy Fyne UI helpers
├── assets/                     # Icons and static assets
└── .github/workflows/
    └── release.yml             # CI/CD pipeline
```

---

## **Developer Notes**

* Git commands are executed using `os/exec` within the selected repository directory.
* Global state is stored in `internal/state`, primarily `RepoPath`.
* Platform-specific code (e.g., `HideWindow` on Windows) uses build tags for cross-platform compilation.
* The frontend communicates with Go via Wails bindings — no REST/WebSocket boilerplate needed.
* The Fyne-based GUI (`main.go`, `internal/ui/`, `utils/`) is legacy and maintained separately.
* Automated tests run via `go test ./gitscope-wails/...`.

---

## **Planned Improvements**

* Persistent repository history across app restarts
* GUI improvements and theme customization
* Installer packages for each OS (via `wails build`)
* Basic Git configuration editor (name, email, remotes)
* Expanded automated test coverage

---

## **Contributing**

Contributions are welcome.

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Open a pull request

For bugs or feature requests, open an issue with details and reproduction steps.

---

## **License**

This project is licensed under the **MIT License**.
See the **LICENSE** file for complete details.
