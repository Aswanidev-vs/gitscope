package doc

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func FormatSection(title, body string) []widget.RichTextSegment {
	return []widget.RichTextSegment{
		&widget.TextSegment{
			Text:  title,
			Style: widget.RichTextStyleHeading,
		},
		&widget.TextSegment{
			Text:  body,
			Style: widget.RichTextStyleParagraph,
		},
	}
}

func Init() fyne.CanvasObject {
	body := `git init creates a new empty Git repository in your current folder. It tells Git to start tracking changes in that directory.

Think of it as telling Git:
"Start a new project here and track everything I do."

Command:
git init

After running this, Git creates a hidden folder:
.git

This folder stores all version history and Git settings.

Example usage:
mkdir myproject
cd myproject
git init

Now 'myproject' becomes a Git repository and you can start adding files and making commits.`

	content := FormatSection("Git Init", body)
	return widget.NewRichText(content...)
}

func Stage() fyne.CanvasObject {
	body := `The stage (also called the staging area or index) is a place where Git stores the changes you want to include in your next commit.

Simple definition:
The stage is where you put files that Git will commit.

Why staging exists:
• choose which files to commit
• commit only specific parts
• prepare commits cleanly

Commands:

Stage a single file:
git add file.txt

Stage all files:
git add .`

	content := FormatSection("Git Stage", body)
	return widget.NewRichText(content...)
}

