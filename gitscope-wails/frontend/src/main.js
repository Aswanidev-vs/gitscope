import './style.css';

import {
    SelectRepo, OpenFolder, GetRepoPath, IsRepoInitialized,
    Init, Status, Stage, Commit, Push, Pull, Log, Revert, Clone,
    CreateBranch, DeleteBranch, SwitchBranch, BranchRename,
    Diff, Reset, Fetch, Stash, Merge, Tag, Remote, Reflog,
    Show, LsFiles, Clean, Shortlog, Blame, Worktree, Rebase,
    CherryPick, UndoLastCommit, MagicSync, GetConflicts,
    ResolveConflict, GetBranches, GetCurrentBranch,
    ReadGitIgnore, WriteGitIgnore, RunCommands, IsGitAvailable,
} from '../wailsjs/go/main/App';

let currentPage = 'dashboard';
let currentCategory = 'Common';
let consoleBuffer = '';
let docContents = {};

const ICONS = {
    repo: '📁', dashboard: '⚡', settings: 'ℹ️', docs: '📖',
    back: '←',
};

// ─── PAGES ──────────────────────────────────────────────────

const PAGES = {
    repo: renderRepoPage,
    dashboard: renderDashboardPage,
    settings: renderSettingsPage,
    docs: renderDocsPage,
};

const CATEGORIES = {
    Common: [
        { label: 'Init', action: () => runGitCmd(Init) },
        { label: 'Stage', action: () => runGitCmd(Stage, 'All (.)') },
        { label: 'Status', action: runStatus },
        { label: 'Commit', action: showCommitDialog },
        { label: 'Push', action: showPushDialog },
        { label: 'Log', action: runLog },
    ],
    Branches: [
        { label: 'Branch', action: showBranchDialog },
        { label: 'Switch', action: showSwitchDialog },
        { label: 'Merge', action: showMergeDialog },
        { label: 'Rename', action: showRenameDialog },
        { label: 'Tag', action: showTagDialog },
    ],
    Remote: [
        { label: 'Remote', action: showRemoteDialog },
        { label: 'Fetch', action: () => runGitCmd(Fetch, 'Default') },
        { label: 'Pull', action: showPullDialog },
        { label: 'Clone', action: showCloneDialog },
        { label: 'Cherry-pick', action: showCherryPickDialog },
    ],
    History: [
        { label: 'Log', action: runLog },
        { label: 'Revert', action: showRevertDialog },
        { label: 'Show', action: showShowDialog },
        { label: 'Shortlog', action: () => runGitCmd(Shortlog, 'Default') },
        { label: 'Reflog', action: () => runGitCmd(Reflog, 'Standard') },
    ],
    Changes: [
        { label: 'Diff', action: showDiffDialog },
        { label: 'Stash', action: showStashDialog },
        { label: 'Clean', action: showCleanDialog },
        { label: 'Ls-Files', action: () => runGitCmd(LsFiles, 'Tracked') },
        { label: '.gitignore', action: showGitIgnoreDialog },
    ],
    Advanced: [
        { label: 'Reset', action: showResetDialog },
        { label: 'Rebase', action: showRebaseDialog },
        { label: 'Undo', action: () => runGitCmd(UndoLastCommit) },
        { label: 'Worktree', action: showWorktreeDialog },
        { label: 'Conflicts', action: showConflictsDialog },
    ],
    Tools: [
        { label: 'Blame', action: showBlameDialog },
        { label: 'Magic Sync', action: () => runGitCmd(MagicSync) },
    ],
};

// ─── APP INIT ───────────────────────────────────────────────

document.addEventListener('DOMContentLoaded', async () => {
    renderApp();
    navigate('dashboard');
    updateRepoPath();
    setInterval(updateRepoPath, 3000);
});

function renderApp() {
    document.getElementById('app').innerHTML = `
        <div class="sidebar" id="sidebar">
            <button class="sidebar-btn active" data-page="repo" data-tooltip="Repository" onclick="navigate('repo')">${ICONS.repo}</button>
            <button class="sidebar-btn" data-page="dashboard" data-tooltip="Dashboard" onclick="navigate('dashboard')">${ICONS.dashboard}</button>
            <button class="sidebar-btn" data-page="settings" data-tooltip="About" onclick="navigate('settings')">${ICONS.settings}</button>
            <button class="sidebar-btn" data-page="docs" data-tooltip="Documentation" onclick="navigate('docs')">${ICONS.docs}</button>
            <div style="flex:1"></div>
        </div>
        <div class="main-area">
            <div class="top-bar">
                <h2>GitScope</h2>
                <span class="repo-path" id="repoPath">No repository selected</span>
                <span class="badge" id="branchBadge"></span>
            </div>
            <div class="content" id="contentArea"></div>
            <div class="console-panel">
                <div class="console-header">
                    <span>Console Output</span>
                    <button onclick="clearConsole()">Clear</button>
                </div>
                <div class="console-output" id="consoleOutput">Ready.</div>
            </div>
        </div>
    `;
    window.navigate = navigate;
    window.clearConsole = clearConsole;
}

