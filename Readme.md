# GitScope

GitScope is a modern, lightweight, and visually intuitive desktop GUI client for Git, built with Go and the Fyne toolkit. It simplifies essential version control operations, making Git more accessible for beginners and experienced developers alike. GitScope provides a user-friendly interface to perform common Git tasks without needing the command line, while still leveraging Git's powerful features under the hood.

## Key Features

- **Repository Management**: Select or create local repositories using a folder picker.
- **Initialization & Status**: Initialize new Git repositories and view repository status.
- **Staging & Committing**: Stage changes and create commits with custom messages.
- **Branching**: Create, delete, switch, and rename branches.
- **Pushing & Pulling**: Push commits to remote repositories and pull changes with branch selection.
- **Cloning**: Clone repositories into empty folders.
- **History & Logs**: View commit logs, reflog, and revert commits by hash.
- **GitIgnore Editing**: Edit or create `.gitignore` files directly in the app.
- **Custom Commands**: Run user-provided Git command sequences for advanced operations.
- **Cross-Platform**: Works on Windows, macOS, and Linux using Fyne's GUI toolkit.

## Prerequisites

- **Go**: Version 1.25.3 or later (module requires go 1.25.3).
- **Git**: Installed and available in your system's PATH.
- **Fyne Dependencies**: Automatically pulled via `go mod` (optional manual setup for development).

## Installation

### Option 1: Build from Source (Recommended for Development)

1. Clone the repository:
   ```bash
   git clone https://github.com/Aswanidev-vs/GitScope.git
   cd GitScope
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run main.go
   ```

### Option 2: Build Executable

1. Build the binary:
   ```bash
   go build -o gitscope main.go
   ```

2. Run the executable:
   - On Windows: `gitscope.exe`
   - On macOS/Linux: `./gitscope`

The built binary can be executed directly on supported platforms.

## Usage

1. **Launch the App**: Run `go run main.go` or the built executable. The app checks for Git availability on startup.

2. **Repository Setup**:
   - Use "Select Repository" to choose an existing local folder.
   - Use "Create New Repository" to paste commands for initializing a new repo (e.g., GitHub creation commands).
   - Use "Existing Repository" for additional setup options.

3. **Dashboard Operations**:
   - **Init**: Initialize a Git repository in the selected folder.
   - **Status**: View the current repository status.
   - **Stage**: Add all changes to the staging area.
   - **Commit**: Create a commit with a custom message.
   - **Push**: Push commits to the remote repository (select branch).
   - **Pull**: Pull changes from the remote (option to reset last commit).
   - **Log**: View commit history in oneline format.
   - **Reflog**: View reference logs.
   - **Clone**: Clone a repository into an empty folder (provide GitHub URL).
   - **Branch**: Create or delete branches.
   - **Switch Branch**: Switch to a different branch.
   - **Rename Branch**: Rename an existing branch.
   - **Revert**: Revert a commit by hash.
   - **GitIgnore (Edit)**: Edit the `.gitignore` file directly.

4. **Settings**: View app information, version, and links.

**Note**: The app executes Git commands via `os/exec`. Ensure Git is installed and in PATH. For security, avoid running untrusted command blocks.

## Project Layout

- `main.go` — Application entry point; initializes the UI.
- `go.mod` & `go.sum` — Go module and dependency declarations.
- `assets/` — Icons and static assets (app icon for the GUI).
- `internal/`
  - `core/manager.go` — GUI components for repository selection and forms.
  - `git/git_go.go` — Thin wrapper functions executing Git commands (init, status, commit, stage, push, clone, branch, pull, revert, log, reflog, etc.).
  - `helpers/` — Utility functions for shell commands, branch selectors, and GUI helpers.
  - `state/` — Global runtime state (e.g., current repository path).
  - `ui/` — Fyne app and views (`fyne_app.go`, `fyne_views.go`) implementing the GUI.
- `utils/` — Additional UI helpers (e.g., branch creation dialogs).
- `Readme.md` — This file.

## Developer Notes

- **Git Execution**: Commands are run using `exec.Command` with `git` in the repository directory.
- **Global State**: `internal/state.RepoPath` stores the selected folder; handle changes carefully.
- **Cross-Platform Shell**: Uses `helpers.RunShellCommand` for portability, with some Windows-specific paths.
- **Dependencies**: Relies on Fyne (fyne.io/fyne/v2) and fyne-tooltip for UI.
- **Testing**: Consider unit tests for Git wrappers and integration tests with temporary repos.

## Suggestions for Future Development

- Add unit tests for Git wrapper functions (mock `exec.Command`).
- Implement integration tests using temporary folders and local Git.
- Persist selected repository across app restarts.
- Create installers/packages for each OS.
- Add support for Git configurations and credentials.

## Contributing

Contributions are welcome! Please:

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Submit a pull request.

For issues or discussions, open an issue on GitHub.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Contact

- **Developer**: Aswanidev VS
- **GitHub**: [https://github.com/Aswanidev-vs/GitScope](https://github.com/Aswanidev-vs/GitScope)

Built with ❤️ using the [Fyne GUI toolkit](https://fyne.io/).
