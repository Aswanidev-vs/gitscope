import './style.css';
import { icon } from './icons.js';

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

const NAV_ITEMS = [
    { id: 'repo',      icon: 'repo',      tip: 'Repository' },
    { id: 'dashboard',  icon: 'dashboard',  tip: 'Dashboard' },
    { id: 'about',      icon: 'about',      tip: 'About' },
    { id: 'docs',       icon: 'docs',       tip: 'Docs' },
];

const CATEGORIES = {
    Common: [
        { icon: 'init',     label: 'Init',     action: () => runGitCmd(Init) },
        { icon: 'stage',    label: 'Stage',    action: () => runGitCmd(Stage, 'All (.)') },
        { icon: 'status',   label: 'Status',   action: runStatus },
        { icon: 'commit',   label: 'Commit',   action: showCommitDialog },
        { icon: 'push',     label: 'Push',     action: showPushDialog },
        { icon: 'log',      label: 'Log',      action: runLog },
    ],
    Branches: [
        { icon: 'branch',   label: 'Branch',   action: showBranchDialog },
        { icon: 'branch',   label: 'Switch',   action: showSwitchDialog },
        { icon: 'merge',    label: 'Merge',    action: showMergeDialog },
        { icon: 'branch',   label: 'Rename',   action: showRenameDialog },
        { icon: 'tag',      label: 'Tag',      action: showTagDialog },
    ],
    Remote: [
        { icon: 'remote',   label: 'Remote',   action: showRemoteDialog },
        { icon: 'fetch',    label: 'Fetch',    action: () => runGitCmd(Fetch, 'Default') },
        { icon: 'pull',     label: 'Pull',     action: showPullDialog },
        { icon: 'clone',    label: 'Clone',    action: showCloneDialog },
        { icon: 'cherry',   label: 'Cherry-pick', action: showCherryPickDialog },
    ],
    History: [
        { icon: 'log',      label: 'Log',      action: runLog },
        { icon: 'revert',   label: 'Revert',   action: showRevertDialog },
        { icon: 'show',     label: 'Show',     action: showShowDialog },
        { icon: 'shortlog', label: 'Shortlog', action: () => runGitCmd(Shortlog, 'Default') },
        { icon: 'log',      label: 'Reflog',   action: () => runGitCmd(Reflog, 'Standard') },
    ],
    Changes: [
        { icon: 'diff',     label: 'Diff',     action: showDiffDialog },
        { icon: 'stash',    label: 'Stash',    action: showStashDialog },
        { icon: 'clean',    label: 'Clean',    action: showCleanDialog },
        { icon: 'lsfiles',  label: 'Ls-Files', action: () => runGitCmd(LsFiles, 'Tracked') },
        { icon: 'gitignore',label: '.gitignore', action: showGitIgnoreDialog },
    ],
    Advanced: [
        { icon: 'reset',    label: 'Reset',    action: showResetDialog },
        { icon: 'rebase',   label: 'Rebase',   action: showRebaseDialog },
        { icon: 'undo',     label: 'Undo',     action: () => runGitCmd(UndoLastCommit) },
        { icon: 'worktree', label: 'Worktree', action: showWorktreeDialog },
        { icon: 'conflict', label: 'Conflicts', action: showConflictsDialog },
    ],
    Tools: [
        { icon: 'blame',    label: 'Blame',    action: showBlameDialog },
        { icon: 'sync',     label: 'Magic Sync', action: () => runGitCmd(MagicSync) },
    ],
};

// ─── APP INIT ───────────────────────────────────

document.addEventListener('DOMContentLoaded', () => {
    renderShell();
    navigate('dashboard');
    updateRepoInfo();
    setInterval(updateRepoInfo, 3000);
});

function renderShell() {
    document.getElementById('app').innerHTML = `
        <div class="sidebar" id="sidebar">
            <div class="sidebar-icon-group">
                ${NAV_ITEMS.map(n => `
                    <button class="sidebar-btn${n.id === 'repo' ? ' active' : ''}"
                            data-page="${n.id}" data-tip="${n.tip}"
                            onclick="navigate('${n.id}')">${icon(n.icon)}</button>
                `).join('')}
            </div>
        </div>
        <div class="main-area">
            <div class="top-bar">
                <div class="top-bar-logo">${icon('git', 20)}<h1>GitScope</h1></div>
                <div class="top-bar-sep"></div>
                <span class="repo-path" id="repoPath">No repository selected</span>
                <div class="branch-badge" id="branchBadge"></div>
            </div>
            <div class="content" id="contentArea"></div>
            <div class="console-panel" id="consolePanel">
                <div class="console-resize-handle" id="consoleResize"></div>
                <div class="console-header">
                    <div class="console-header-label">${icon('code', 12)}<span>Console</span></div>
                    <button class="btn btn-ghost btn-sm" onclick="clearConsole()">${icon('clear', 12)} Clear</button>
                </div>
                <div class="console-output" id="consoleOutput">Ready.</div>
                <div class="console-status">
                    <div class="console-status-dot" id="consoleDot"></div>
                    <span id="consoleStatusText">Idle</span>
                </div>
            </div>
        </div>
        <div class="toast-container" id="toastContainer"></div>
    `;
    window.navigate = navigate;
    window.clearConsole = clearConsole;
    initConsoleResize();
}

