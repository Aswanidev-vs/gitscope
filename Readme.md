<p align="center">
  <img width="150" height="150" src="https://github.com/user-attachments/assets/992700cd-57ce-4866-83c0-8d71cf27bbec" />
</p>


# **GitScope**

GitScope is a lightweight and intuitive Git desktop client built with **Go** and the **Fyne** GUI toolkit. It provides a clean interface for performing essential Git operations without needing the command line, while still using Git’s full power internally.

GitScope is ideal for developers who want a simple, cross-platform Git companion for everyday tasks such as committing, branching, pushing, pulling, and browsing repository history.

## **Features**

* **Repository Setup**
  Select, create, or initialize repositories using a simple folder picker.

* **Initialization and Status**
  Initialize new Git repositories and view repository status directly in the UI.

* **Staging and Committing**
  Stage changes and create commits with a custom message.

* **Branch Management**
  Create, delete, switch, and rename branches.

* **Push and Pull**
  Push commits to remotes and pull changes with branch selection.

* **Clone Repositories**
  Clone any remote repository into an empty folder.

* **Logs and History**
  View commit history, reflog entries, and revert specific commits.

* **GitIgnore Editor**
  Create or edit the `.gitignore` file inside the app.

* **Custom Git Commands**
  Run user-defined Git command sequences for advanced workflows.

* **Cross-Platform**
  Runs on Windows, macOS, and Linux using Fyne’s cross-platform GUI engine.

## **Prerequisites**

* **Go**: Version 1.25.3 or later
* **Git**: Installed and accessible from PATH
* **Fyne**: Pulled automatically through Go modules; additional setup may be required for some platforms

---

## **Installation**

### **Build from Source**

```bash
git clone https://github.com/Aswanidev-vs/GitScope.git
cd GitScope
go mod tidy
go run main.go
```

### **Build Executable**

```bash
go build -o gitscope main.go
```

Run with:

* Windows: `gitscope.exe`
* macOS/Linux: `./gitscope`

---

## **Usage**

1. **Start the App**
   Launch using `go run main.go` or the built executable.

2. **Select a Repository**

   * Pick an existing folder
   * Or create a new folder and initialize it using the "Init" function

3. **Perform Git Operations**
   Use the dashboard to run all supported operations:

   * Init
   * Status
   * Stage All
   * Commit
   * Branch Create/Delete
   * Branch Switch
   * Branch Rename
   * Push
   * Pull
   * Clone
   * Log
   * Reflog
   * Revert commit by hash
   * Edit `.gitignore`
   * Run custom Git commands

The app uses `os/exec` to run Git commands. Ensure Git is installed and on your PATH.

---

## **Project Structure**

```
GitScope/
│
├── main.go                 # Entry point
├── assets/                 # Icons and static assets
├── go.mod / go.sum         # Module files
│
├── internal/
│   ├── core/               # GUI forms, manager logic
│   │   └── manager.go
│   ├── git/                # Git integration and command wrappers
│   │   └── git_go.go
│   ├── helpers/            # Shell, branch selectors, utility helpers
│   ├── state/              # Runtime repo state
│   └── ui/                 # Fyne UI implementation
│       ├── fyne_app.go
│       └── fyne_views.go
│
└── utils/                  # Dialogs and UI helpers
```

---

## **Developer Notes**

* Git commands are executed using `exec.Command` within the selected repository directory.
* Global state is stored in `internal/state`, primarily `RepoPath`.
* Commands are designed to be cross-platform, with shell helpers ensuring portability.
* Future tests should include:

  * Unit tests for Git command wrappers
  * Integration tests using temporary Git repositories

---

## **Planned Improvements**

* Persistent repository history across app restarts
* GUI improvements and theme customization
* Simple merge and conflict-viewing support
* Installer packages for each OS
* Basic Git configuration editor (name, email, remotes)
* Automated unit and integration tests

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

---

## **Contact**

**Developer:** Aswanidev VS
**GitHub:** [https://github.com/Aswanidev-vs/GitScope](https://github.com/Aswanidev-vs/GitScope)