// ─── NAVIGATION ─────────────────────────────────────────────

function navigate(page) {
    currentPage = page;
    document.querySelectorAll('.sidebar-btn').forEach(b => b.classList.toggle('active', b.dataset.page === page));
    const area = document.getElementById('contentArea');
    area.innerHTML = '';
    const renderer = PAGES[page];
    if (renderer) renderer(area);
}

// ─── CONSOLE ────────────────────────────────────────────────

function consoleLog(msg, type) {
    const el = document.getElementById('consoleOutput');
    if (!el) return;
    const prefix = type === 'error' ? '❌ ' : type === 'success' ? '✅ ' : type === 'warning' ? '⚠️ ' : '';
    const line = `${prefix}${msg}`;
    consoleBuffer += line + '\n';
    el.textContent = consoleBuffer;
    el.scrollTop = el.scrollHeight;
}

function clearConsole() {
    consoleBuffer = '';
    const el = document.getElementById('consoleOutput');
    if (el) el.textContent = '';
}

// ─── HELPERS ────────────────────────────────────────────────

async function updateRepoPath() {
    try {
        const path = await GetRepoPath();
        const el = document.getElementById('repoPath');
        if (el) el.textContent = path || 'No repository selected';
        if (path) {
            const branch = await GetCurrentBranch();
            const badge = document.getElementById('branchBadge');
            if (badge) badge.textContent = branch || '';
        }
    } catch (_) {}
}

async function checkRepo() {
    const path = await GetRepoPath();
    if (!path) {
        consoleLog('Please select a repository first.', 'warning');
        return false;
    }
    return true;
}

async function runGitCmd(fn, ...args) {
    if (!await checkRepo()) return;
    consoleLog(`Running ${fn.name}...`);
    try {
        const result = await fn(...args);
        consoleLog(result || `${fn.name} completed.`, 'success');
    } catch (err) {
        consoleLog(`${fn.name} failed: ${err}`, 'error');
    }
}

function showModal(html) {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.innerHTML = html;
    overlay.addEventListener('click', e => { if (e.target === overlay) overlay.remove(); });
    document.body.appendChild(overlay);
    return overlay;
}

// ─── REPO PAGE ──────────────────────────────────────────────

function renderRepoPage(area) {
    let commands = '';

    area.innerHTML = `
        <div class="page-area">
            <h3 class="section-title">Repository Setup</h3>
            <div class="repo-actions">
                <div style="display:flex;gap:8px;flex-wrap:wrap">
                    <button class="btn btn-primary" onclick="selectAndOpen()">📂 Select Repository</button>
                    <button class="btn btn-secondary" onclick="initRepo()">🏁 Initialize Git</button>
                </div>
                <div class="repo-card">
                    <h4>📝 Create New Repository</h4>
                    <p style="font-size:12px;color:var(--text-secondary);margin-bottom:8px">Paste GitHub commands to create and push a new repository</p>
                    <textarea id="newRepoCmds" placeholder="git init&#10;git add .&#10;git commit -m &quot;initial commit&quot;&#10;gh repo create ...">${commands}</textarea>
                    <button class="btn btn-primary" onclick="runNewRepoCmds()">Run Commands</button>
                </div>
                <div class="repo-card">
                    <h4>🔗 Push Existing Repository</h4>
                    <textarea id="existingRepoCmds" placeholder="git remote add origin https://github.com/yourname/repo.git&#10;git branch -M main&#10;git push -u origin main"></textarea>
                    <button class="btn btn-primary" onclick="runExistingRepoCmds()">Run Commands</button>
                </div>
            </div>
        </div>
    `;

    window.selectAndOpen = async () => {
        try {
            const path = await SelectRepo();
            if (path) {
                consoleLog(`Repository selected: ${path}`, 'success');
                updateRepoPath();
            }
        } catch (err) {
            consoleLog(`Error: ${err}`, 'error');
        }
    };

    window.initRepo = () => runGitCmd(Init);

    window.runNewRepoCmds = async () => {
        if (!await checkRepo()) return;
        const cmds = document.getElementById('newRepoCmds').value;
        consoleLog('Running setup commands...');
        try {
            const result = await RunCommands(cmds);
            consoleLog(result || 'Commands completed.', 'success');
        } catch (err) {
            consoleLog(`Error: ${err}`, 'error');
        }
    };

    window.runExistingRepoCmds = async () => {
        if (!await checkRepo()) return;
        const cmds = document.getElementById('existingRepoCmds').value;
        consoleLog('Running push commands...');
        try {
            const result = await RunCommands(cmds);
            consoleLog(result || 'Commands completed.', 'success');
        } catch (err) {
            consoleLog(`Error: ${err}`, 'error');
        }
    };
}