// ─── NAVIGATION ─────────────────────────────────

function navigate(page) {
    currentPage = page;
    document.querySelectorAll('.sidebar-btn').forEach(b =>
        b.classList.toggle('active', b.dataset.page === page));
    const area = document.getElementById('contentArea');
    area.innerHTML = '';
    const renderers = {
        repo: renderRepoPage,
        dashboard: renderDashboardPage,
        about: renderAboutPage,
        docs: renderDocsPage,
    };
    if (renderers[page]) renderers[page](area);
}

// ─── CONSOLE ────────────────────────────────────

let consoleActive = false;

function consoleLog(msg, type) {
    const prefixes = { error: '[err] ', success: '[ok] ', warning: '[warn] ' };
    const line = (prefixes[type] || '') + msg;
    consoleBuffer += line + '\n';
    const el = document.getElementById('consoleOutput');
    if (el) {
        el.textContent = consoleBuffer;
        el.scrollTop = el.scrollHeight;
    }
    window.console && window.console.log('[GitScope]', line);
    try { showToast(msg, type || 'info'); } catch (_) {}
}

function clearConsole() {
    consoleBuffer = '';
    const el = document.getElementById('consoleOutput');
    if (el) el.textContent = '';
}

function setConsoleBusy(busy) {
    consoleActive = busy;
    const dot = document.getElementById('consoleDot');
    const txt = document.getElementById('consoleStatusText');
    if (dot) dot.classList.toggle('active', busy);
    if (txt) txt.textContent = busy ? 'Running...' : 'Idle';
}

// ─── CONSOLE RESIZE ─────────────────────────────

function initConsoleResize() {
    const handle = document.getElementById('consoleResize');
    const panel = document.getElementById('consolePanel');
    if (!handle || !panel) return;

    let startY, startH, dragging = false;

    const onMouseMove = (e) => {
        if (!dragging) return;
        const delta = startY - e.clientY;
        const newH = Math.max(60, Math.min(window.innerHeight * 0.5, startH + delta));
        panel.style.height = newH + 'px';
        panel.style.setProperty('--console-height', newH + 'px');
    };

    const onMouseUp = () => {
        dragging = false;
        handle.classList.remove('active');
        document.removeEventListener('mousemove', onMouseMove);
        document.removeEventListener('mouseup', onMouseUp);
        document.body.style.cursor = '';
        document.body.style.userSelect = '';
    };

    handle.addEventListener('mousedown', (e) => {
        e.preventDefault();
        dragging = true;
        startY = e.clientY;
        startH = panel.offsetHeight;
        handle.classList.add('active');
        document.body.style.cursor = 'ns-resize';
        document.body.style.userSelect = 'none';
        document.addEventListener('mousemove', onMouseMove);
        document.addEventListener('mouseup', onMouseUp);
    });
}

// ─── TOAST ──────────────────────────────────────

function showToast(msg, type) {
    const container = document.getElementById('toastContainer');
    if (!container) return;
    const toastIconMap = { error: 'xCircle', success: 'check', warning: 'warn', info: 'code' };
    const toast = document.createElement('div');
    toast.className = `toast ${type || 'info'}`;
    toast.innerHTML = `${icon(toastIconMap[type] || 'code', 14)}<span>${msg}</span>`;
    container.appendChild(toast);
    setTimeout(() => toast.remove(), 3500);
}

// ─── HELPERS ────────────────────────────────────