func Status() fyne.CanvasObject {
	body := `
git status shows the current state of your working directory and staging area.It tells you what changed, what is staged, what is
unstaged, and which files are untracked.Think of it as a quick dashboard that shows what will and won't be included in your next
commit.It helps you confirm which files you modified, which ones you staged, and whether you are ready to commit. It is one of
the safest and most frequently used Git commands.

Command

Show current status:

git status

Example output

Changes not staged for commit:
	modified:   main.go

Untracked files:
	newfile.txt

Changes to be committed:
	modified:   README.md
	`
	content := FormatSection("Git Status", body)
	return widget.NewRichText(content...)
}
func Commit() fyne.CanvasObject {
	body := `
A commit is a snapshot of your project at a specific point in time. It records the changes you staged and creates a permanent
history entry.

Simple definition:
A commit saves your staged changes with a message describing what changed.

Why commits are important:
• they create a history of your work  
• you can go back to earlier versions  
• they help track progress  
• easier collaboration and debugging  

Commands:

Create a commit:
git commit -m "your message"

Commit with detailed message editor:

git commit

Commit all tracked changes (skip staging):

git commit -a -m "message"

	`
	content := FormatSection("Git Commit", body)
	return widget.NewRichText(content...)
}
func Push() fyne.CanvasObject {
	body := `
Push is the action of sending your local commits to a remote repository like GitHub or GitLab.

Simple definition:
Push uploads your commits so others can see them.

Why push is needed:
• share your work with others  
• back up your commits online  
• update the remote branch with your latest changes  
• collaborate smoothly across devices

Commands:

Push current branch:
git push

Push a specific branch:
git push origin branchname

Push and create the branch on remote:
git push -u origin branchname

Force push only if needed:
git push --force

View all remotes:
git remote -v
	`
	content := FormatSection("Git Push", body)
	return widget.NewRichText(content...)
}
func Log() fyne.CanvasObject {
	body := `Log shows the full history of commits in your repository.

Simple definition:
Log is the timeline of all commits.

Why log is useful:
• see what changed over time  
• view commit messages and authors  
• inspect commit IDs for operations like reset or checkout  
• understand project progress  
• debug by reviewing past states

Common commands:

View full commit history:
git log

View one line per commit:
git log --oneline

View graph view:
git log --oneline --graph --decorate --all

View changes made in each commit:
git log -p

View commits by a specific author:
git log --author="name"

View commits for a single file:
git log file.txt`

	content := FormatSection("Git Log", body)
	return widget.NewRichText(content...)
}
func Revert() fyne.CanvasObject {
	body := `Revert lets you undo a specific commit safely by creating a new commit.

Simple definition:
Revert makes a new commit that cancels the changes of an older commit without changing history.

Why revert is useful:
• safe way to undo mistakes  
• does not rewrite commit history  
• good for shared or public branches  
• fixes issues without breaking collaborators  
• keeps the timeline clean and traceable

Common commands:

Revert a single commit:
git revert <commit-id>

Revert multiple commits (interactive):
git revert --no-commit <old-id>..<new-id>
git commit

Abort a revert if something goes wrong:
git revert --abort`

	content := FormatSection("Git Revert", body)
	return widget.NewRichText(content...)
}
func Clone() fyne.CanvasObject {
	body := `Clone creates a local copy of a remote repository on your system.

Simple definition:
Clone downloads a full repository including its history so you can work on it locally.

Why clone is useful:
• get a complete working copy of any remote repo  
• access full commit history  
• start contributing or editing immediately  
• set up the remote connection automatically (origin)

Common commands:

Clone a repository:
git clone <repo-url>

Clone into a custom folder:
git clone <repo-url> myproject

Clone only the latest commit (faster):
git clone --depth 1 <repo-url>

Check the remote linked to your cloned repo:
git remote -v`

	content := FormatSection("Git Clone", body)
	return widget.NewRichText(content...)
}
func Branch() fyne.CanvasObject {
	body := `A branch is a separate line of development in your project.

Simple definition:
A branch lets you work on new features or fixes without affecting the main code.

Why branches exist:
• isolate new work safely  
• experiment without breaking main  
• switch between tasks easily  
• collaborate without conflicts  

Common commands:

List all branches:
git branch or git branch --list

Create a new branch:
git branch feature-login

Switch to a branch:
git switch feature-login

Create and switch in one step:
git switch -c feature-login

Delete a branch:
git branch -d feature-login

Rename a branch:
git branch -m new-name

If you need to create a new branch and switch to it simultaneously, use the -b flag
git checkout -b <new_branch_name>

Running git checkout -b feature-branch is exactly the same as running these two commands back-to-back:

git branch feature-branch (Creates the branch)

git checkout feature-branch (Switches your working directory to that branch)
`
	content := FormatSection("Git Branch", body)
	return widget.NewRichText(content...)
}
func Pull() fyne.CanvasObject {
	body := `git pull brings the latest changes from a remote repository into your current branch.

Simple definition:
It downloads new commits from the remote and automatically merges them into your local branch.

Why pull is important:
• keeps your copy up to date  
• prevents conflicts during push  
• syncs your work with your team  

What git pull actually does:
It runs two operations:
1. git fetch  (download changes)
2. git merge  (combine them with your local branch)

Common usage:

Pull latest changes:
git pull

Pull from a specific branch:
git pull origin main

Pull without merging (rebase instead):
git pull --rebase`

	content := FormatSection("Git Pull", body)
	return widget.NewRichText(content...)
}
func Reflog() fyne.CanvasObject {
	body := `git reflog shows the full history of where your HEAD and branches have been.

Simple definition:
Reflog is a recovery log that records every movement of HEAD, even if commits were deleted or branches were changed.

Why reflog is important:
• recover lost commits  
• undo mistakes  
• restore deleted branches  
• go back to any previous state  

Think of it like a time machine for your repository.

Common usage:

Show full reflog:
git reflog

Reset to an earlier state:
git reset --hard <reflog-id>

Example:
If you messed up a rebase or deleted a commit, reflog helps you find the old commit and restore it.`

	content := FormatSection("Git Reflog", body)
	return widget.NewRichText(content...)
}
func GitIgnore() fyne.CanvasObject {
	body := `A .gitignore file tells Git which files and folders it should NOT track.

Simple definition:
.gitignore is a list of files Git must ignore.

Why .gitignore is important:
• keeps unwanted files out of commits  
• prevents IDE, OS, and build files from polluting the repo  
• protects sensitive files from being pushed  
• keeps the repo clean and lightweight  

Common examples:

Ignore OS files:
*.DS_Store
Thumbs.db

Ignore build folders:
build/
dist/
bin/

Ignore logs and temporary files:
*.log
*.tmp

Ignore environment files:
.env
config.local.json

Ignore IDE settings:
.vscode/
.idea/

Ignore all files of a type:
*.exe
*.dll
*.zip

Important note:
If a file is already tracked, .gitignore will NOT untrack it.
Use:
git rm --cached <file>`

	content := FormatSection(".gitignore", body)
	return widget.NewRichText(content...)
}