// ─── DASHBOARD PAGE ─────────────────────────────────────────

function renderDashboardPage(area) {
    const cats = Object.keys(CATEGORIES);
    area.innerHTML = `
        <div class="dashboard-layout">
            <div class="dashboard-left">
                <div style="font-size:11px;text-transform:uppercase;color:var(--text-secondary);padding:0 8px 4px;letter-spacing:0.5px">Categories</div>
                ${cats.map(c => `<button class="cat-btn ${c === currentCategory ? 'active' : ''}" onclick="switchCategory('${c}')">${c}</button>`).join('')}
            </div>
            <div class="dashboard-right" id="dashContent"></div>
        </div>
    `;

    window.switchCategory = (cat) => {
        currentCategory = cat;
        document.querySelectorAll('.cat-btn').forEach(b => b.classList.toggle('active', b.textContent === cat));
        renderCategoryButtons(document.getElementById('dashContent'), cat);
    };

    renderCategoryButtons(document.getElementById('dashContent'), currentCategory);
}

function renderCategoryButtons(container, category) {
    const btns = CATEGORIES[category] || [];
    const withSelect = ['Status', 'Log', 'Shortlog', 'Diff', 'Ls-Files'];
    container.innerHTML = `
        <h3 class="section-title">${category}</h3>
        <div class="btn-grid">
            ${btns.map(b => {
                if (withSelect.includes(b.label)) {
                    return `<button class="git-btn" onclick="window['on_${b.label}'] ? window['on_${b.label}']() : ${b.action.name}()">${b.label}</button>`;
                }
                return `<button class="git-btn" onclick="(${b.action.toString()})()">${b.label}</button>`;
            }).join('')}
        </div>
    `;
    btns.forEach(b => {
        if (withSelect.includes(b.label)) {
            window[`on_${b.label}`] = b.action;
        }
    });
}

// ─── DIALOG ACTIONS ─────────────────────────────────────────

function runStatus() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Status</h3>
            <div class="field">
                <label>Mode</label>
                <select id="statusMode">
                    <option value="Standard">Standard</option>
                    <option value="Short (-s)">Short (-s)</option>
                    <option value="Branch (-b)">Branch (-b)</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runStatusCmd()">Run</button>
            </div>
        </div>
    `);
    window.runStatusCmd = async () => {
        modal.remove();
        const opt = document.getElementById('statusMode').value;
        await runGitCmd(Status, opt);
    };
}

function runLog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Log</h3>
            <div class="field">
                <label>Format</label>
                <select id="logMode">
                    <option value="Oneline">Oneline</option>
                    <option value="Graph">Graph</option>
                    <option value="Pretty">Pretty</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runLogCmd()">Run</button>
            </div>
        </div>
    `);
    window.runLogCmd = async () => {
        modal.remove();
        const opt = document.getElementById('logMode').value;
        await runGitCmd(Log, opt);
    };
}

function showCommitDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Commit</h3>
            <div class="field">
                <label>Message</label>
                <input id="commitMsg" placeholder="Enter commit message" />
            </div>
            <div class="field">
                <label>Options</label>
                <select id="commitMode">
                    <option value="Standard (-m)">Standard (-m)</option>
                    <option value="Stage All (-a)">Stage All (-a)</option>
                    <option value="Amend (--amend)">Amend (--amend)</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runCommit()">Commit</button>
            </div>
        </div>
    `);
    window.runCommit = async () => {
        const msg = document.getElementById('commitMsg').value.trim();
        if (!msg) { consoleLog('Commit message cannot be empty', 'error'); return; }
        modal.remove();
        const opt = document.getElementById('commitMode').value;
        await runGitCmd(Commit, msg, opt);
    };
}

function showPushDialog() {
    (async () => {
        let branches = [];
        try { branches = await GetBranches(); } catch (_) {}
        const modal = showModal(`
            <div class="modal">
                <h3>Git Push</h3>
                <div class="field">
                    <label>Branch</label>
                    <select id="pushBranch">${branches.map(b => `<option>${b}</option>`).join('')}</select>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button class="btn btn-primary" onclick="runPush()">Push</button>
                </div>
            </div>
        `);
        window.runPush = async () => {
            modal.remove();
            const branch = document.getElementById('pushBranch').value;
            await runGitCmd(Push, branch);
        };
    })();
}

function showPullDialog() {
    (async () => {
        let branches = [];
        try { branches = await GetBranches(); } catch (_) {}
        const modal = showModal(`
            <div class="modal">
                <h3>Git Pull</h3>
                <div class="field">
                    <label>Branch</label>
                    <select id="pullBranch">${branches.map(b => `<option>${b}</option>`).join('')}</select>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button class="btn btn-primary" onclick="runPull()">Pull</button>
                </div>
            </div>
        `);
        window.runPull = async () => {
            modal.remove();
            const branch = document.getElementById('pullBranch').value;
            await runGitCmd(Pull, branch);
        };
    })();
}

function showBranchDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Branch Actions</h3>
            <div class="field">
                <label>Branch Name</label>
                <input id="branchName" placeholder="Enter branch name" />
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runBranchCreate()">Create</button>
                <button class="btn btn-danger" onclick="runBranchDelete()">Delete</button>
            </div>
        </div>
    `);
    window.runBranchCreate = async () => {
        const name = document.getElementById('branchName').value.trim();
        if (!name) return;
        modal.remove();
        await runGitCmd(CreateBranch, name);
    };
    window.runBranchDelete = async () => {
        const name = document.getElementById('branchName').value.trim();
        if (!name) return;
        modal.remove();
        await runGitCmd(DeleteBranch, name);
    };
}

function showSwitchDialog() {
    (async () => {
        let branches = [];
        try { branches = await GetBranches(); } catch (_) {}
        const modal = showModal(`
            <div class="modal">
                <h3>Switch Branch</h3>
                <div class="field">
                    <label>Branch</label>
                    <select id="switchBranch">${branches.map(b => `<option>${b}</option>`).join('')}</select>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button class="btn btn-primary" onclick="runSwitch()">Switch</button>
                </div>
            </div>
        `);
        window.runSwitch = async () => {
            modal.remove();
            const branch = document.getElementById('switchBranch').value;
            await runGitCmd(SwitchBranch, branch);
            updateRepoPath();
        };
    })();
}

function showMergeDialog() {
    (async () => {
        let branches = [];
        try { branches = await GetBranches(); } catch (_) {}
        const modal = showModal(`
            <div class="modal">
                <h3>Merge Branch</h3>
                <div class="field">
                    <label>Branch to merge</label>
                    <select id="mergeBranch">${branches.map(b => `<option>${b}</option>`).join('')}</select>
                </div>
                <div class="modal-actions">
                    <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                    <button class="btn btn-primary" onclick="runMerge()">Merge</button>
                </div>
            </div>
        `);
        window.runMerge = async () => {
            modal.remove();
            const branch = document.getElementById('mergeBranch').value;
            await runGitCmd(Merge, branch);
        };
    })();
}

function showRenameDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Rename Branch</h3>
            <div class="field"><label>Current Name</label><input id="oldName" placeholder="Current branch name" /></div>
            <div class="field"><label>New Name</label><input id="newName" placeholder="New branch name" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runRename()">Rename</button>
            </div>
        </div>
    `);
    window.runRename = async () => {
        const old = document.getElementById('oldName').value.trim();
        const n = document.getElementById('newName').value.trim();
        if (!old || !n) return;
        modal.remove();
        await runGitCmd(BranchRename, old, n);
    };
}

function showTagDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Tag Actions</h3>
            <div class="field">
                <label>Action</label>
                <select id="tagAction">
                    <option value="list">List</option>
                    <option value="create">Create</option>
                    <option value="delete">Delete</option>
                    <option value="push">Push</option>
                </select>
            </div>
            <div class="field"><label>Tag Name</label><input id="tagName" placeholder="Tag name" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runTag()">Run</button>
            </div>
        </div>
    `);
    window.runTag = async () => {
        const action = document.getElementById('tagAction').value;
        const name = document.getElementById('tagName').value.trim();
        modal.remove();
        if (action === 'list') {
            await runGitCmd(Tag, action, '');
        } else {
            if (!name) { consoleLog('Tag name required', 'error'); return; }
            await runGitCmd(Tag, action, name);
        }
    };
}

function showRemoteDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Remote Actions</h3>
            <div class="field">
                <label>Action</label>
                <select id="remoteAction">
                    <option value="list">List</option>
                    <option value="add">Add</option>
                    <option value="remove">Remove</option>
                </select>
            </div>
            <div class="field"><label>Remote URL / Name</label><input id="remoteUrl" placeholder="URL for add / name for remove" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runRemote()">Run</button>
            </div>
        </div>
    `);
    window.runRemote = async () => {
        const action = document.getElementById('remoteAction').value;
        const val = document.getElementById('remoteUrl').value.trim();
        modal.remove();
        if (action === 'list') {
            await runGitCmd(Remote, 'list', '');
        } else if (action === 'remove') {
            if (!val) { consoleLog('Remote name required', 'error'); return; }
            await runGitCmd(Remote, 'remove', val);
        } else {
            await runGitCmd(Remote, 'add', val);
        }
    };
}

function showRevertDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Revert Commit</h3>
            <div class="field"><label>Commit Hash</label><input id="revertHash" placeholder="e.g. a1s4fd6" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runRevert()">Revert</button>
            </div>
        </div>
    `);
    window.runRevert = async () => {
        const hash = document.getElementById('revertHash').value.trim();
        if (!hash) { consoleLog('Commit hash required', 'error'); return; }
        modal.remove();
        await runGitCmd(Revert, hash);
    };
}

function showCloneDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Clone Repository</h3>
            <div class="field"><label>Repository URL</label><input id="cloneUrl" placeholder="https://github.com/yourname/repo.git" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runClone()">Clone</button>
            </div>
        </div>
    `);
    window.runClone = async () => {
        const url = document.getElementById('cloneUrl').value.trim();
        if (!url) { consoleLog('Repository URL required', 'error'); return; }
        modal.remove();
        await runGitCmd(Clone, url);
    };
}

function showDiffDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Diff</h3>
            <div class="field">
                <label>Mode</label>
                <select id="diffMode">
                    <option value="Unstaged">Unstaged</option>
                    <option value="Staged (--cached)">Staged (--cached)</option>
                    <option value="Names (--name-only)">Names (--name-only)</option>
                    <option value="Summary (--stat)">Summary (--stat)</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runDiff()">Run</button>
            </div>
        </div>
    `);
    window.runDiff = async () => {
        modal.remove();
        const opt = document.getElementById('diffMode').value;
        await runGitCmd(Diff, opt);
    };
}

function showResetDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Reset</h3>
            <div class="field">
                <label>Mode</label>
                <select id="resetMode">
                    <option value="--mixed">--mixed</option>
                    <option value="--soft">--soft</option>
                    <option value="--hard">--hard</option>
                </select>
            </div>
            <div class="field"><label>Target</label><input id="resetTarget" value="HEAD~1" placeholder="HEAD~1 or hash" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runReset()">Reset</button>
            </div>
        </div>
    `);
    window.runReset = async () => {
        modal.remove();
        const mode = document.getElementById('resetMode').value;
        const target = document.getElementById('resetTarget').value.trim();
        await runGitCmd(Reset, mode, target);
    };
}

function showStashDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Stash</h3>
            <div class="field">
                <label>Action</label>
                <select id="stashAction">
                    <option value="Save">Save</option>
                    <option value="Pop">Pop</option>
                    <option value="List">List</option>
                    <option value="Drop">Drop</option>
                    <option value="Apply">Apply</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runStash()">Run</button>
            </div>
        </div>
    `);
    window.runStash = async () => {
        modal.remove();
        const action = document.getElementById('stashAction').value;
        await runGitCmd(Stash, action);
    };
}

function showCleanDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Clean</h3>
            <div class="field">
                <label>Mode</label>
                <select id="cleanMode">
                    <option value="Preview (-n)">Preview (-n)</option>
                    <option value="Remove Dir (-d)">Remove Dir (-d)</option>
                    <option value="Force (-f)">Force (-f)</option>
                    <option value="Full (-fdx)">Full (-fdx)</option>
                </select>
            </div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runClean()">Run</button>
            </div>
        </div>
    `);
    window.runClean = async () => {
        modal.remove();
        const opt = document.getElementById('cleanMode').value;
        await runGitCmd(Clean, opt);
    };
}

function showRebaseDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Rebase</h3>
            <div class="field">
                <label>Action</label>
                <select id="rebaseAction">
                    <option value="Interactive (-i)">Interactive (-i)</option>
                    <option value="Onto">Onto</option>
                    <option value="Continue">Continue</option>
                    <option value="Abort">Abort</option>
                    <option value="Skip">Skip</option>
                </select>
            </div>
            <div class="field"><label>Target</label><input id="rebaseTarget" placeholder="target branch or commit" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runRebase()">Run</button>
            </div>
        </div>
    `);
    window.runRebase = async () => {
        modal.remove();
        const action = document.getElementById('rebaseAction').value;
        const target = document.getElementById('rebaseTarget').value.trim();
        await runGitCmd(Rebase, action, target);
    };
}

function showCherryPickDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Cherry-pick</h3>
            <div class="field"><label>Commit Hash</label><input id="cherryHash" placeholder="Enter commit hash" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runCherryPick()">Apply</button>
            </div>
        </div>
    `);
    window.runCherryPick = async () => {
        const hash = document.getElementById('cherryHash').value.trim();
        if (!hash) return;
        modal.remove();
        await runGitCmd(CherryPick, hash);
    };
}

function showShowDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Show</h3>
            <div class="field">
                <label>Mode</label>
                <select id="showMode">
                    <option value="Head">Head</option>
                    <option value="Last 5">Last 5</option>
                    <option value="Specific">Specific</option>
                </select>
            </div>
            <div class="field"><label>Commit Hash (for Specific)</label><input id="showHash" placeholder="Commit hash" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runShow()">Show</button>
            </div>
        </div>
    `);
    window.runShow = async () => {
        modal.remove();
        const mode = document.getElementById('showMode').value;
        const hash = document.getElementById('showHash').value.trim();
        const target = mode === 'Specific' ? hash : '';
        await runGitCmd(Show, mode, target);
    };
}

function showBlameDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Git Blame</h3>
            <div class="field"><label>File Path</label><input id="blameFile" placeholder="path/to/file.go" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runBlame()">Blame</button>
            </div>
        </div>
    `);
    window.runBlame = async () => {
        const file = document.getElementById('blameFile').value.trim();
        if (!file) return;
        modal.remove();
        await runGitCmd(Blame, file);
    };
}

function showWorktreeDialog() {
    const modal = showModal(`
        <div class="modal">
            <h3>Worktree</h3>
            <div class="field">
                <label>Action</label>
                <select id="wtAction">
                    <option value="List">List</option>
                    <option value="Add">Add</option>
                    <option value="Remove">Remove</option>
                    <option value="Prune">Prune</option>
                </select>
            </div>
            <div class="field"><label>Args (for Add/Remove)</label><input id="wtArgs" placeholder="path branch" /></div>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="runWt()">Run</button>
            </div>
        </div>
    `);
    window.runWt = async () => {
        modal.remove();
        const action = document.getElementById('wtAction').value;
        const args = document.getElementById('wtArgs').value.trim();
        await runGitCmd(Worktree, action, args);
    };
}

function showConflictsDialog() {
    (async () => {
        try {
            const conflicts = await GetConflicts();
            if (!conflicts || conflicts.length === 0) {
                consoleLog('No merge conflicts detected.', 'success');
                return;
            }
            let html = `<div class="modal"><h3>Resolve Conflicts</h3>`;
            conflicts.forEach(f => {
                html += `<div class="field" style="display:flex;gap:8px;align-items:center">
                    <span style="flex:1">${f}</span>
                    <button class="btn btn-primary" style="padding:4px 12px;font-size:11px" onclick="resolveConflict('${f}','ours')">Keep Mine</button>
                    <button class="btn btn-danger" style="padding:4px 12px;font-size:11px" onclick="resolveConflict('${f}','theirs')">Take Theirs</button>
                </div>`;
            });
            html += `<div class="modal-actions"><button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Close</button></div></div>`;
            const modal = showModal(html);
            window.resolveConflict = async (file, strategy) => {
                modal.remove();
                await runGitCmd(ResolveConflict, file, strategy);
            };
        } catch (err) {
            consoleLog(`Error: ${err}`, 'error');
        }
    })();
}

function showGitIgnoreDialog() {
    (async () => {
        try {
            const content = await ReadGitIgnore() || '';
            const modal = showModal(`
                <div class="modal" style="min-width:500px">
                    <h3>.gitignore</h3>
                    <textarea id="gitignoreContent" style="width:100%;min-height:200px;font-family:monospace;font-size:12px">${content}</textarea>
                    <div class="modal-actions">
                        <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                        <button class="btn btn-primary" onclick="saveGitIgnore()">Save</button>
                    </div>
                </div>
            `);
            window.saveGitIgnore = async () => {
                const text = document.getElementById('gitignoreContent').value;
                try {
                    await WriteGitIgnore(text);
                    consoleLog('.gitignore saved.', 'success');
                    modal.remove();
                } catch (err) {
                    consoleLog(`Error: ${err}`, 'error');
                }
            };
        } catch (err) {
            consoleLog(`Error: ${err}`, 'error');
        }
    })();
}

// ─── SETTINGS / ABOUT PAGE ─────────────────────────────────

function renderSettingsPage(area) {
    area.innerHTML = `
        <div class="page-area about-page">
            <h2>🔮 GitScope</h2>
            <p>Modern, lightweight, and visually intuitive Git client</p>
            <p>Version: 1.0.0</p>
            <p style="color:var(--text-secondary)">Built with Go + Wails</p>
            <p style="margin-top:8px">
                <a href="https://github.com/Aswanidev-vs/GitScope" target="_blank">🔗 View Project on GitHub</a>
            </p>
        </div>
    `;
}

// ─── DOCS PAGE ──────────────────────────────────────────────

const DOC_ITEMS = [
    "Init","Stage","Status","Commit","Push","Log","Revert","Clone","Branch","Pull",
    "Reflog","GitIgnore","Remote","Diff","Reset","Fetch","Stash","Merge","Tag",
    "Cherry-pick","Rebase","Clean","Show","Ls-files","Worktree","Shortlog","Blame"
];

const DOC_TEXT = {
    Init: `Git Init\n\n${'='.repeat(40)}\n\ngit init creates a new empty Git repository in your current folder.\n\nThink of it as telling Git:\n"Start a new project here and track everything I do."\n\nCommand:\n  git init\n\nAfter running this, Git creates a hidden folder:\n  .git\n\nExample:\n  mkdir myproject\n  cd myproject\n  git init`,
    Stage: `Git Stage\n\n${'='.repeat(40)}\n\nThe stage (staging area/index) is where Git stores changes\nyou want to include in your next commit.\n\nSimple definition:\n  The stage is where you put files that Git will commit.\n\nWhy staging exists:\n  - choose which files to commit\n  - commit only specific parts\n  - prepare commits cleanly\n\nCommands:\n  git add file.txt    (stage a single file)\n  git add .           (stage all files)`,
    Status: `Git Status\n\n${'='.repeat(40)}\n\ngit status shows the current state of your working directory\nand staging area.\n\nIt tells you:\n  - what changed\n  - what is staged\n  - what is unstaged\n  - which files are untracked\n\nCommand:\n  git status\n\nExample output:\n  Changes not staged for commit:\n    modified:   main.go\n  Untracked files:\n    newfile.txt`,
    Commit: `Git Commit\n\n${'='.repeat(40)}\n\nA commit is a snapshot of your project at a specific point in time.\n\nCommands:\n  git commit -m "message"\n  git commit -a -m "message"   (stage + commit)\n  git commit --amend            (amend last commit)`,
    Push: `Git Push\n\n${'='.repeat(40)}\n\nPush sends your local commits to a remote repository.\n\nCommands:\n  git push\n  git push origin branchname\n  git push -u origin branchname\n  git push --force`,
    Log: `Git Log\n\n${'='.repeat(40)}\n\nLog shows the full history of commits in your repository.\n\nCommands:\n  git log\n  git log --oneline\n  git log --graph --oneline --decorate --all\n  git log --author="name"`,
    Revert: `Git Revert\n\n${'='.repeat(40)}\n\nRevert undoes a specific commit by creating a new commit.\n\nCommands:\n  git revert <commit-id>\n  git revert --no-commit <old>..<new>`,
    Clone: `Git Clone\n\n${'='.repeat(40)}\n\nClone creates a local copy of a remote repository.\n\nCommands:\n  git clone <repo-url>\n  git clone <repo-url> myproject\n  git clone --depth 1 <repo-url>`,
    Branch: `Git Branch\n\n${'='.repeat(40)}\n\nA branch is a separate line of development.\n\nCommands:\n  git branch                        (list)\n  git branch feature-login          (create)\n  git switch feature-login          (switch)\n  git switch -c feature-login       (create + switch)\n  git branch -d feature-login       (delete)\n  git branch -m new-name            (rename)`,
    Pull: `Git Pull\n\n${'='.repeat(40)}\n\ngit pull brings changes from remote to your current branch.\n\nActually does:\n  1. git fetch  (download changes)\n  2. git merge  (combine with local)\n\nCommands:\n  git pull\n  git pull origin main\n  git pull --rebase`,
    Reflog: `Git Reflog\n\n${'='.repeat(40)}\n\nReflog shows the history of where HEAD has been.\n\nThink of it as a time machine for your repository.\n\nCommands:\n  git reflog\n  git reset --hard <reflog-id>`,
    GitIgnore: `.gitignore\n\n${'='.repeat(40)}\n\nTells Git which files/folders to NOT track.\n\nCommon entries:\n  *.DS_Store\n  build/\n  dist/\n  *.log\n  .env\n  .vscode/\n  *.exe`,
    Remote: `Remote\n\n${'='.repeat(40)}\n\ngit remote -v shows all connected remote URLs.\n\nCommands:\n  git remote -v\n  git remote add origin <url>\n  git remote remove origin`,
    Diff: `Git Diff\n\n${'='.repeat(40)}\n\nShows differences between file versions.\n\nCommands:\n  git diff                    (unstaged)\n  git diff --staged          (staged)\n  git diff HEAD              (all changes)\n  git diff main feature      (between branches)`,
    Reset: `Git Reset\n\n${'='.repeat(40)}\n\nMoves HEAD and controls commit history/staging/working dir.\n\nModes:\n  --soft   → undo commit, keep staged\n  --mixed  → undo commit, unstage (default)\n  --hard   → delete everything (DANGEROUS)\n\nCommands:\n  git reset HEAD^\n  git reset --soft HEAD^\n  git reset --hard HEAD^`,
    Fetch: `Git Fetch\n\n${'='.repeat(40)}\n\nDownloads changes from remote without merging.\n\nCommands:\n  git fetch origin\n  git fetch --all`,
    Stash: `Git Stash\n\n${'='.repeat(40)}\n\nTemporarily shelves changes to get a clean directory.\n\nCommands:\n  git stash\n  git stash list\n  git stash pop\n  git stash apply`,
    Merge: `Git Merge\n\n${'='.repeat(40)}\n\nJoins two development histories together.\n\nCommands:\n  git merge feature-x\n  git merge --no-ff feature-x`,
    Tag: `Git Tag\n\n${'='.repeat(40)}\n\nMarks specific points as important (releases).\n\nCommands:\n  git tag v1.0\n  git push origin v1.0`,
    "Cherry-pick": `Cherry-pick\n\n${'='.repeat(40)}\n\nApplies changes from existing commits to current branch.\n\nCommands:\n  git cherry-pick <commit-hash>`,
    Rebase: `Git Rebase\n\n${'='.repeat(40)}\n\nReapplies commits on top of another base tip.\n\nCommands:\n  git rebase main\n  git rebase -i HEAD~3\n  git rebase --continue\n  git rebase --abort`,
    Clean: `Git Clean\n\n${'='.repeat(40)}\n\nRemoves untracked files.\n\nCommands:\n  git clean -n    (preview)\n  git clean -f    (remove files)\n  git clean -fd   (remove files & dirs)`,
    Show: `Git Show\n\n${'='.repeat(40)}\n\nShows details about a Git object (commit, tag, etc.).\n\nCommands:\n  git show HEAD\n  git show <hash>\n  git show --stat`,
    "Ls-files": `Git Ls-files\n\n${'='.repeat(40)}\n\nShows information about files in the index and working tree.\n\nCommands:\n  git ls-files\n  git ls-files --cached\n  git ls-files --others`,
    Worktree: `Git Worktree\n\n${'='.repeat(40)}\n\nManage multiple working trees attached to one repo.\n\nCommands:\n  git worktree list\n  git worktree add <path> <branch>\n  git worktree remove <name>`,
    Shortlog: `Git Shortlog\n\n${'='.repeat(40)}\n\nSummarizes git log output grouped by author.\n\nCommands:\n  git shortlog\n  git shortlog -s\n  git shortlog -n`,
    Blame: `Git Blame\n\n${'='.repeat(40)}\n\nShows what revision/author last modified each line.\n\nCommands:\n  git blame <file>\n  git blame -L 10,20 <file>`,
};

function renderDocsPage(area) {
    area.innerHTML = `
        <div class="page-area">
            <h3 class="section-title">Documentation</h3>
            <div class="doc-list" id="docList"></div>
            <div class="doc-content" id="docContent">Select a topic from above to view documentation.</div>
        </div>
    `;

    const list = document.getElementById('docList');
    DOC_ITEMS.forEach(item => {
        const btn = document.createElement('button');
        btn.className = 'doc-btn';
        btn.textContent = item;
        btn.onclick = () => {
            document.querySelectorAll('.doc-btn').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            const content = DOC_TEXT[item] || 'No documentation available.';
            document.getElementById('docContent').textContent = content;
        };
        list.appendChild(btn);
    });
}
