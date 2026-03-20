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
func Remote() fyne.CanvasObject {
	body := `
git remote -v shows all the remote repository URLs that your local project is connected to.

remote → means remote repository

-v → means "verbose" (show the URLs)

🟦 Simple Real Example

Suppose you cloned a repo:

git clone https://github.com/someone/project.git


Now run:

git remote -v


You will see:

origin  https://github.com/someone/project.git (fetch)
origin  https://github.com/someone/project.git (push)

This means:

Your local project is linked to one remote named origin

Both fetching and pushing use the same URL

🟩 Example With Multiple Remotes

If you add another remote:

git remote add backup https://github.com/you/backup-project.git


Now run:

git remote -v


Output becomes:

origin  https://github.com/someone/project.git (fetch)
origin  https://github.com/someone/project.git (push)
backup  https://github.com/you/backup-project.git (fetch)
backup  https://github.com/you/backup-project.git (push)

Meaning:

origin → main repo

backup → second repo

🟥 Why git remote -v matters

It helps you check:

✔ Which repo you are going to push to
✔ Which repo you are going to pull from
✔ If you cloned a repo and want to change the remote

🟨 Quick Use Case: Changing Remote

If this shows:

origin  https://github.com/someone/project.git (fetch)


but you want to push to your own repo:

git remote remove origin
git remote add origin https://github.com/you/myrepo.git


	`
	content := FormatSection("Remote", body)
	return widget.NewRichText(content...)
}
func Diff() fyne.CanvasObject {
	body := `
git diff shows the difference between changes in your files.  
It helps you see what you changed before committing.

diff → difference between versions of files

🟦 Simple Meaning

Think of git diff as:

"What exactly did I change?"

It compares:
• Old version vs new version
• Saved files vs last commit
• Staged vs unstaged changes

🟦 Simple Real Example

Suppose you have a file app.go:

Before:
fmt.Println("Hello")

You change it to:
fmt.Println("Hello World")

Now run:

git diff

Output:
- fmt.Println("Hello")
+ fmt.Println("Hello World")

Meaning:
- line was removed
+ line was added

🟩 Understanding Symbols

+  added line  
-  removed line  

No symbol means unchanged context.

🟦 git diff (most common)

Command:
git diff

Shows:
Changes that are NOT staged yet

Use case:
You edited files but did not run git add

🟩 Example

You edit main.go but do not add it.

git diff

Shows what you changed in main.go.

🟦 git diff --staged (or --cached)

Command:
git diff --staged

Shows:
Changes that ARE staged and ready to commit

Use case:
You already ran git add and want to review before commit.

🟩 Example

git add main.go
git diff --staged

Shows what will go into the next commit.

🟦 git diff HEAD

Command:
git diff HEAD

Shows:
All changes compared to last commit
Includes staged + unstaged changes.

🟦 git diff filename

Command:
git diff main.go

Shows:
Differences only for that file.

Helpful when:
You changed many files but want to check one.

🟦 git diff branch1 branch2

Command:
git diff main feature-login

Shows:
Difference between two branches.

Use case:
Before merging a feature branch.

🟦 git diff --stat

Command:
git diff --stat

Shows:
Summary instead of full code.

Example output:
 main.go | 5 +++--
 utils.go | 2 ++

Meaning:
• main.go changed 5 lines
• utils.go changed 2 lines

🟥 Why git diff is important

It helps you:
✔ Avoid committing mistakes
✔ Review changes clearly
✔ Understand what broke your code
✔ Learn what exactly you modified

🟨 Best Beginner Workflow

1. Edit files
2. Run git diff
3. Run git add .
4. Run git diff --staged
5. git commit -m "message"

This keeps your commits clean and safe.
`
	content := FormatSection("Diff", body)
	return widget.NewRichText(content...)
}
func Reset() fyne.CanvasObject {
	body := `
git reset is used to move HEAD and control what happens to:
• commit history
• staging area (index)
• working directory (files)

reset → go back to a previous state

🟦 Simple Meaning

Think of git reset as:

"I want to undo something, but in different levels."

Git reset has THREE common modes:
• --soft
• --mixed (default)
• --hard

Each one affects different areas.

🟦 Areas to remember (VERY IMPORTANT)

1. HEAD → last commit pointer
2. Staging area → git add
3. Working directory → actual files

🟦 git reset --soft (least destructive)

Command:
git reset --soft HEAD^

What it does:
✔ Moves HEAD back
✔ Keeps changes staged
✔ Files are NOT changed

Use case:
You want to redo the last commit message or combine commits.

🟩 Example

You committed too early:
git commit -m "oops"

Undo commit but keep everything staged:
git reset --soft HEAD^

Now you can recommit:
git commit -m "correct message"

🟦 git reset (or git reset --mixed) MOST COMMON

Command:
git reset HEAD^

This is the default mode.

What it does:
✔ Moves HEAD back
✔ Unstages files
✔ Keeps file changes

Use case:
You added files but do not want them staged yet.

🟩 Example

git add .
git commit -m "wrong commit"

Undo commit and unstage changes:
git reset HEAD^

Now files are edited but not staged.

🟦 git reset --hard (DANGEROUS)

Command:
git reset --hard HEAD^

What it does:
✔ Moves HEAD back
✔ Clears staging area
✔ Deletes file changes permanently

⚠ Warning:
This will DESTROY uncommitted changes.

🟩 Example

You completely messed up:
git reset --hard HEAD^

Your project becomes exactly like the previous commit.

🟥 NEVER use --hard unless you are 100% sure.

🟦 git reset filename (very common)

Command:
git reset main.go

What it does:
✔ Unstages a file
✔ Keeps changes in file

Use case:
You accidentally added a file.

🟩 Example

git add main.go
git reset main.go

main.go is no longer staged.

🟦 git reset HEAD (quick unstaging)

Command:
git reset HEAD

What it does:
✔ Unstages all files
✔ Keeps changes

Equivalent to:
git reset --mixed HEAD

🟦 Why git reset is important

It helps you:
✔ Fix wrong commits
✔ Unstage files safely
✔ Clean your commit history
✔ Recover from small mistakes

🟨 Beginner Safe Rule

Use these confidently:
git reset
git reset filename
git reset --soft HEAD^

Avoid unless sure:
git reset --hard

🟢 One-line memory trick

soft  → undo commit, keep staged  
mixed → undo commit, unstage  
hard  → delete everything  

`
	content := FormatSection("Reset", body)
	return widget.NewRichText(content...)
}