async function updateRepoInfo() {
    try {
        const path = await GetRepoPath();
        const el = document.getElementById('repoPath');
        if (el) el.textContent = path || 'No repository selected';
        const badge = document.getElementById('branchBadge');
        if (path) {
            const branch = await GetCurrentBranch();
            badge.innerHTML = icon('branch', 12) + ' ' + (branch || '');
            badge.style.display = branch ? 'inline-flex' : 'none';
        } else {
            badge.style.display = 'none';
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
    setConsoleBusy(true);
    try {
        const result = await fn(...args);
        const output = (typeof result === 'string' ? result.trim() : '') || `${fn.name} completed.`;
        consoleLog(output, 'success');
    } catch (err) {
        const errMsg = (err && err.message) ? err.message : String(err);
        consoleLog(`${fn.name} failed: ${errMsg}`, 'error');
    }
    setConsoleBusy(false);
}

function showModal(html) {
    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    overlay.innerHTML = `<div class="modal">${html}</div>`;
    overlay.addEventListener('click', e => {
        if (e.target === overlay) overlay.remove();
    });
    document.body.appendChild(overlay);
    return overlay;
}

function field(id, label, inputHtml) {
    return `<div class="field"><label for="${id}">${label}</label>${inputHtml}</div>`;
}

function modalActions(primaryLabel, primaryFn) {
    return `<div class="modal-actions">
        <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
        <button class="btn btn-primary" onclick="${primaryFn}">${primaryLabel}</button>
    </div>`;
}

function optSelect(id, opts) {
    return `<select id="${id}">${opts.map(o => `<option value="${o}">${o}</option>`).join('')}</select>`;
}

// ─── REPO PAGE ──────────────────────────────────

function renderRepoPage(area) {
    area.innerHTML = `
        <div class="page-area">
            <div class="page-title">Repository Setup</div>
            <div style="display:flex;gap:8px;margin-bottom:20px">
                <button class="btn btn-primary" onclick="openRepo()">${icon('select', 14)} Select Repository</button>
                <button class="btn btn-secondary" onclick="initRepo()">${icon('init', 14)} Initialize</button>
            </div>
            <div class="repo-grid">
                <div class="repo-card">
                    <div class="repo-card-header">
                        <div class="repo-card-icon">${icon('stage', 16)}</div>
                        <div>
                            <h3>Create New Repository</h3>
                            <p>Paste GitHub commands to create and push</p>
                        </div>
                    </div>
                    <textarea id="newRepoCmds" placeholder="git init\ngit add .\ngit commit -m &quot;initial commit&quot;\ngh repo create ..."></textarea>
                    <div class="repo-actions"><button class="btn btn-primary btn-sm" onclick="runNewRepo()">${icon('play', 12)} Run</button></div>
                </div>
                <div class="repo-card">
                    <div class="repo-card-header">
                        <div class="repo-card-icon">${icon('remote', 16)}</div>
                        <div>
                            <h3>Push Existing Repository</h3>
                            <p>Connect and push an existing local repo</p>
                        </div>
                    </div>
                    <textarea id="existRepoCmds" placeholder="git remote add origin https://github.com/user/repo.git\ngit branch -M main\ngit push -u origin main"></textarea>
                    <div class="repo-actions"><button class="btn btn-primary btn-sm" onclick="runExistRepo()">${icon('play', 12)} Run</button></div>
                </div>
            </div>
        </div>
    `;
    window.openRepo = async () => {
        try {
            const path = await SelectRepo();
            if (path) { consoleLog(`Repository: ${path}`, 'success'); updateRepoInfo(); }
        } catch (e) { consoleLog(`Error: ${e}`, 'error'); }
    };
    window.initRepo = () => runGitCmd(Init);
    window.runNewRepo = async () => {
        if (!await checkRepo()) return;
        consoleLog('Running setup commands...');
        try { const r = await RunCommands(document.getElementById('newRepoCmds').value); consoleLog(r || 'Done.', 'success'); }
        catch (e) { consoleLog(`Error: ${e}`, 'error'); }
    };
    window.runExistRepo = async () => {
        if (!await checkRepo()) return;
        consoleLog('Running push commands...');
        try { const r = await RunCommands(document.getElementById('existRepoCmds').value); consoleLog(r || 'Done.', 'success'); }
        catch (e) { consoleLog(`Error: ${e}`, 'error'); }
    };
}

// ─── DASHBOARD PAGE ─────────────────────────────

function renderDashboardPage(area) {
    const cats = Object.keys(CATEGORIES);
    area.innerHTML = `
        <div class="dashboard-layout">
            <div class="dash-sidebar">
                <div class="dash-sidebar-label">Commands</div>
                ${cats.map(c => `
                    <button class="cat-btn${c === currentCategory ? ' active' : ''}"
                            onclick="switchCat('${c}')">${c}</button>
                `).join('')}
            </div>
            <div class="dash-content" id="dashContent"></div>
        </div>
    `;
    window.switchCat = (cat) => {
        currentCategory = cat;
        document.querySelectorAll('.cat-btn').forEach(b =>
            b.classList.toggle('active', b.textContent === cat));
        renderCatBtns(document.getElementById('dashContent'), cat);
    };
    renderCatBtns(document.getElementById('dashContent'), currentCategory);
}

function renderCatBtns(container, category) {
    const btns = CATEGORIES[category] || [];
    container.innerHTML = `
        <div class="section-header">
            <h3>${category}</h3>
            <span class="count">${btns.length} commands</span>
        </div>
        <div class="cmd-grid">
            ${btns.map((b, i) => `
                <button class="cmd-card" onclick="runCmd_${category}_${i}()">
                    ${icon(b.icon, 18)}
                    ${b.label}
                </button>
            `).join('')}
        </div>
    `;
    btns.forEach((b, i) => {
        window[`runCmd_${category}_${i}`] = b.action;
    });
}

// ─── DIALOG HELPERS ─────────────────────────────

function optSelectWithDefault(id, opts, def) {
    return `<select id="${id}">${opts.map(o =>
        `<option value="${o}"${o === def ? ' selected' : ''}>${o}</option>`
    ).join('')}</select>`;
}

function textInput(id, placeholder, value) {
    return `<input type="text" id="${id}" placeholder="${placeholder}" value="${value || ''}" />`;
}

async function loadBranches() {
    try { return await GetBranches(); } catch (_) { return []; }
}

// ─── DIALOGS ────────────────────────────────────

function runStatus() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('status', 14)}</div>
            <h3>Git Status</h3>
        </div>
        ${field('statusMode', 'Format', optSelect('statusMode', ['Standard', 'Short (-s)', 'Branch (-b)']))}
        ${modalActions('Run', `window._statusRun()`)}
    `);
    window._statusRun = async () => {
        const opt = document.getElementById('statusMode').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Status, opt);
    };
}

function runLog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('log', 14)}</div>
            <h3>Git Log</h3>
        </div>
        ${field('logMode', 'Format', optSelect('logMode', ['Oneline', 'Graph', 'Pretty']))}
        ${modalActions('Run', `window._logRun()`)}
    `);
    window._logRun = async () => {
        const opt = document.getElementById('logMode').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Log, opt);
    };
}

function showCommitDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('commit', 14)}</div>
            <h3>Git Commit</h3>
        </div>
        ${field('commitMsg', 'Message', textInput('commitMsg', 'Commit message'))}
        ${field('commitMode', 'Options', optSelect('commitMode', ['Standard (-m)', 'Stage All (-a)', 'Amend (--amend)']))}
        ${modalActions('Commit', `window._commitRun()`)}
    `);
    window._commitRun = async () => {
        const msg = document.getElementById('commitMsg').value.trim();
        const opt = document.getElementById('commitMode').value;
        if (!msg) { consoleLog('Commit message cannot be empty', 'error'); return; }
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Commit, msg, opt);
    };
}

async function showPushDialog() {
    const branches = await loadBranches();
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('push', 14)}</div>
            <h3>Git Push</h3>
        </div>
        ${field('pushBranch', 'Branch', optSelect('pushBranch', branches))}
        ${modalActions('Push', `window._pushRun()`)}
    `);
    window._pushRun = async () => {
        const branch = document.getElementById('pushBranch').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Push, branch);
    };
}

async function showPullDialog() {
    const branches = await loadBranches();
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('pull', 14)}</div>
            <h3>Git Pull</h3>
        </div>
        ${field('pullBranch', 'Branch', optSelect('pullBranch', branches))}
        ${modalActions('Pull', `window._pullRun()`)}
    `);
    window._pullRun = async () => {
        const branch = document.getElementById('pullBranch').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Pull, branch);
    };
}

function showBranchDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('branch', 14)}</div>
            <h3>Branch</h3>
        </div>
        ${field('branchName', 'Branch name', textInput('branchName', 'feature/my-branch'))}
        <div class="modal-actions">
            <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
            <button class="btn btn-danger btn-sm" onclick="window._branchDel()">Delete</button>
            <button class="btn btn-primary" onclick="window._branchCreate()">Create</button>
        </div>
    `);
    window._branchCreate = async () => {
        const n = document.getElementById('branchName').value.trim();
        if (!n) return;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(CreateBranch, n);
    };
    window._branchDel = async () => {
        const n = document.getElementById('branchName').value.trim();
        if (!n) return;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(DeleteBranch, n);
    };
}

async function showSwitchDialog() {
    const branches = await loadBranches();
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('branch', 14)}</div>
            <h3>Switch Branch</h3>
        </div>
        ${field('switchBranch', 'Branch', optSelect('switchBranch', branches))}
        ${modalActions('Switch', `window._switchRun()`)}
    `);
    window._switchRun = async () => {
        const branch = document.getElementById('switchBranch').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(SwitchBranch, branch);
        updateRepoInfo();
    };
}

async function showMergeDialog() {
    const branches = await loadBranches();
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('merge', 14)}</div>
            <h3>Merge Branch</h3>
        </div>
        ${field('mergeBranch', 'Branch to merge', optSelect('mergeBranch', branches))}
        ${modalActions('Merge', `window._mergeRun()`)}
    `);
    window._mergeRun = async () => {
        const branch = document.getElementById('mergeBranch').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Merge, branch);
    };
}

function showRenameDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('branch', 14)}</div>
            <h3>Rename Branch</h3>
        </div>
        ${field('oldName', 'Current name', textInput('oldName', 'main'))}
        ${field('newName', 'New name', textInput('newName', 'develop'))}
        ${modalActions('Rename', `window._renameRun()`)}
    `);
    window._renameRun = async () => {
        const o = document.getElementById('oldName').value.trim();
        const n = document.getElementById('newName').value.trim();
        if (!o || !n) return;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(BranchRename, o, n);
    };
}

function showTagDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('tag', 14)}</div>
            <h3>Tag</h3>
        </div>
        ${field('tagAction', 'Action', optSelect('tagAction', ['list', 'create', 'delete', 'push']))}
        ${field('tagName', 'Tag name', textInput('tagName', 'v1.0.0'))}
        ${modalActions('Run', `window._tagRun()`)}
    `);
    window._tagRun = async () => {
        const action = document.getElementById('tagAction').value;
        const name = document.getElementById('tagName').value.trim();
        document.querySelector('.modal-overlay').remove();
        if (action === 'list') { await runGitCmd(Tag, action, ''); }
        else if (!name) { consoleLog('Tag name required', 'error'); }
        else { await runGitCmd(Tag, action, name); }
    };
}

function showRemoteDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('remote', 14)}</div>
            <h3>Remote</h3>
        </div>
        ${field('remoteAction', 'Action', optSelect('remoteAction', ['list', 'add', 'remove']))}
        ${field('remoteVal', 'URL / Name', textInput('remoteVal', 'https://github.com/...'))}
        ${modalActions('Run', `window._remoteRun()`)}
    `);
    window._remoteRun = async () => {
        const action = document.getElementById('remoteAction').value;
        const val = document.getElementById('remoteVal').value.trim();
        document.querySelector('.modal-overlay').remove();
        if (action === 'list') { await runGitCmd(Remote, 'list', ''); }
        else if (action === 'remove' && !val) { consoleLog('Remote name required', 'error'); }
        else { await runGitCmd(Remote, action, val); }
    };
}

function showRevertDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('revert', 14)}</div>
            <h3>Revert</h3>
        </div>
        ${field('revertHash', 'Commit hash', textInput('revertHash', 'a1b2c3d'))}
        ${modalActions('Revert', `window._revertRun()`)}
    `);
    window._revertRun = async () => {
        const h = document.getElementById('revertHash').value.trim();
        if (!h) { consoleLog('Commit hash required', 'error'); return; }
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Revert, h);
    };
}

function showCloneDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('clone', 14)}</div>
            <h3>Clone Repository</h3>
        </div>
        ${field('cloneUrl', 'Repository URL', textInput('cloneUrl', 'https://github.com/user/repo.git'))}
        ${modalActions('Clone', `window._cloneRun()`)}
    `);
    window._cloneRun = async () => {
        const url = document.getElementById('cloneUrl').value.trim();
        if (!url) { consoleLog('URL required', 'error'); return; }
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Clone, url);
    };
}

function showDiffDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('diff', 14)}</div>
            <h3>Git Diff</h3>
        </div>
        ${field('diffMode', 'Mode', optSelect('diffMode', ['Unstaged', 'Staged (--cached)', 'Names (--name-only)', 'Summary (--stat)']))}
        ${modalActions('Run', `window._diffRun()`)}
    `);
    window._diffRun = async () => {
        const opt = document.getElementById('diffMode').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Diff, opt);
    };
}

function showResetDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('reset', 14)}</div>
            <h3>Git Reset</h3>
        </div>
        ${field('resetMode', 'Mode', optSelect('resetMode', ['--mixed', '--soft', '--hard']))}
        ${field('resetTarget', 'Target', textInput('resetTarget', 'HEAD~1', 'HEAD~1'))}
        ${modalActions('Reset', `window._resetRun()`)}
    `);
    window._resetRun = async () => {
        const mode = document.getElementById('resetMode').value;
        const target = document.getElementById('resetTarget').value.trim();
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Reset, mode, target);
    };
}

function showStashDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('stash', 14)}</div>
            <h3>Git Stash</h3>
        </div>
        ${field('stashAction', 'Action', optSelect('stashAction', ['Save', 'Pop', 'List', 'Drop', 'Apply']))}
        ${modalActions('Run', `window._stashRun()`)}
    `);
    window._stashRun = async () => {
        const action = document.getElementById('stashAction').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Stash, action);
    };
}

function showCleanDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('clean', 14)}</div>
            <h3>Git Clean</h3>
        </div>
        ${field('cleanMode', 'Mode', optSelect('cleanMode', ['Preview (-n)', 'Remove Dir (-d)', 'Force (-f)', 'Full (-fdx)']))}
        ${modalActions('Run', `window._cleanRun()`)}
    `);
    window._cleanRun = async () => {
        const opt = document.getElementById('cleanMode').value;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Clean, opt);
    };
}

function showRebaseDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('rebase', 14)}</div>
            <h3>Git Rebase</h3>
        </div>
        ${field('rebaseAction', 'Action', optSelect('rebaseAction', ['Interactive (-i)', 'Onto', 'Continue', 'Abort', 'Skip']))}
        ${field('rebaseTarget', 'Target', textInput('rebaseTarget', 'main'))}
        ${modalActions('Run', `window._rebaseRun()`)}
    `);
    window._rebaseRun = async () => {
        const action = document.getElementById('rebaseAction').value;
        const target = document.getElementById('rebaseTarget').value.trim();
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Rebase, action, target);
    };
}

function showCherryPickDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('cherry', 14)}</div>
            <h3>Cherry-pick</h3>
        </div>
        ${field('cherryHash', 'Commit hash', textInput('cherryHash', 'a1b2c3d'))}
        ${modalActions('Apply', `window._cherryRun()`)}
    `);
    window._cherryRun = async () => {
        const h = document.getElementById('cherryHash').value.trim();
        if (!h) return;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(CherryPick, h);
    };
}

function showShowDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('show', 14)}</div>
            <h3>Git Show</h3>
        </div>
        ${field('showMode', 'Mode', optSelect('showMode', ['Head', 'Last 5', 'Specific']))}
        ${field('showHash', 'Commit hash (for Specific)', textInput('showHash', ''))}
        ${modalActions('Show', `window._showRun()`)}
    `);
    window._showRun = async () => {
        const mode = document.getElementById('showMode').value;
        const hash = document.getElementById('showHash').value.trim();
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Show, mode, mode === 'Specific' ? hash : '');
    };
}

function showBlameDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('blame', 14)}</div>
            <h3>Git Blame</h3>
        </div>
        ${field('blameFile', 'File path', textInput('blameFile', 'src/main.go'))}
        ${modalActions('Blame', `window._blameRun()`)}
    `);
    window._blameRun = async () => {
        const f = document.getElementById('blameFile').value.trim();
        if (!f) return;
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Blame, f);
    };
}

function showWorktreeDialog() {
    showModal(`
        <div class="modal-header">
            <div class="modal-header-icon">${icon('worktree', 14)}</div>
            <h3>Worktree</h3>
        </div>
        ${field('wtAction', 'Action', optSelect('wtAction', ['List', 'Add', 'Remove', 'Prune']))}
        ${field('wtArgs', 'Args (for Add/Remove)', textInput('wtArgs', '/path branch-name'))}
        ${modalActions('Run', `window._wtRun()`)}
    `);
    window._wtRun = async () => {
        const action = document.getElementById('wtAction').value;
        const args = document.getElementById('wtArgs').value.trim();
        document.querySelector('.modal-overlay').remove();
        await runGitCmd(Worktree, action, args);
    };
}

async function showConflictsDialog() {
    try {
        const conflicts = await GetConflicts();
        if (!conflicts || conflicts.length === 0) {
            consoleLog('No merge conflicts detected.', 'success');
            return;
        }
        const overlay = document.createElement('div');
        overlay.className = 'modal-overlay';
        overlay.innerHTML = `<div class="modal">
            <div class="modal-header">
                <div class="modal-header-icon">${icon('conflict', 14)}</div>
                <h3>Resolve Conflicts</h3>
            </div>
            ${conflicts.map(f => `
                <div class="conflict-row">
                    <span>${f}</span>
                    <button class="btn btn-primary btn-sm" onclick="window._conflictResolve('${f}','ours')">Keep Mine</button>
                    <button class="btn btn-danger btn-sm" onclick="window._conflictResolve('${f}','theirs')">Take Theirs</button>
                </div>
            `).join('')}
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Close</button>
            </div>
        </div>`;
        overlay.addEventListener('click', e => { if (e.target === overlay) overlay.remove(); });
        document.body.appendChild(overlay);
        window._conflictResolve = async (file, strategy) => {
            overlay.remove();
            await runGitCmd(ResolveConflict, file, strategy);
        };
    } catch (err) { consoleLog(`Error: ${err}`, 'error'); }
}

async function showGitIgnoreDialog() {
    try {
        const content = await ReadGitIgnore() || '';
        const overlay = document.createElement('div');
        overlay.className = 'modal-overlay';
        overlay.innerHTML = `<div class="modal" style="min-width:500px">
            <div class="modal-header">
                <div class="modal-header-icon">${icon('gitignore', 14)}</div>
                <h3>.gitignore</h3>
            </div>
            <textarea id="gitignoreContent" style="min-height:220px;font-family:var(--font-mono);font-size:12px;line-height:1.7">${content}</textarea>
            <div class="modal-actions">
                <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">Cancel</button>
                <button class="btn btn-primary" onclick="window._saveGitignore()">Save</button>
            </div>
        </div>`;
        overlay.addEventListener('click', e => { if (e.target === overlay) overlay.remove(); });
        document.body.appendChild(overlay);
        window._saveGitignore = async () => {
            const text = document.getElementById('gitignoreContent').value;
            try {
                await WriteGitIgnore(text);
                consoleLog('.gitignore saved.', 'success');
                overlay.remove();
            } catch (err) { consoleLog(`Error: ${err}`, 'error'); }
        };
    } catch (err) { consoleLog(`Error: ${err}`, 'error'); }
}

// ─── ABOUT PAGE ─────────────────────────────────

function renderAboutPage(area) {
    area.innerHTML = `
        <div class="page-area about-page">
            <div class="about-logo">${icon('git', 32)}</div>
            <h2>GitScope</h2>
            <p>A modern, lightweight Git client built with Go and Wails. Simplifies version control operations for both beginners and experienced developers.</p>
            <p style="font-size:11px;color:var(--text-muted)">Version 1.0.0</p>
            <a class="about-link" href="https://github.com/Aswanidev-vs/GitScope" target="_blank">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                View on GitHub
            </a>
        </div>
    `;
}

// ─── DOCS PAGE ──────────────────────────────────

const DOC_ITEMS = [
    "Init","Stage","Status","Commit","Push","Log","Revert","Clone","Branch","Pull",
    "Reflog","GitIgnore","Remote","Diff","Reset","Fetch","Stash","Merge","Tag",
    "Cherry-pick","Rebase","Clean","Show","Ls-files","Worktree","Shortlog","Blame"
];

const DOC_TEXT = {
    Init: `Git Init\n${'='.repeat(40)}\n\ngit init creates a new empty Git repository in your current folder.\n\nCommand:\n  git init\n\nAfter running this, Git creates a hidden .git folder\nthat stores all version history and settings.\n\nExample:\n  mkdir myproject && cd myproject && git init`,
    Stage: `Git Stage\n${'='.repeat(40)}\n\nThe staging area (index) is where Git stores changes\nyou want to include in your next commit.\n\nCommands:\n  git add file.txt    (stage a single file)\n  git add .           (stage all files)\n  git add -u          (stage modified/deleted)`,
    Status: `Git Status\n${'='.repeat(40)}\n\ngit status shows the current state of your working\ndirectory and staging area.\n\nCommand:\n  git status\n  git status -s       (short format)\n  git status -b       (show branch info)`,
    Commit: `Git Commit\n${'='.repeat(40)}\n\nA commit is a snapshot of your project at a point in time.\n\nCommands:\n  git commit -m "message"\n  git commit -a -m "message"   (stage + commit)\n  git commit --amend            (amend last commit)`,
    Push: `Git Push\n${'='.repeat(40)}\n\nPush sends your local commits to a remote repository.\n\nCommands:\n  git push\n  git push origin branchname\n  git push -u origin branchname`,
    Log: `Git Log\n${'='.repeat(40)}\n\nLog shows the full history of commits.\n\nCommands:\n  git log\n  git log --oneline\n  git log --graph --oneline --decorate --all\n  git log -p                  (show diffs)`,
    Revert: `Git Revert\n${'='.repeat(40)}\n\nRevert undoes a specific commit by creating a new one.\n\nCommands:\n  git revert <commit-id>\n  git revert --no-commit <old>..<new>`,
    Clone: `Git Clone\n${'='.repeat(40)}\n\nClone creates a local copy of a remote repository.\n\nCommands:\n  git clone <url>\n  git clone <url> myproject\n  git clone --depth 1 <url>`,
    Branch: `Git Branch\n${'='.repeat(40)}\n\nA branch is a separate line of development.\n\nCommands:\n  git branch                      (list)\n  git branch feature-login        (create)\n  git switch feature-login        (switch)\n  git switch -c feature-login     (create + switch)\n  git branch -d feature-login     (delete)\n  git branch -m new-name          (rename)`,
    Pull: `Git Pull\n${'='.repeat(40)}\n\ngit pull brings remote changes into your current branch.\n\nCommands:\n  git pull\n  git pull origin main\n  git pull --rebase`,
    Reflog: `Git Reflog\n${'='.repeat(40)}\n\nReflog shows the history of where HEAD has been.\n\nCommands:\n  git reflog\n  git reflog show HEAD`,
    GitIgnore: `.gitignore\n${'='.repeat(40)}\n\nTells Git which files/folders to NOT track.\n\nCommon entries:\n  *.DS_Store\n  build/\n  dist/\n  *.log\n  .env\n  .vscode/`,
    Remote: `Remote\n${'='.repeat(40)}\n\ngit remote manages connections to other repositories.\n\nCommands:\n  git remote -v              (list)\n  git remote add name url   (add)\n  git remote remove name    (remove)`,
    Diff: `Git Diff\n${'='.repeat(40)}\n\nShows differences between file versions.\n\nCommands:\n  git diff                    (unstaged)\n  git diff --staged          (staged)\n  git diff HEAD              (all changes)\n  git diff --stat            (summary)`,
    Reset: `Git Reset\n${'='.repeat(40)}\n\nMoves HEAD and controls commit history.\n\nModes:\n  --soft   undo commit, keep staged\n  --mixed  undo commit, unstage (default)\n  --hard   delete everything (destructive)\n\nCommands:\n  git reset HEAD^\n  git reset --soft HEAD^\n  git reset --hard HEAD^`,
    Fetch: `Git Fetch\n${'='.repeat(40)}\n\nDownloads changes from remote without merging.\n\nCommands:\n  git fetch origin\n  git fetch --all`,
    Stash: `Git Stash\n${'='.repeat(40)}\n\nTemporarily shelves changes for a clean directory.\n\nCommands:\n  git stash\n  git stash list\n  git stash pop\n  git stash apply`,
    Merge: `Git Merge\n${'='.repeat(40)}\n\nJoins two development histories together.\n\nCommands:\n  git merge feature-x\n  git merge --no-ff feature-x`,
    Tag: `Git Tag\n${'='.repeat(40)}\n\nMarks specific points as important (releases).\n\nCommands:\n  git tag v1.0\n  git push origin v1.0`,
    "Cherry-pick": `Cherry-pick\n${'='.repeat(40)}\n\nApplies changes from existing commits to current branch.\n\nCommands:\n  git cherry-pick <commit-hash>`,
    Rebase: `Git Rebase\n${'='.repeat(40)}\n\nReapplies commits on top of another base tip.\n\nCommands:\n  git rebase main\n  git rebase -i HEAD~3\n  git rebase --continue\n  git rebase --abort`,
    Clean: `Git Clean\n${'='.repeat(40)}\n\nRemoves untracked files from working tree.\n\nCommands:\n  git clean -n    (preview)\n  git clean -f    (remove files)\n  git clean -fd   (remove files + dirs)`,
    Show: `Git Show\n${'='.repeat(40)}\n\nShows details about a Git object.\n\nCommands:\n  git show HEAD\n  git show <hash>\n  git show --stat`,
    "Ls-files": `Git Ls-files\n${'='.repeat(40)}\n\nShows files in the index and working tree.\n\nCommands:\n  git ls-files\n  git ls-files --cached\n  git ls-files --others`,
    Worktree: `Git Worktree\n${'='.repeat(40)}\n\nManage multiple working trees.\n\nCommands:\n  git worktree list\n  git worktree add <path> <branch>\n  git worktree remove <name>`,
    Shortlog: `Git Shortlog\n${'='.repeat(40)}\n\nSummarizes git log grouped by author.\n\nCommands:\n  git shortlog\n  git shortlog -s\n  git shortlog -n`,
    Blame: `Git Blame\n${'='.repeat(40)}\n\nShows what revision/author last modified each line.\n\nCommands:\n  git blame <file>\n  git blame -L 10,20 <file>`,
};

function renderDocsPage(area) {
    area.innerHTML = `
        <div class="page-area">
            <div class="page-title">Documentation</div>
            <div class="doc-grid" id="docGrid"></div>
            <div class="doc-body" id="docBody">Select a command from above to view its documentation.</div>
        </div>
    `;
    const grid = document.getElementById('docGrid');
    DOC_ITEMS.forEach(item => {
        const btn = document.createElement('button');
        btn.className = 'doc-pill';
        btn.textContent = item;
        btn.onclick = () => {
            document.querySelectorAll('.doc-pill').forEach(b => b.classList.remove('active'));
            btn.classList.add('active');
            document.getElementById('docBody').textContent = DOC_TEXT[item] || 'No documentation available.';
        };
        grid.appendChild(btn);
    });
}
