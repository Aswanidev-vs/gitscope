# gitscope

Gitscope is a small desktop GUI tool written in Go that provides a lightweight visual interface for common Git operations. It uses the Fyne GUI toolkit and executes Git commands under the hood, making it useful as a helper app for initializing repositories, staging/committing changes, pushing/pulling, branch management, and running custom git command sequences.

Key features
- Select or create a repository folder using a folder picker
- Initialize a Git repo, view status, stage changes, create commits
- Push and pull with branch selection
- Clone repositories into an empty folder
- Create / delete branches through a dialog
- Revert a commit by hash
- Run a block of user-provided git commands or repo-creation commands (useful to paste a set of GitHub creation commands)

Quick start

Prerequisites
- Go 1.25+ installed (module requires go 1.25.3)
- Git installed and available on PATH
- (Optional) Fyne development dependencies are pulled by `go mod` automatically

Run locally (development)

1. Clone the repo and change to its directory

```bash
cd g:/gitscope
```

2. Run the app

```bash
go run main.go
```

Build

```bash
go build -o gitscope.exe main.go
```

On supported OSes the produced binary can be executed directly. The GUI uses Fyne and should work across desktop platforms (Windows, macOS, Linux). Ensure you have Git installed and accessible.

How to use
- Start the app.
- Use "Select Repository" to choose or create a local target folder.
- Use the Dashboard tab to run Init, Stage, Commit, Push, Pull, Log, Clone, Revert and Branch operations.
- Use "Create New Repository" to paste a set of commands (for creating a GitHub repo and pushing). The app will attempt to run those commands in the selected folder.

Project layout

- `main.go` — program entrypoint; calls `ui.App()`
- `go.mod` — module and dependency declarations
- `assets/` — icons and static assets (app icon used by the GUI)
- `internal/`
	- `core/manager.go` — higher-level GUI components for repo selection and new/existing repo forms
	- `git/` — thin wrapper functions executing git commands (`init`, `status`, `commit`, `stage`, `push`, `clone`, `branch`, `pull`, `revert`, `log`)
	- `helpers/` — helper utilities to run shell commands, branch selectors, run pasted command blocks and other GUI helpers
	- `state/` — global runtime state (currently contains `RepoPath` string)
	- `ui/` — Fyne app and views (`fyne_app.go`, `fyne_views.go`) implementing the GUI
- `utils/` — small UI helpers (e.g. branch creation dialog)
- `Readme.md` — this file

Developer notes & known behavior
- The app runs Git by calling external `git` commands (via `os/exec`). Git must be present in PATH.
- Global state: `internal/state.RepoPath` is used widely to store the current selected folder. Be careful when changing it.
- The helpers package has a `RullshellCommand` helper for cross-platform shell commands, but some code paths call `cmd /C` directly on Windows-specific commands; review these paths if running on non-Windows platforms.
- The UI uses Fyne’s widget library plus a tooltip helper (`dweymouth/fyne-tooltip`).

Security & safety notes
- The app executes arbitrary shell/git commands pasted in by the user (intended feature). Do not run untrusted command blocks.

Suggestions / next steps
- Add unit tests for the git wrapper functions by mocking `exec.Command` in a wrapper.
- Add integration tests that run on CI with a temporary folder and a local Git binary.
- Add a proper configuration file / app settings to persist selected repo across runs.
- Add an installer / packaging instructions for each OS.

Contributing

Contributions are welcome. Please open an issue for discussion or submit a pull request.

License

 Add a LICENSE (for example MIT or Apache-2.0) to make the project terms explicit.

Contact

Author: project repository (see commit history). For questions, open an issue in the repository.

Acknowledgements

Built using the Fyne GUI toolkit (fyne.io/fyne/v2).