func Fetch() fyne.CanvasObject {
	body := `git fetch downloads commits, files, and refs from a remote repository into your local repo.

Simple definition:
Fetch gets updates from the remote but does NOT merge them into your local work.

Why fetch is useful:
• see what others have worked on without changing your files
• check for branch updates safely
• update remote-tracking branches

Common commands:
git fetch origin
git fetch --all`

	content := FormatSection("Git Fetch", body)
	return widget.NewRichText(content...)
}

func Stash() fyne.CanvasObject {
	body := `git stash temporarily shelves (or stashes) changes you've made to your working copy so you can work on something else.

Simple definition:
Stash is a "drawer" where you put unfinished work to get a clean directory.

Common usage:
git stash          (save changes)
git stash list     (see saved stashes)
git stash pop      (apply and remove latest)
git stash apply    (apply and keep in stash)`

	content := FormatSection("Git Stash", body)
	return widget.NewRichText(content...)
}

func Merge() fyne.CanvasObject {
	body := `git merge joins two or more development histories together.

Simple definition:
Merge combines changes from one branch into another (usually into main).

Example:
git branch feature-x
# ... work ...
git switch main
git merge feature-x

If changes conflict, Git will ask you to resolve them manually.`

	content := FormatSection("Git Merge", body)
	return widget.NewRichText(content...)
}

func Tag() fyne.CanvasObject {
	body := `Tags are used to mark specific points in history as being important, typically for releases (v1.0, v2.0).

Simple definition:
A tag is a permanent label for a specific commit.

Commands:
git tag v1.0
git push origin v1.0`

	content := FormatSection("Git Tag", body)
	return widget.NewRichText(content...)
}

func CherryPick() fyne.CanvasObject {
	body := `git cherry-pick applies the changes introduced by some existing commits to your current branch.

Simple definition:
"Copy and paste" a single commit from one branch to another.

Command:
git cherry-pick <commit-hash>`

	content := FormatSection("Git Cherry-pick", body)
	return widget.NewRichText(content...)
}
