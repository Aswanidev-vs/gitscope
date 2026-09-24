/* ════════════════════════════════════════════════════════════════
   GitScope — frontend controller
   Every value in this file comes from the Go backend through the
   Wails bindings below. There is no sample/mock data: with no
   repository selected the UI shows empty states, and all git output
   is whatever `git` actually printed.

   Pages: Repository · Dashboard · History · Docs · About  (keys 1-5)
   Shortcuts: Ctrl/⌘ K palette · 1-5 pages · / filter · Esc close
   ════════════════════════════════════════════════════════════════ */

import './style.css';
import { icon } from './icons.js';

import {
    Blame, BranchRename, CherryPick, Clean, Clone, Commit, CreateBranch, DeleteBranch,
    Diff, Fetch, GetBranches, GetConflicts, GetCurrentBranch, GetRepoPath, Init,
    IsGitAvailable, IsRepoInitialized, Log, LsFiles, MagicSync, Merge, Pull, Push,
    ReadGitIgnore, Rebase, Reflog, Remote, Reset, ResolveConflict, Revert, RunCommands,
    SelectRepo, Shortlog, Show, Stage, Stash, Status, SwitchBranch, Tag, UndoLastCommit,
    Worktree, WriteGitIgnore,
} from '../wailsjs/go/main/App';

/* ─── BACKEND STATE (mirrors what Go reports) ─────── */

const state = {
    repoPath: '',
    branch: '',
    upstream: '',
    ahead: null,          // null = no upstream / unknown
    behind: null,
    gitAvailable: true,
    initialized: false,
    branches: [],
    conflicts: [],
    staged: [],           // parsed from `git status -s`
    changed: [],
    busy: false,
    history: { format: 'Oneline', raw: '', rows: [], loaded: false },
};

let currentPage = 'repo';
let currentCategory = 'Common';
let lastFocused = null;
let pollTimer = null;

/* ─── HELPERS ─────────────────────────────────────── */

const $ = (sel, root = document) => root.querySelector(sel);
const $$ = (sel, root = document) => [...root.querySelectorAll(sel)];
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) =>
    ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

function errMsg(err) {
    if (!err) return 'Unknown error';
    if (typeof err === 'string') return err;
    if (err.message) return err.message;
    try { return JSON.stringify(err); } catch (_) { return String(err); }
}

/* run a backend call, falling back instead of throwing */
async function safe(fn, fallback) {
    try {
        const out = await fn();
        return out === undefined || out === null ? fallback : out;
    } catch (_) {
        return fallback;
    }
}

function hydrateIcons(root) {
    $$('[data-icon]', root || document).forEach((el) => {
        el.innerHTML = icon(el.dataset.icon, Number(el.dataset.size) || 18);
        el.style.display = 'inline-flex';
    });
}

/* ─── CONSOLE ─────────────────────────────────────── */

function log(msg, type) {
    const out = $('#consoleOut');
    if (!out) return;
    const cls = type || 'default';
    const prefix = type === 'error' ? '[err] ' : type === 'warning' ? '[warn] ' : '';
    const span = document.createElement('span');
    span.className = 'cl-' + cls;
    span.textContent = prefix + String(msg ?? '') + '\n';
    out.appendChild(span);
    out.scrollTop = out.scrollHeight;
    if (type === 'error' || type === 'warning') toast(msg, type === 'error' ? 'error' : 'warning');
}

function clearConsole() {
    const out = $('#consoleOut');
    if (out) out.innerHTML = '';
    log('Console cleared.', 'dim');
}

function setBusy(on) {
    state.busy = on;
    $('#busyLed')?.classList.toggle('is-busy', on);
    const t = $('#busyText');
    if (t) t.textContent = on ? 'Running…' : 'Idle';
}

/* Run a backend call and echo the real output to the console. */
async function run(label, fn, opts = {}) {
    if (opts.needsRepo !== false && !state.repoPath) {
        log('Select a repository first.', 'warning');
        return null;
    }
    setBusy(true);
    log('$ ' + label, 'cmd');
    try {
        const result = await fn();
        const text = typeof result === 'string' ? result.trim() : '';
        if (text) log(text, 'default');
        else log('Done.', 'dim');
        if (opts.refresh) await loadRepo();
        return result;
    } catch (err) {
        log(errMsg(err), 'error');
        return null;
    } finally {
        setBusy(false);
    }
}

/* ─── TOASTS ──────────────────────────────────────── */

function toast(msg, type) {
    const root = $('#toastRoot');
    if (!root) return;
    const kind = type || 'info';
    const ic = { success: 'check', error: 'xCircle', warning: 'warn', info: 'info' }[kind] || 'info';
    const el = document.createElement('div');
    el.className = 'toast t-' + kind;
    el.setAttribute('role', kind === 'error' ? 'alert' : 'status');
    el.innerHTML = icon(ic, 15);
    const span = document.createElement('span');
    span.textContent = String(msg ?? '');
    el.appendChild(span);
    root.appendChild(el);
    setTimeout(() => {
        el.classList.add('is-out');
        setTimeout(() => el.remove(), 220);
    }, 3600);
}

/* ─── CHROME (top bar + console status) ───────────── */

function syncChrome() {
    const pathEl = $('#repoPath');
    if (pathEl) {
        pathEl.textContent = state.repoPath || 'No repository selected';
        pathEl.title = state.repoPath || '';
        pathEl.classList.toggle('is-empty', !state.repoPath);
    }

    const branchChip = $('#branchChip');
    if (branchChip) {
        branchChip.hidden = !state.repoPath;
        const nameEl = $('#branchName');
        if (nameEl) nameEl.textContent = state.branch || 'detached';
    }

    const ahead = $('#aheadChip');
    const behind = $('#behindChip');
    const known = state.ahead !== null || state.behind !== null;
    if (ahead) {
        ahead.hidden = !state.repoPath || !known;
        ahead.textContent = '↑ ' + (state.ahead ?? 0);
        ahead.classList.toggle('is-zero', !state.ahead);
        ahead.title = state.upstream ? 'Ahead of ' + state.upstream : 'Ahead of upstream';
    }
    if (behind) {
        behind.hidden = !state.repoPath || !known;
        behind.textContent = '↓ ' + (state.behind ?? 0);
        behind.classList.toggle('is-zero', !state.behind);
        behind.title = state.upstream ? 'Behind ' + state.upstream : 'Behind upstream';
    }

    const magic = $('#magicSyncBtn');
    if (magic) magic.disabled = !state.repoPath;

    const copyBtn = $('#copyPathBtn');
    if (copyBtn) copyBtn.hidden = !state.repoPath;

    const led = $('#busyLed');
    if (led && !state.busy) led.classList.toggle('is-bad', !state.gitAvailable);

    const meta = $('#statusMeta');
    if (meta) {
        if (!state.gitAvailable) meta.textContent = 'git not found on PATH';
        else if (!state.repoPath) meta.textContent = 'no repository selected';
        else {
            meta.textContent = [
                state.branch || 'detached',
                state.staged.length + ' staged',
                state.changed.length + ' changed',
                state.conflicts.length ? state.conflicts.length + ' conflicts' : null,
            ].filter(Boolean).join(' · ');
        }
    }
}

/* ─── DIALOG ENGINE ─────────────────────────────────
   spec = { title, icon, desc, wide, fields:[{id,label,type,options,
   value,placeholder,required,mono,errorText}],
   danger:{when(values),text}, actions:[{id,label,kind,onRun(values)}] }
   role=dialog + focus trap + Esc + focus restore + inline validation +
   two-step confirmation for destructive actions. */

function collectValues(dlg) {
    const v = {};
    (dlg.dataset.fields || '').split(',').filter(Boolean).forEach((id) => {
        const el = dlg.querySelector('#f_' + id);
        if (el) v[id] = el.value;
    });
    return v;
}

function openDialog(spec) {
    closeDialog();
    lastFocused = document.activeElement;

    const overlay = document.createElement('div');
    overlay.className = 'modal-overlay';
    const dlg = document.createElement('div');
    dlg.className = 'dialog' + (spec.wide ? ' dialog-wide' : '');
    dlg.setAttribute('role', 'dialog');
    dlg.setAttribute('aria-modal', 'true');
    dlg.setAttribute('aria-labelledby', 'dlgTitle');
    dlg.dataset.fields = (spec.fields || []).map((f) => f.id).join(',');

    const fieldsHTML = (spec.fields || []).map((f) => {
        const id = 'f_' + f.id;
        let control;
        if (f.type === 'select') {
            const opts = (typeof f.options === 'function' ? f.options() : f.options) || [];
            control = `<select id="${id}"${f.mono ? ' class="mono"' : ''}>` +
                opts.map((o) => `<option${o === f.value ? ' selected' : ''}>${esc(o)}</option>`).join('') +
                '</select>';
        } else if (f.type === 'textarea') {
            control = `<textarea id="${id}"${f.mono ? ' class="mono"' : ''} placeholder="${esc(f.placeholder || '')}">${esc(f.value || '')}</textarea>`;
        } else {
            control = `<input type="text" id="${id}"${f.mono ? ' class="mono"' : ''}` +
                ` placeholder="${esc(f.placeholder || '')}" value="${esc(f.value || '')}"` +
                (f.required ? ' required aria-required="true"' : '') + ' autocomplete="off" />';
        }
        return `<div class="field" data-field="${f.id}">` +
            `<label class="field-label" for="${id}">${esc(f.label)}${f.required ? ' *' : ''}</label>` +
            control +
            `<p class="field-error" id="err_${f.id}" role="alert"></p></div>`;
    }).join('');

    const actions = spec.actions || [{ id: 'run', label: 'Run', kind: 'primary' }];
    const footHTML = actions.map((a) => {
        const cls = a.kind === 'danger' ? 'btn-danger' : a.kind === 'secondary' ? 'btn-secondary' : 'btn-primary';
        return `<button type="button" class="btn ${cls} btn-sm" data-action="${a.id}">${esc(a.label)}</button>`;
    }).join('');

    dlg.innerHTML =
        `<div class="dialog-head">` +
        `<span class="dialog-icon">${icon(spec.icon || 'git', 16)}</span>` +
        `<h2 id="dlgTitle">${esc(spec.title)}</h2>` +
        `<button type="button" class="dialog-close" data-action="__cancel" aria-label="Close dialog">${icon('clear', 15)}</button>` +
        `</div>` +
        `<div class="dialog-body">` +
        (spec.desc ? `<p class="dialog-desc">${spec.desc}</p>` : '') +
        fieldsHTML +
        (spec.bodyHTML || '') +
        `<div class="danger-strip" id="dangerStrip" hidden>${icon('warn', 15)}<span id="dangerText"></span></div>` +
        `<label class="check" id="confirmWrap" hidden><input type="checkbox" id="confirmChk" /> <span>I understand the consequences</span></label>` +
        `<p class="field-error" id="confirmErr" role="alert" style="margin-top:8px"></p>` +
        `</div>` +
        `<div class="dialog-foot">` +
        `<button type="button" class="btn btn-secondary btn-sm" data-action="__cancel">Cancel</button>` +
        `<span class="spacer"></span>${footHTML}</div>`;

    overlay.appendChild(dlg);
    $('#modalRoot').appendChild(overlay);
    hydrateIcons(dlg);

    const strip = dlg.querySelector('#dangerStrip');
    const stripText = dlg.querySelector('#dangerText');
    const confirmWrap = dlg.querySelector('#confirmWrap');
    const confirmChk = dlg.querySelector('#confirmChk');
    const confirmErr = dlg.querySelector('#confirmErr');
    let forcedDanger = false;

    function evalDanger() {
        const active = forcedDanger ||
            !!(spec.danger && spec.danger.when && spec.danger.when(collectValues(dlg)));
        strip.hidden = !active;
        confirmWrap.hidden = !active;
        if (active) stripText.innerHTML = (spec.danger && spec.danger.text) || 'This operation cannot be undone.';
        if (!active) { confirmChk.checked = false; confirmErr.style.display = 'none'; }
        return active;
    }
    dlg.addEventListener('change', evalDanger);
    dlg.addEventListener('input', evalDanger);
    evalDanger();

    let focusFirstInvalid = false;
    function validate() {
        let ok = true;
        focusFirstInvalid = false;
        (spec.fields || []).forEach((f) => {
            if (!f.required) return;
            const el = dlg.querySelector('#f_' + f.id);
            const wrap = dlg.querySelector(`[data-field="${f.id}"]`);
            const err = dlg.querySelector('#err_' + f.id);
            if (!el.value.trim()) {
                ok = false;
                el.setAttribute('aria-invalid', 'true');
                el.setAttribute('aria-describedby', 'err_' + f.id);
                wrap.classList.add('has-error');
                err.textContent = f.errorText || f.label + ' is required.';
                if (!focusFirstInvalid) { el.focus(); focusFirstInvalid = true; }
            } else {
                el.removeAttribute('aria-invalid');
                wrap.classList.remove('has-error');
                err.textContent = '';
            }
        });
        return ok;
    }

    dlg.addEventListener('click', (e) => {
        const btn = e.target.closest('[data-action]');
        if (!btn) return;
        const actId = btn.dataset.action;
        if (actId === '__cancel') { closeDialog(); return; }
        const action = actions.find((a) => a.id === actId);
        if (!action) return;
        if (!validate()) return;

        if (action.kind === 'danger' && !forcedDanger) {
            forcedDanger = true;
            if (!spec.danger) spec.danger = {};
            if (!spec.danger.text) spec.danger.text = action.confirmText || 'This operation cannot be undone.';
            evalDanger();
            btn.textContent = 'Confirm — ' + action.label;
            confirmChk.focus();
            return;
        }
        if (evalDanger() && !confirmChk.checked) {
            confirmErr.textContent = 'Tick the confirmation box to proceed.';
            confirmErr.style.display = 'block';
            confirmChk.focus();
            return;
        }
        const values = collectValues(dlg);
        closeDialog();
        if (action.onRun) action.onRun(values);
    });

    dlg.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') { e.stopPropagation(); closeDialog(); return; }
        if (e.key !== 'Tab') return;
        const focusables = $$('button, input, select, textarea, a[href]', dlg)
            .filter((el) => !el.disabled && el.offsetParent !== null);
        if (!focusables.length) return;
        const first = focusables[0];
        const last = focusables[focusables.length - 1];
        if (e.shiftKey && document.activeElement === first) { e.preventDefault(); last.focus(); }
        else if (!e.shiftKey && document.activeElement === last) { e.preventDefault(); first.focus(); }
    });
    overlay.addEventListener('mousedown', (e) => { if (e.target === overlay) closeDialog(); });

    const firstInput = dlg.querySelector('.field input, .field select, .field textarea') || dlg.querySelector('.dialog-close');
    firstInput.focus();
    return dlg;
}

function closeDialog() {
    const root = $('#modalRoot');
    if (root && root.firstChild) {
        root.innerHTML = '';
        if (lastFocused && document.contains(lastFocused)) lastFocused.focus();
        lastFocused = null;
    }
}

function dialogOpen() { return !!( $('#modalRoot') && $('#modalRoot').firstChild); }

/* ─── REPO DATA (all values are parsed from real git output) ─── */

/* `git status -s` → { staged, changed }. Line = XY<space>path */
function parseShortStatus(text) {
    const staged = [];
    const changed = [];
    String(text || '').split('\n').forEach((line) => {
        if (!line || line.length < 3 || line.startsWith('##')) return;
        const x = line[0];
        const y = line[1];
        let file = line.slice(3).trim();
        if (file.includes(' -> ')) file = file.split(' -> ').pop();
        if (x !== ' ' && x !== '?' && x !== '!') staged.push({ code: x, file });
        if (y !== ' ' || x === '?') changed.push({ code: x === '?' ? 'U' : y, file });
    });
    return { staged, changed };
}

/* `git status -b` → { upstream, ahead, behind } (null when unknown) */
function parseBranchStatus(text) {
    const t = String(text || '');
    const out = { upstream: '', ahead: null, behind: null };
    const upToDate = /up to date with '([^']+)'/.exec(t);
    const aheadRe = /ahead of '([^']+)' by (\d+) commit/.exec(t);
    const behindRe = /behind '([^']+)' by (\d+) commit/.exec(t);
    const divergedRe = /and '([^']+)' have diverged,\s*and have (\d+) and (\d+) different commits? each/.exec(t);
    const basedOn = /Your branch is based on '([^']+)'/.exec(t);
    if (upToDate) { out.upstream = upToDate[1]; out.ahead = 0; out.behind = 0; }
    if (aheadRe) { out.upstream = aheadRe[1]; out.ahead = Number(aheadRe[2]); out.behind = 0; }
    if (behindRe) { out.upstream = behindRe[1]; out.behind = Number(behindRe[2]); out.ahead = 0; }
    if (divergedRe) { out.upstream = divergedRe[1]; out.ahead = Number(divergedRe[2]); out.behind = Number(divergedRe[3]); }
    if (!out.upstream && basedOn) out.upstream = basedOn[1];
    return out;
}

/* git can leak messages (e.g. autocrlf "LF will be replaced by CRLF …")
   into command output; keep only entries that are real unmerged paths. */
function sanitizeConflicts(list) {
    return (Array.isArray(list) ? list : [])
        .map((s) => String(s).trim())
        .filter((s) => s && !/^(warning|error|fatal|hint|note):/i.test(s));
}

/* Preserve in-progress text across page re-renders (the 5s status poll and
   action-triggered refreshes call renderPage() — without this, a rebuild
   recreates the textarea and wipes whatever the user was typing). */
let pageSig = '';
function captureDraft() {
    const d = {};
    document.querySelectorAll('#pageRoot textarea, #pageRoot input[type="text"]').forEach((el) => {
        if (el.id) d[el.id] = { v: el.value, s: el.selectionStart, e: el.selectionEnd, f: document.activeElement === el };
    });
    return d;
}
function restoreDraft(d) {
    Object.keys(d).forEach((id) => {
        const el = document.getElementById(id);
        if (!el) return;
        el.value = d[id].v;
        try { el.setSelectionRange(d[id].s, d[id].e); } catch (_) { /* inputs without selection support */ }
        if (d[id].f) el.focus();
    });
}

/* Pull everything the UI shows straight from the backend. */
async function loadRepo() {
    state.gitAvailable = await safe(IsGitAvailable, false);
    state.repoPath = await safe(GetRepoPath, '');

    if (!state.repoPath) {
        state.branch = '';
        state.upstream = '';
        state.ahead = null;
        state.behind = null;
        state.initialized = false;
        state.branches = [];
        state.conflicts = [];
        state.staged = [];
        state.changed = [];
        syncChrome();
        return;
    }

    state.initialized = await safe(IsRepoInitialized, false);

    if (!state.initialized) {
        state.branch = '';
        state.upstream = '';
        state.ahead = null;
        state.behind = null;
        state.branches = [];
        state.conflicts = [];
        state.staged = [];
        state.changed = [];
        syncChrome();
        return;
    }

    state.branch = await safe(GetCurrentBranch, '');
    state.branches = await safe(GetBranches, []);

    /* files come from `git status -s`; sync state from `git status -b` */
    const shortText = await safe(() => Status('Short (-s)'), '');
    const parsed = parseShortStatus(shortText);
    state.staged = parsed.staged;
    state.changed = parsed.changed;

    const branchText = await safe(() => Status('Branch (-b)'), '');
    Object.assign(state, parseBranchStatus(branchText));

    state.conflicts = sanitizeConflicts(await safe(GetConflicts, []));

    syncChrome();
    if (currentPage === 'repo') {
        /* Rebuild only when git state actually changed, and keep any text
           the user is composing while the page swaps out. */
        const sig = [state.gitAvailable, state.repoPath, state.initialized, state.branch,
            state.upstream, state.ahead, state.behind, state.branches.join('\n'),
            state.conflicts.join('\n'), state.staged.join('\n'), state.changed.join('\n')].join('\u0000');
        if (sig !== pageSig) {
            const draft = captureDraft();
            pageSig = sig;
            renderPage();
            restoreDraft(draft);
        }
    }
}

function startPolling() {
    if (pollTimer) clearInterval(pollTimer);
    pollTimer = setInterval(() => {
        if (state.busy || dialogOpen() || paletteOpen) return;
        loadRepo();
    }, 5000);
}

/* ─── NAVIGATION / SHELL ──────────────────────────── */

const PAGES = [
    { id: 'repo', label: 'Repository', icon: 'repo', num: '1' },
    { id: 'dashboard', label: 'Dashboard', icon: 'dashboard', num: '2' },
    { id: 'history', label: 'History', icon: 'history', num: '3' },
    { id: 'docs', label: 'Docs', icon: 'docs', num: '4' },
    { id: 'about', label: 'About', icon: 'about', num: '5' },
];

function renderShell() {
    $('#app').innerHTML = `
        <a class="skip-link" href="#pageRoot">Skip to content</a>
        <div class="grain" aria-hidden="true"></div>

        <nav class="rail" aria-label="Primary">
            <div class="rail-mark"><img src="./gitscope-logo.png" alt="GitScope" /></div>
            <ul class="rail-list" role="list">
                ${PAGES.map((p) => `
                    <li>
                        <button class="rail-btn${p.id === currentPage ? ' is-active' : ''}" data-page="${p.id}"
                                ${p.id === currentPage ? 'aria-current="page"' : ''}>
                            ${icon(p.icon, 18)}
                            <span class="rail-tip" role="tooltip">${p.label} <kbd>${p.num}</kbd></span>
                        </button>
                    </li>`).join('')}
            </ul>
            <div class="rail-foot">
                <button class="rail-btn" id="railPalette" aria-label="Open command palette">
                    ${icon('search', 18)}
                    <span class="rail-tip" role="tooltip">Commands <kbd>Ctrl K</kbd></span>
                </button>
            </div>
        </nav>

        <div class="workspace">
            <header class="topbar">
                <div class="brand">
                    <img class="brand-logo" src="./gitscope-logo.png" alt="GitScope" />
                    <span class="brand-name">GitScope</span>
                </div>
                <div class="top-sep" aria-hidden="true"></div>
                <div class="repo-slot">
                    <span class="repo-path mono is-empty" id="repoPath">No repository selected</span>
                    <button class="icon-btn" id="copyPathBtn" aria-label="Copy repository path" hidden>${icon('copy', 13)}</button>
                    <span class="branch-chip mono" id="branchChip" hidden>${icon('branch', 12)}<span id="branchName"></span></span>
                    <span class="sync-chip mono" id="aheadChip" hidden></span>
                    <span class="sync-chip mono" id="behindChip" hidden></span>
                </div>
                <div class="top-actions">
                    <button class="btn btn-ghost btn-sm" id="paletteBtn">
                        ${icon('search', 13)}<span>Commands</span><kbd>Ctrl K</kbd>
                    </button>
                    <button class="btn btn-primary btn-sm" id="magicSyncBtn" disabled>
                        ${icon('sync', 13)}<span>Magic Sync</span>
                    </button>
                </div>
            </header>

            <main class="page" id="pageRoot" tabindex="-1"></main>

            <section class="console" id="consolePanel" aria-label="Console">
                <div class="console-grip" id="consoleResize" role="separator" aria-label="Resize console"
                     aria-orientation="horizontal" tabindex="0"></div>
                <div class="console-bar">
                    <span class="micro">${icon('code', 12)} Console</span>
                    <div class="console-bar-right">
                        <span class="console-hint mono">git output</span>
                        <button class="btn btn-ghost btn-sm" id="clearConsoleBtn">${icon('clear', 12)} Clear</button>
                    </div>
                </div>
                <div class="console-out" id="consoleOut" role="log" aria-live="polite" aria-label="Console output"></div>
                <div class="console-status">
                    <span class="led" id="busyLed" aria-hidden="true"></span>
                    <span id="busyText">Idle</span>
                    <span class="console-status-sep" aria-hidden="true">·</span>
                    <span class="mono" id="statusMeta">no repository selected</span>
                </div>
            </section>
        </div>

        <div class="palette-overlay" id="paletteRoot" hidden>
            <div class="palette" role="dialog" aria-modal="true" aria-label="Command palette">
                <div class="palette-input-row">
                    ${icon('search', 15)}
                    <input id="paletteInput" type="text" placeholder="Run a command, jump to a page…"
                           autocomplete="off" spellcheck="false" aria-label="Search commands"
                           aria-controls="paletteList" />
                    <kbd>Esc</kbd>
                </div>
                <ul class="palette-list" id="paletteList" role="listbox" aria-label="Results"></ul>
                <div class="palette-foot">
                    <span><kbd>↑</kbd><kbd>↓</kbd> navigate</span>
                    <span><kbd>Enter</kbd> run</span>
                    <span><kbd>Esc</kbd> close</span>
                </div>
            </div>
        </div>

        <div class="modal-root" id="modalRoot"></div>
        <div class="toast-stack" id="toastRoot"></div>`;
}

function navigate(page) {
    currentPage = page;
    $$('.rail-btn[data-page]').forEach((b) => {
        const on = b.dataset.page === page;
        b.classList.toggle('is-active', on);
        if (on) b.setAttribute('aria-current', 'page');
        else b.removeAttribute('aria-current');
    });
    renderPage(true);
}

function renderPage(focus) {
    const root = $('#pageRoot');
    if (!root) return;
    root.style.animation = 'none';
    void root.offsetWidth;
    root.style.animation = '';
    ({ repo: renderRepo, dashboard: renderDashboard, history: renderHistory,
       docs: renderDocs, about: renderAbout }[currentPage] || renderRepo)(root);
    hydrateIcons(root);
    if (focus) root.focus({ preventScroll: true });
}

/* ─── COMMAND REGISTRY ──────────────────────────────
   Every entry maps 1:1 onto a Go binding. Option strings are the
   exact values the Go switch statements accept. */

const COMMANDS = {};
function reg(key, def) { COMMANDS[key] = def; }
const branchOpts = () => state.branches.slice();

/* ── Common ── */
reg('init', { cat: 'Common', label: 'Init', icon: 'init',
    run: () => run('git init', () => Init(), { refresh: true }) });

reg('stage', { cat: 'Common', label: 'Stage', icon: 'stage',
    open: () => openDialog({
        title: 'Stage Changes', icon: 'stage',
        desc: 'Adds working-tree changes to the index. <code>All (.)</code> includes untracked files; <code>Untracked (-u)</code> updates tracked files only.',
        fields: [{ id: 'mode', label: 'Scope', type: 'select', options: ['All (.)', 'Untracked (-u)'] }],
        actions: [{ id: 'run', label: 'Stage', kind: 'primary',
            onRun: (v) => run('git add ' + v.mode, () => Stage(v.mode), { refresh: true }) }],
    }) });

reg('status', { cat: 'Common', label: 'Status', icon: 'status',
    open: () => openDialog({
        title: 'Git Status', icon: 'status',
        fields: [{ id: 'mode', label: 'Format', type: 'select', options: ['Standard', 'Short (-s)', 'Branch (-b)'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git status', () => Status(v.mode)) }],
    }) });

reg('commit', { cat: 'Common', label: 'Commit', icon: 'commit',
    open: () => openDialog({
        title: 'Git Commit', icon: 'commit',
        fields: [
            { id: 'msg', label: 'Message', type: 'textarea', mono: true, required: true,
              placeholder: 'Describe your change…', errorText: 'Commit message cannot be empty.' },
            { id: 'opt', label: 'Options', type: 'select',
              options: ['Standard (-m)', 'Stage All (-a)', 'Amend (--amend)'] },
        ],
        actions: [{ id: 'run', label: 'Commit', kind: 'primary',
            onRun: (v) => run(`git commit "${v.msg}"`, () => Commit(v.msg, v.opt), { refresh: true }) }],
    }) });

reg('push', { cat: 'Common', label: 'Push', icon: 'push',
    open: () => openDialog({
        title: 'Git Push', icon: 'push',
        desc: 'Pushes to <code>origin</code>. If the remote is ahead, GitScope pulls and retries automatically.',
        fields: [{ id: 'branch', label: 'Branch', type: 'select', options: branchOpts, value: state.branch }],
        actions: [{ id: 'run', label: 'Push', kind: 'primary',
            onRun: (v) => run(`git push -u origin ${v.branch}`, () => Push(v.branch), { refresh: true }) }],
    }) });

reg('log', { cat: 'Common', label: 'Log', icon: 'log',
    open: () => openDialog({
        title: 'Git Log', icon: 'log',
        fields: [{ id: 'mode', label: 'Format', type: 'select', options: ['Oneline', 'Graph', 'Pretty'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git log · ' + v.mode, () => Log(v.mode)) }],
    }) });

/* ── Branches ── */
reg('branch', { cat: 'Branches', label: 'Branch', icon: 'branch',
    open: () => openDialog({
        title: 'Branch', icon: 'branch',
        desc: 'Creates a local branch. Deleting discards unmerged commits.',
        fields: [{ id: 'name', label: 'Branch name', required: true, mono: true, placeholder: 'feature/my-branch' }],
        actions: [
            { id: 'create', label: 'Create', kind: 'primary',
              onRun: (v) => run(`git branch ${v.name}`, () => CreateBranch(v.name), { refresh: true }) },
            { id: 'delete', label: 'Delete', kind: 'danger',
              confirmText: '<strong>Deleting a branch discards unmerged commits.</strong> GitScope cannot restore them.',
              onRun: (v) => run(`git branch -d ${v.name}`, () => DeleteBranch(v.name), { refresh: true }) },
        ],
    }) });

reg('switch', { cat: 'Branches', label: 'Switch', icon: 'branch',
    open: () => openDialog({
        title: 'Switch Branch', icon: 'branch',
        fields: [{ id: 'branch', label: 'Branch', type: 'select', options: branchOpts, value: state.branch }],
        actions: [{ id: 'run', label: 'Switch', kind: 'primary',
            onRun: (v) => run(`git switch ${v.branch}`, () => SwitchBranch(v.branch), { refresh: true }) }],
    }) });

reg('merge', { cat: 'Branches', label: 'Merge', icon: 'merge',
    open: () => openDialog({
        title: 'Merge Branch', icon: 'merge',
        fields: [{ id: 'branch', label: 'Branch to merge', type: 'select',
            options: () => branchOpts().filter((b) => b !== state.branch) }],
        actions: [{ id: 'run', label: 'Merge', kind: 'primary',
            onRun: (v) => run(`git merge ${v.branch}`, () => Merge(v.branch), { refresh: true }) }],
    }) });

reg('rename', { cat: 'Branches', label: 'Rename', icon: 'branch',
    open: () => openDialog({
        title: 'Rename Branch', icon: 'branch',
        fields: [
            { id: 'old', label: 'Current name', type: 'select', options: branchOpts, value: state.branch },
            { id: 'new', label: 'New name', required: true, mono: true, placeholder: 'develop' },
        ],
        actions: [{ id: 'run', label: 'Rename', kind: 'primary',
            onRun: (v) => run(`git branch -m ${v.old} ${v.new}`, () => BranchRename(v.old, v.new), { refresh: true }) }],
    }) });

reg('tag', { cat: 'Branches', label: 'Tag', icon: 'tag',
    open: () => openDialog({
        title: 'Tag', icon: 'tag',
        fields: [
            { id: 'action', label: 'Action', type: 'select', options: ['list', 'create', 'delete', 'push'] },
            { id: 'name', label: 'Tag name', mono: true, placeholder: 'v1.0.0' },
        ],
        danger: { when: (v) => v.action === 'delete',
            text: '<strong>Deleting a tag removes the release marker</strong> locally — and on the remote if it was pushed.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary', onRun: (v) => {
            if (v.action === 'list') { run('git tag', () => Tag('list', '')); return; }
            if (!v.name.trim()) { log('Tag name is required.', 'warning'); return; }
            if (v.action === 'push') { run(`git push origin ${v.name}`, () => Tag('push', v.name)); return; }
            run(`git tag ${v.action === 'delete' ? '-d ' : ''}${v.name}`, () => Tag(v.action, v.name));
        } }],
    }) });

/* ── Remote ── */
reg('remote', { cat: 'Remote', label: 'Remote', icon: 'remote',
    open: () => openDialog({
        title: 'Remote', icon: 'remote',
        fields: [
            { id: 'action', label: 'Action', type: 'select', options: ['list', 'add', 'remove'] },
            { id: 'val', label: 'URL / remote name', mono: true, placeholder: 'https://github.com/user/repo.git' },
        ],
        danger: { when: (v) => v.action === 'remove',
            text: '<strong>Removing a remote drops its fetch/push endpoints.</strong> Local branches and commits are kept.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary', onRun: (v) => {
            if (v.action === 'list') { run('git remote -v', () => Remote('list', '')); return; }
            if (!v.val.trim()) { log(v.action === 'add' ? 'URL is required.' : 'Remote name is required.', 'warning'); return; }
            run(`git remote ${v.action} ${v.val}`, () => Remote(v.action, v.val), { refresh: true });
        } }],
    }) });

reg('fetch', { cat: 'Remote', label: 'Fetch', icon: 'fetch',
    open: () => openDialog({
        title: 'Git Fetch', icon: 'fetch',
        fields: [
            { id: 'mode', label: 'Scope', type: 'select', options: ['Default', 'All (--all)'] },
            { id: 'remoteSel', label: 'Remote name (optional)', mono: true, placeholder: 'origin' },
        ],
        actions: [{ id: 'run', label: 'Fetch', kind: 'primary', onRun: (v) => {
            const opt = v.remoteSel.trim() ? v.remoteSel.trim() : v.mode;
            run('git fetch ' + opt, () => Fetch(opt), { refresh: true });
        } }],
    }) });

reg('pull', { cat: 'Remote', label: 'Pull', icon: 'pull',
    open: () => openDialog({
        title: 'Git Pull', icon: 'pull',
        fields: [{ id: 'branch', label: 'Branch', type: 'select', options: branchOpts, value: state.branch }],
        actions: [{ id: 'run', label: 'Pull', kind: 'primary',
            onRun: (v) => run(`git pull origin ${v.branch}`, () => Pull(v.branch), { refresh: true }) }],
    }) });

reg('clone', { cat: 'Remote', label: 'Clone', icon: 'clone',
    open: () => openDialog({
        title: 'Clone Repository', icon: 'clone',
        desc: 'Clones into the currently selected folder.',
        fields: [{ id: 'url', label: 'Repository URL', required: true, mono: true,
            placeholder: 'https://github.com/user/repo.git', errorText: 'A clone URL is required.' }],
        actions: [{ id: 'run', label: 'Clone', kind: 'primary',
            onRun: (v) => run(`git clone ${v.url}`, () => Clone(v.url), { refresh: true }) }],
    }) });

reg('cherry', { cat: 'Remote', label: 'Cherry-pick', icon: 'cherry',
    open: () => openDialog({
        title: 'Cherry-pick', icon: 'cherry',
        fields: [{ id: 'hash', label: 'Commit hash', required: true, mono: true,
            placeholder: 'a1b2c3d', errorText: 'A commit hash is required.' }],
        actions: [{ id: 'run', label: 'Apply', kind: 'primary',
            onRun: (v) => run(`git cherry-pick ${v.hash}`, () => CherryPick(v.hash), { refresh: true }) }],
    }) });

/* ── History ── */
reg('log-history', { cat: 'History', label: 'Log', icon: 'log',
    open: () => COMMANDS.log.open() });

reg('revert', { cat: 'History', label: 'Revert', icon: 'revert',
    open: () => openDialog({
        title: 'Revert Commit', icon: 'revert',
        desc: 'Creates a new commit that undoes the given commit (<code>--no-edit</code>).',
        fields: [{ id: 'hash', label: 'Commit hash', required: true, mono: true,
            placeholder: 'a1b2c3d', errorText: 'A commit hash is required.' }],
        actions: [{ id: 'run', label: 'Revert', kind: 'danger',
            confirmText: '<strong>Revert writes a new commit.</strong> The history is extended, not rewritten.',
            onRun: (v) => run(`git revert ${v.hash}`, () => Revert(v.hash), { refresh: true }) }],
    }) });

reg('show', { cat: 'History', label: 'Show', icon: 'show',
    open: () => openDialog({
        title: 'Git Show', icon: 'show',
        desc: 'Shows the current HEAD. To inspect a specific commit, open the History page and click it.',
        fields: [{ id: 'mode', label: 'Detail', type: 'select',
            options: ['Full', 'Stats (--stat)', 'Name Only (--name-only)', 'Patch (--patch)'] }],
        actions: [{ id: 'run', label: 'Show', kind: 'primary',
            onRun: (v) => run('git show ' + v.mode, () => Show(v.mode, '')) }],
    }) });

reg('shortlog', { cat: 'History', label: 'Shortlog', icon: 'shortlog',
    open: () => openDialog({
        title: 'Git Shortlog', icon: 'shortlog',
        fields: [{ id: 'mode', label: 'Format', type: 'select',
            options: ['Summary (-s)', 'Numeric (-n)', 'By email (-e)', 'All (-a)'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git shortlog ' + v.mode, () => Shortlog(v.mode)) }],
    }) });

reg('reflog', { cat: 'History', label: 'Reflog', icon: 'log',
    open: () => openDialog({
        title: 'Git Reflog', icon: 'log',
        desc: 'Every position HEAD has pointed at — the safety net for lost commits.',
        fields: [{ id: 'mode', label: 'Format', type: 'select',
            options: ['Standard', 'Limit 20 (-n 20)', 'Date Relative'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git reflog ' + v.mode, () => Reflog(v.mode)) }],
    }) });

/* ── Changes ── */
reg('diff', { cat: 'Changes', label: 'Diff', icon: 'diff',
    open: () => openDialog({
        title: 'Git Diff', icon: 'diff',
        fields: [{ id: 'mode', label: 'Mode', type: 'select',
            options: ['Unstaged', 'Staged (--cached)', 'Names (--name-only)', 'Summary (--stat)'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git diff ' + v.mode, () => Diff(v.mode)) }],
    }) });

reg('stash', { cat: 'Changes', label: 'Stash', icon: 'stash',
    open: () => openDialog({
        title: 'Git Stash', icon: 'stash',
        desc: '<code>Save</code> shelves current changes; <code>Pop</code> restores and drops the newest stash.',
        fields: [{ id: 'action', label: 'Action', type: 'select', options: ['Save', 'Pop', 'List', 'Drop', 'Apply'] }],
        danger: { when: (v) => v.action === 'Drop',
            text: '<strong>Dropping a stash deletes the shelved changes permanently.</strong> GitScope cannot recover them.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run(`git stash ${v.action.toLowerCase()}`, () => Stash(v.action), { refresh: true }) }],
    }) });

reg('clean', { cat: 'Changes', label: 'Clean', icon: 'clean', danger: true,
    open: () => openDialog({
        title: 'Git Clean', icon: 'clean',
        desc: 'Removes untracked files. <code>Dry Run (-n)</code> only lists what would be deleted.',
        fields: [{ id: 'mode', label: 'Mode', type: 'select',
            options: ['Dry Run (-n)', 'Directories (-d)', 'Force (-f)', 'Full (-fdx)'] }],
        danger: { when: (v) => v.mode !== 'Dry Run (-n)',
            text: '<strong>Untracked files are deleted permanently.</strong> Git cannot restore them afterwards.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary',
            onRun: (v) => run('git clean ' + v.mode, () => Clean(v.mode), { refresh: true }) }],
    }) });

reg('lsfiles', { cat: 'Changes', label: 'Ls-Files', icon: 'lsfiles',
    open: () => openDialog({
        title: 'Git Ls-Files', icon: 'lsfiles',
        fields: [{ id: 'mode', label: 'Filter', type: 'select',
            options: ['Tracked', 'Cached (--cached)', 'Modified (--modified)', 'Others (--others)', 'Deleted (--deleted)', 'Staged'] }],
        actions: [{ id: 'run', label: 'Run', kind: 'primary', onRun: (v) => {
            const opt = v.mode === 'Tracked' ? 'Default' : v.mode;
            run('git ls-files ' + v.mode, () => LsFiles(opt));
        } }],
    }) });

reg('gitignore', { cat: 'Changes', label: '.gitignore', icon: 'gitignore',
    open: async () => {
        const content = await safe(ReadGitIgnore, '');
        openDialog({
            title: '.gitignore', icon: 'gitignore', wide: true,
            desc: 'Read from and written to the repository root.',
            fields: [{ id: 'content', label: 'Contents', type: 'textarea', mono: true,
                value: content, placeholder: 'build/\n*.log' }],
            actions: [{ id: 'save', label: 'Save', kind: 'primary', onRun: async (v) => {
                try {
                    await WriteGitIgnore(v.content);
                    log('.gitignore saved.', 'success');
                    await loadRepo();
                } catch (err) {
                    log(errMsg(err), 'error');
                }
            } }],
        });
    } });

/* ── Advanced ── */
reg('reset', { cat: 'Advanced', label: 'Reset', icon: 'reset', danger: true,
    open: () => openDialog({
        title: 'Git Reset', icon: 'reset',
        desc: 'Moves HEAD. <code>--soft</code> keeps changes staged, <code>--mixed</code> unstages them, <code>--hard</code> discards them.',
        fields: [
            { id: 'mode', label: 'Mode', type: 'select', options: ['--mixed', '--soft', '--hard'] },
            { id: 'target', label: 'Target', mono: true, value: 'HEAD~1' },
        ],
        danger: { when: (v) => v.mode === '--hard',
            text: '<strong>--hard discards uncommitted work in the working tree.</strong> It cannot be undone.' },
        actions: [{ id: 'run', label: 'Reset', kind: 'primary',
            onRun: (v) => run(`git reset ${v.mode} ${v.target}`, () => Reset(v.mode, v.target), { refresh: true }) }],
    }) });

reg('rebase', { cat: 'Advanced', label: 'Rebase', icon: 'rebase', danger: true,
    open: () => openDialog({
        title: 'Git Rebase', icon: 'rebase',
        desc: 'Replays commits on a new base. <code>-i</code> uses the interactive editor; <code>--onto</code> takes the target as the new base.',
        fields: [
            { id: 'option', label: 'Operation', type: 'select',
                options: ['-i', '--onto', 'Continue', 'Abort', 'Skip'] },
            { id: 'target', label: 'Target (for -i / --onto)', mono: true, placeholder: 'main' },
        ],
        danger: { when: (v) => v.option === '-i' || v.option === '--onto',
            text: '<strong>Rebase rewrites commit history.</strong> A force-push is required to update the remote.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary', onRun: (v) => {
            const needsTarget = v.option === '-i' || v.option === '--onto';
            if (needsTarget && !v.target.trim()) { log('A target revision is required for ' + v.option + '.', 'warning'); return; }
            run(`git rebase ${v.option}${needsTarget ? ' ' + v.target : ''}`,
                () => Rebase(v.option, needsTarget ? v.target : ''), { refresh: true });
        } }],
    }) });

reg('undo', { cat: 'Advanced', label: 'Undo', icon: 'undo', danger: true,
    open: () => openDialog({
        title: 'Undo Last Commit', icon: 'undo',
        desc: 'Runs <code>git reset --soft HEAD~1</code>: the commit is removed and its changes stay staged.',
        danger: { when: () => true,
            text: '<strong>The last commit is removed from the branch.</strong> Its changes remain staged, so nothing is lost.' },
        actions: [{ id: 'run', label: 'Undo commit', kind: 'primary',
            onRun: () => run('git reset --soft HEAD~1', () => UndoLastCommit(), { refresh: true }) }],
    }) });

reg('worktree', { cat: 'Advanced', label: 'Worktree', icon: 'worktree',
    open: () => openDialog({
        title: 'Worktree', icon: 'worktree',
        fields: [
            { id: 'action', label: 'Action', type: 'select', options: ['List', 'Add', 'Remove', 'Prune'] },
            { id: 'path', label: 'Path (Add / Remove)', mono: true, placeholder: 'C:\\dev\\repo-wt' },
            { id: 'branch', label: 'Branch (Add only)', mono: true, placeholder: 'feature/x' },
        ],
        danger: { when: (v) => v.action === 'Remove',
            text: '<strong>Removing a worktree deletes its directory.</strong> Uncommitted work inside it is lost.' },
        actions: [{ id: 'run', label: 'Run', kind: 'primary', onRun: (v) => {
            if (v.action === 'Add' && (!v.path.trim() || !v.branch.trim())) {
                log('Add needs both a path and a branch.', 'warning'); return;
            }
            if (v.action === 'Remove' && !v.path.trim()) {
                log('Remove needs the worktree path.', 'warning'); return;
            }
            const args = v.action === 'Add' ? `${v.path.trim()} ${v.branch.trim()}`
                : v.action === 'Remove' ? v.path.trim() : '';
            run(`git worktree ${v.action.toLowerCase()} ${args}`.trim(), () => Worktree(v.action, args));
        } }],
    }) });

reg('conflicts', { cat: 'Advanced', label: 'Conflicts', icon: 'conflict',
    open: async () => {
        const conflicts = sanitizeConflicts(await safe(GetConflicts, []));
        state.conflicts = conflicts;
        if (!conflicts.length) {
            log('No merge conflicts detected.', 'success');
            return;
        }
        const rows = conflicts.map((f) =>
            `<div class="conflict-row"><span class="conflict-path">${esc(f)}</span>` +
            `<button type="button" class="btn btn-secondary btn-sm" data-file="${esc(f)}" data-side="ours">Keep Mine</button>` +
            `<button type="button" class="btn btn-danger btn-sm" data-file="${esc(f)}" data-side="theirs">Take Theirs</button></div>`).join('');
        const dlg = openDialog({
            title: 'Resolve Conflicts', icon: 'conflict', wide: true,
            desc: `${conflicts.length} file(s) are unmerged. Resolving checks out one side and stages the file.`,
            bodyHTML: rows,
            actions: [{ id: 'close', label: 'Close', kind: 'secondary', onRun: () => {} }],
        });
        dlg.addEventListener('click', (e) => {
            const btn = e.target.closest('[data-file]');
            if (!btn) return;
            const file = btn.dataset.file;
            const side = btn.dataset.side;
            closeDialog();
            run(`git checkout --${side} ${file} && git add ${file}`,
                async () => {
                    const out = await ResolveConflict(file, side);
                    await loadRepo();
                    return out;
                }, { needsRepo: true });
        });
    } });

/* ── Tools ── */
reg('blame', { cat: 'Tools', label: 'Blame', icon: 'blame',
    open: () => openDialog({
        title: 'Git Blame', icon: 'blame',
        fields: [{ id: 'file', label: 'File path', required: true, mono: true,
            placeholder: 'src/main.go', errorText: 'A file path is required.' }],
        actions: [{ id: 'run', label: 'Blame', kind: 'primary',
            onRun: (v) => run(`git blame ${v.file}`, () => Blame(v.file)) }],
    }) });

reg('magic', { cat: 'Tools', label: 'Magic Sync', icon: 'sync',
    run: () => run('git stash push -u · git fetch · git pull --rebase · git stash pop', () => MagicSync(), { refresh: true }) });

/* category order + which tiles get the danger treatment */
const CATEGORY_ORDER = ['Common', 'Branches', 'Remote', 'History', 'Changes', 'Advanced', 'Tools'];
const DANGER_KEYS = new Set(['reset', 'rebase', 'undo', 'clean', 'revert', 'branch', 'tag']);

/* Open a command's dialog, or execute it directly. */
function exec(key) {
    const cmd = COMMANDS[key];
    if (!cmd) { log('Unknown command: ' + key, 'error'); return; }
    closePalette();
    if (!state.repoPath && key !== 'init') { log('Select a repository first.', 'warning'); return; }
    if (cmd.open) cmd.open();
    else cmd.run();
}

/* ─── REPOSITORY PAGE ─────────────────────────────── */

function fileList(items, kind) {
    if (!items.length) {
        return `<li class="empty-note">${kind === 'staged'
            ? 'Nothing staged — use Stage or Stage all.'
            : 'Working tree clean.'}</li>`;
    }
    const seen = new Set();
    return items.filter((it) => {
        const sig = it.code + it.file;
        if (seen.has(sig)) return false;
        seen.add(sig);
        return true;
    }).map((it) => `
        <li class="file-row">
            <span class="file-code c-${['M', 'A', 'D', 'R'].includes(it.code) ? it.code : 'U'}" aria-hidden="true">${esc(it.code)}</span>
            <span class="file-path" title="${esc(it.file)}">${esc(it.file)}</span>
        </li>`).join('');
}

function setupPanelHTML() {
    return `
        <section class="panel" aria-label="Quick setup">
            <div class="panel-head"><span class="panel-title">${icon('folder', 14)} Quick setup</span></div>
            <div class="panel-body setup-list">
                <div class="setup-card">
                    <div class="setup-card-head">
                        <span class="setup-icon">${icon('stage', 15)}</span>
                        <div><h4>Create new repository</h4><p>Commands run inside the selected folder</p></div>
                    </div>
                    <textarea id="newRepoCmds" class="mono" aria-label="New repository commands"
                        placeholder='git init&#10;git add .&#10;git commit -m "initial commit"'></textarea>
                    <button type="button" class="btn btn-primary btn-sm" id="runNewRepoBtn">${icon('play', 12)} Run</button>
                </div>
                <div class="setup-card">
                    <div class="setup-card-head">
                        <span class="setup-icon">${icon('remote', 15)}</span>
                        <div><h4>Push existing repository</h4><p>Connect this folder to a remote</p></div>
                    </div>
                    <textarea id="existRepoCmds" class="mono" aria-label="Existing repository commands"
                        placeholder="git remote add origin https://github.com/user/repo.git&#10;git branch -M main&#10;git push -u origin main"></textarea>
                    <button type="button" class="btn btn-primary btn-sm" id="runExistRepoBtn">${icon('play', 12)} Run</button>
                </div>
            </div>
            <p class="panel-note">Runs each line through <span class="mono">RunCommands()</span> — <span class="mono">git …</span> lines go straight to git, others run as shell commands.</p>
        </section>`;
}

function renderRepo(root) {
    const head = `
        <header class="page-head">
            <div>
                <p class="micro page-kicker">${icon('repo', 12)} Working tree</p>
                <h1 class="page-title">Repository</h1>
                <p class="page-sub">${state.repoPath
                    ? 'Status, staging and setup for this working copy.'
                    : 'Select a folder to begin — everything on this page is read from git.'}</p>
            </div>
            <div class="page-actions">
                <button class="btn btn-secondary" id="selectRepoBtn">${icon('select', 14)} Select repository</button>
                <button class="btn btn-ghost" id="openFolderBtn">${icon('folder', 14)} Open folder</button>
                <button class="btn btn-ghost" id="refreshRepoBtn">${icon('refresh', 14)} Refresh</button>
            </div>
        </header>`;

    /* git missing — nothing else can work */
    if (!state.gitAvailable) {
        root.innerHTML = head + `
            <div class="banner is-error">${icon('xCircle', 16)}
                <div><strong>Git was not found on PATH.</strong> GitScope needs the git executable to read repositories.
                Install Git, then press Refresh.</div>
                <div class="banner-actions"><button class="btn btn-secondary btn-sm" id="refreshRepoBtn2">Refresh</button></div>
            </div>`;
        wireRepo(root);
        return;
    }

    /* no folder selected yet */
    if (!state.repoPath) {
        root.innerHTML = head + `
            <div class="empty-state">
                <div class="empty-icon">${icon('repo', 26)}</div>
                <h2>No repository selected</h2>
                <p>Pick a folder that contains a Git repository — or an empty folder you want to turn into one.</p>
                <div class="empty-actions">
                    <button class="btn btn-primary" id="selectRepoBtn2">${icon('select', 14)} Select repository</button>
                    <button class="btn btn-secondary" id="initHintBtn">${icon('init', 14)} Initialize a folder</button>
                </div>
            </div>`;
        wireRepo(root);
        return;
    }

    /* folder selected but not a git repository yet */
    if (!state.initialized) {
        root.innerHTML = head + `
            <div class="banner is-warn">${icon('warn', 16)}
                <div><strong>Not a Git repository yet.</strong> <code>${esc(state.repoPath)}</code> has no <code>.git</code> directory.</div>
                <div class="banner-actions">
                    <button class="btn btn-primary btn-sm" id="initRepoBtn">${icon('init', 13)} Initialize</button>
                </div>
            </div>
            <div class="repo-grid"><div class="repo-col">${setupPanelHTML()}</div></div>`;
        wireRepo(root);
        return;
    }

    /* full view — repository is initialised and readable */
    const syncUnknown = state.ahead === null && state.behind === null;
    root.innerHTML = head + `
        <div class="stat-row">
            <div class="stat-card">
                <span class="micro stat-label">${icon('branch', 12)} Branch</span>
                <span class="stat-value accent" title="${esc(state.branch)}">${esc(state.branch || 'detached HEAD')}</span>
                <span class="stat-note">${state.upstream ? 'tracking ' + esc(state.upstream) : 'no upstream configured'}</span>
            </div>
            <div class="stat-card${syncUnknown ? ' is-quiet' : ''}">
                <span class="micro stat-label">${icon('push', 12)} Sync</span>
                <span class="stat-value">${syncUnknown
                    ? '<span class="muted">no upstream</span>'
                    : `<span class="up">↑${state.ahead ?? 0}</span> <span class="down">↓${state.behind ?? 0}</span>`}</span>
                <span class="stat-note">${syncUnknown ? 'push once to set tracking' : 'relative to ' + esc(state.upstream)}</span>
            </div>
            <div class="stat-card">
                <span class="micro stat-label">${icon('stage', 12)} Working tree</span>
                <span class="stat-value">${state.staged.length} <span class="muted">staged</span> · ${state.changed.length} <span class="muted">changed</span></span>
                <span class="stat-note">${state.changed.length ? 'review before committing' : 'nothing to commit'}</span>
            </div>
            <div class="stat-card${state.conflicts.length ? ' is-alert' : ' is-quiet'}">
                <span class="micro stat-label">${icon('conflict', 12)} Conflicts</span>
                <span class="stat-value">${state.conflicts.length
                    ? `<span class="alert">${state.conflicts.length}</span>`
                    : '<span class="muted">none</span>'}</span>
                <span class="stat-note">${state.conflicts.length ? 'needs resolution' : 'no unmerged files'}</span>
            </div>
        </div>
        ${state.conflicts.length ? `
        <div class="banner is-error">${icon('conflict', 16)}
            <div><strong>${state.conflicts.length} unresolved file(s).</strong> Resolve them before committing or merging.</div>
            <div class="banner-actions"><button class="btn btn-primary btn-sm" id="resolveConflictsBtn">Resolve conflicts</button></div>
        </div>` : ''}
        <div class="repo-grid">
            <section class="panel" aria-label="Changes">
                <form class="composer" id="commitComposer" novalidate>
                    <label class="field-label" for="commitMsg">Commit message</label>
                    <textarea id="commitMsg" class="mono" placeholder="Describe your change…"
                              aria-required="true" aria-describedby="commitMsgErr"></textarea>
                    <p class="field-error" id="commitMsgErr" role="alert"></p>
                    <div class="composer-bar">
                        <div class="field-inline">
                            <label class="field-label" for="commitOpt" style="margin:0">Options</label>
                            <select id="commitOpt">
                                <option>Standard (-m)</option>
                                <option>Stage All (-a)</option>
                                <option>Amend (--amend)</option>
                            </select>
                        </div>
                        <span class="spacer"></span>
                        <span class="composer-hint">${state.staged.length} staged</span>
                        <button type="submit" class="btn btn-primary btn-sm">${icon('commit', 13)} Commit</button>
                    </div>
                </form>

                <div class="file-group">
                    <div class="file-group-head">
                        <span class="dot dot-staged"></span>
                        <span class="micro">Staged</span>
                        <span class="panel-count">${state.staged.length}</span>
                        <span class="spacer"></span>
                        <button type="button" class="btn btn-ghost btn-sm" id="stageAllBtn">Stage all</button>
                    </div>
                    <ul class="file-list">${fileList(state.staged, 'staged')}</ul>
                </div>

                <div class="file-group">
                    <div class="file-group-head">
                        <span class="dot dot-unstaged"></span>
                        <span class="micro">Changed</span>
                        <span class="panel-count">${state.changed.length}</span>
                        <span class="spacer"></span>
                        <button type="button" class="btn btn-ghost btn-sm" id="diffBtn">Diff</button>
                    </div>
                    <ul class="file-list">${fileList(state.changed, 'changed')}</ul>
                </div>

                <p class="panel-note">Read live from <span class="mono">git status -s</span>. Per-file stage/unstage needs a new
                   backend binding — <span class="mono">Stage()</span> accepts <span class="mono">All (.)</span> or
                   <span class="mono">Untracked (-u)</span> today.</p>
            </section>

            <aside class="repo-col">
                ${setupPanelHTML()}
            </aside>
        </div>`;
    wireRepo(root);
}

function wireRepo(root) {
    const byId = (id) => root.querySelector('#' + id);

    /* shared in every state of this page */
    byId('selectRepoBtn')?.addEventListener('click', pickFolder);
    byId('selectRepoBtn2')?.addEventListener('click', pickFolder);
    byId('openFolderBtn')?.addEventListener('click', pickFolder);
    byId('refreshRepoBtn')?.addEventListener('click', () => loadRepo());
    byId('refreshRepoBtn2')?.addEventListener('click', () => loadRepo());
    byId('initRepoBtn')?.addEventListener('click', () => exec('init'));
    byId('initHintBtn')?.addEventListener('click', pickFolder);
    byId('runNewRepoBtn')?.addEventListener('click', () => runSetup('runNewRepoBtn', 'newRepoCmds'));
    byId('runExistRepoBtn')?.addEventListener('click', () => runSetup('runExistRepoBtn', 'existRepoCmds'));

    byId('stageAllBtn')?.addEventListener('click', () => exec('stage'));
    byId('diffBtn')?.addEventListener('click', () => exec('diff'));
    byId('resolveConflictsBtn')?.addEventListener('click', () => exec('conflicts'));

    byId('commitComposer')?.addEventListener('submit', (e) => {
        e.preventDefault();
        const ta = byId('commitMsg');
        const err = byId('commitMsgErr');
        const msg = ta.value.trim();
        if (!msg) {
            ta.setAttribute('aria-invalid', 'true');
            err.textContent = 'Commit message cannot be empty.';
            err.classList.add('is-visible');
            ta.focus();
            toast('Commit message cannot be empty.', 'error');
            return;
        }
        ta.removeAttribute('aria-invalid');
        err.classList.remove('is-visible');
        const opt = byId('commitOpt').value;
        run(`git commit "${msg}"`, () => Commit(msg, opt), { refresh: true });
    });
}

async function pickFolder() {
    try {
        const path = await SelectRepo();
        if (path) {
            log('Repository: ' + path, 'success');
            const draft = captureDraft();
            await loadRepo();
            renderPage();
            restoreDraft(draft);
        }
    } catch (err) {
        const m = errMsg(err);
        if (!/no directory selected|canceled|cancelled/i.test(m)) log(m, 'error');
    }
}

async function runSetup(btnId, textareaId) {
    const btn = $('#' + btnId);
    const ta = $('#' + textareaId);
    if (!ta) return;
    const text = ta.value.trim();
    if (!text) { log('Paste at least one command first.', 'warning'); return; }
    if (btn) btn.disabled = true;
    try {
        const out = await RunCommands(text);
        String(out || '').split('\n').forEach((line) => {
            if (!line) return;
            if (line.startsWith('>')) log(line, 'cmd');
            else if (/^Error:/i.test(line)) log(line, 'error');
            else log(line, 'default');
        });
        const draft = captureDraft();
        await loadRepo();
        renderPage();
        restoreDraft(draft);
    } catch (err) {
        log(errMsg(err), 'error');
    } finally {
        if (btn) btn.disabled = false;
    }
}

/* ─── DASHBOARD PAGE ──────────────────────────────── */

function renderDashboard(root) {
    const keys = Object.keys(COMMANDS);
    root.innerHTML = `
        <header class="page-head">
            <div>
                <p class="micro page-kicker">${icon('dashboard', 12)} Command center</p>
                <h1 class="page-title">Dashboard</h1>
                <p class="page-sub">${keys.length} git operations across ${CATEGORY_ORDER.length} groups — every tile calls a
                   Go binding. Or press <kbd>Ctrl K</kbd> to search them all.</p>
            </div>
        </header>
        <div class="dash-layout">
            <nav class="dash-cats" aria-label="Command categories">
                ${CATEGORY_ORDER.map((c) => {
                    const n = keys.filter((k) => COMMANDS[k].cat === c).length;
                    return `<button type="button" class="cat-btn" data-cat="${c}"
                        aria-pressed="${c === currentCategory}">${c}<span class="cat-n">${n}</span></button>`;
                }).join('')}
            </nav>
            <div class="dash-main" id="dashMain"></div>
        </div>`;

    renderCat(root.querySelector('#dashMain'), currentCategory);
    root.querySelectorAll('.cat-btn').forEach((b) => {
        b.addEventListener('click', () => {
            currentCategory = b.dataset.cat;
            root.querySelectorAll('.cat-btn').forEach((x) =>
                x.setAttribute('aria-pressed', String(x.dataset.cat === currentCategory)));
            renderCat(root.querySelector('#dashMain'), currentCategory);
            hydrateIcons(root);
        });
    });
}

function renderCat(container, cat) {
    const keys = Object.keys(COMMANDS).filter((k) => COMMANDS[k].cat === cat);
    container.innerHTML = `
        <div class="dash-section-head">
            <h2>${cat}</h2>
            <span class="count">${keys.length} command${keys.length === 1 ? '' : 's'}</span>
            <span class="rule"></span>
        </div>
        <div class="cmd-grid">
            ${keys.map((k, i) => {
                const c = COMMANDS[k];
                return `<button type="button" class="cmd-card${DANGER_KEYS.has(k) ? ' is-danger' : ''}"
                    data-cmd="${k}" style="animation:rowIn .3s ease ${i * 0.03}s both">
                    <span class="cmd-ico">${icon(c.icon, 16)}</span>
                    ${esc(c.label)}
                </button>`;
            }).join('')}
        </div>
        ${cat === 'Advanced' ? `<p class="panel-note" style="border:1px solid var(--line-1);border-radius:10px;margin-top:14px">
            Destructive tiles require an explicit confirmation step before they run.</p>` : ''}`;
    container.querySelectorAll('[data-cmd]').forEach((b) =>
        b.addEventListener('click', () => exec(b.dataset.cmd)));
}

/* ─── HISTORY PAGE (real `git log`) ───────────────── */

/* Oneline:  <hash> <subject>  ·  Pretty:  <hash> - <meta> : <subject> */
function parseLog(text, format) {
    const rows = [];
    const rawLines = [];
    String(text || '').split('\n').forEach((line) => {
        if (!line.trim()) return;
        rawLines.push(line);
        if (format === 'Oneline') {
            const m = /^([0-9a-f]{7,40})\s+(.*)$/.exec(line);
            if (m) rows.push({ hash: m[1], subject: m[2] });
        } else if (format === 'Pretty') {
            const m = /^([0-9a-f]{7,40})\s+-\s+(.*)$/.exec(line);
            if (m) rows.push({ hash: m[1], subject: m[2] });
        }
    });
    return { rows, rawLines };
}

async function loadHistory(format) {
    if (!state.repoPath) {
        state.history.raw = '';
        state.history.rows = [];
        state.history.loaded = false;
        return;
    }
    state.history.format = format;
    setBusy(true);
    try {
        const out = await Log(format);
        state.history.raw = String(out || '');
        const parsed = parseLog(state.history.raw, format);
        state.history.rows = parsed.rows;
        state.history.rawLines = parsed.rawLines;
        state.history.loaded = true;
    } catch (err) {
        log(errMsg(err), 'error');
        state.history.raw = '';
        state.history.rows = [];
        state.history.loaded = false;
    } finally {
        setBusy(false);
    }
}

function renderHistory(root) {
    root.innerHTML = `
        <header class="page-head">
            <div>
                <p class="micro page-kicker">${icon('history', 12)} git log</p>
                <h1 class="page-title">History</h1>
                <p class="page-sub">Commit history for <em>${esc(state.branch || 'HEAD')}</em> — click a commit to inspect
                   or revert it.</p>
            </div>
            <div class="page-actions">
                <button class="btn btn-secondary" id="reloadHistoryBtn">${icon('refresh', 14)} Reload</button>
                <button class="btn btn-ghost" id="historyConsoleBtn">${icon('code', 14)} Send to console</button>
            </div>
        </header>

        ${!state.repoPath ? `
            <div class="empty-state">
                <div class="empty-icon">${icon('history', 26)}</div>
                <h2>No repository selected</h2>
                <p>Select a repository on the Repository page to load its commit history.</p>
                <div class="empty-actions"><button class="btn btn-primary" id="historyGoRepo">Go to Repository</button></div>
            </div>` : `
            <div class="history-toolbar">
                <label class="micro" for="historyFormat">${icon('log', 12)} Format</label>
                <select id="historyFormat">
                    <option${state.history.format === 'Oneline' ? ' selected' : ''}>Oneline</option>
                    <option${state.history.format === 'Graph' ? ' selected' : ''}>Graph</option>
                    <option${state.history.format === 'Pretty' ? ' selected' : ''}>Pretty</option>
                </select>
                <label class="micro" for="historyFilter">${icon('search', 12)} Filter</label>
                <input type="text" id="historyFilter" placeholder="hash or subject…" autocomplete="off" />
                <span class="spacer"></span>
                <span class="micro" id="historyCount">${state.history.rows.length} commits</span>
            </div>
            <div id="historyBody"></div>`}`;

    hydrateIcons(root);
    if (!state.repoPath) {
        root.querySelector('#historyGoRepo')?.addEventListener('click', () => navigate('repo'));
        return;
    }

    const body = root.querySelector('#historyBody');
    const count = root.querySelector('#historyCount');

    const draw = (filter) => {
        const q = (filter || '').toLowerCase().trim();
        const rows = state.history.rows.filter((r) =>
            !q || r.hash.toLowerCase().includes(q) || r.subject.toLowerCase().includes(q));
        if (count) count.textContent = rows.length + ' commit' + (rows.length === 1 ? '' : 's');

        if (!state.history.loaded) {
            body.innerHTML = `<div class="panel"><div class="empty-note" style="padding:34px">Loading history…</div></div>`;
            return;
        }
        /* graph output keeps its layout, so show it verbatim */
        if (state.history.format === 'Graph' || (state.history.rows.length === 0 && state.history.raw)) {
            body.innerHTML = `<div class="panel"><div class="log-raw">${esc(state.history.raw) || 'No commits yet.'}</div></div>`;
            return;
        }
        if (!rows.length) {
            body.innerHTML = `<div class="panel"><div class="empty-note" style="padding:34px">${
                state.history.raw ? 'No commits match that filter.' : 'No commits yet — make your first commit.'}</div></div>`;
            return;
        }
        body.innerHTML = `<div class="panel"><ul class="log-list">
            ${rows.map((r, i) => `
                <li><button type="button" class="log-row${i === 0 && !q ? ' is-head' : ''}" data-hash="${esc(r.hash)}"
                    style="animation:rowIn .28s ease ${Math.min(i, 20) * 0.03}s both">
                    <span class="log-node" aria-hidden="true"></span>
                    <span class="log-hash">${esc(r.hash)}</span>
                    <span class="log-msg" title="${esc(r.subject)}">${esc(r.subject)}</span>
                </button></li>`).join('')}
        </ul></div>`;
        body.querySelectorAll('[data-hash]').forEach((b) =>
            b.addEventListener('click', () => openCommitDialog(b.dataset.hash)));
    };

    const reload = async (format) => {
        await loadHistory(format);
        draw(root.querySelector('#historyFilter')?.value || '');
    };

    if (!state.history.loaded) reload(state.history.format);
    else draw('');

    root.querySelector('#historyFormat')?.addEventListener('change', (e) => reload(e.target.value));
    root.querySelector('#historyFilter')?.addEventListener('input', (e) => draw(e.target.value));
    root.querySelector('#reloadHistoryBtn')?.addEventListener('click', () =>
        reload(root.querySelector('#historyFormat').value));
    root.querySelector('#historyConsoleBtn')?.addEventListener('click', () => {
        if (state.history.raw) log(state.history.raw, 'default');
        else log('Run Reload first.', 'warning');
    });
}

function openCommitDialog(hash) {
    openDialog({
        title: 'Commit ' + hash, icon: 'show', wide: true,
        desc: `<span class="mono" style="color:var(--signal)">${esc(hash)}</span> — <code>Show</code> prints
               <span class="mono">git show ${esc(hash)}</span> to the console.`,
        actions: [
            { id: 'show', label: 'Show in console', kind: 'primary',
              onRun: () => run(`git show ${hash}`, () => Show('Specific', hash)) },
            { id: 'copy', label: 'Copy hash', kind: 'secondary',
              onRun: () => { copyText(hash); toast('Copied ' + hash, 'success'); } },
            { id: 'revert', label: 'Revert', kind: 'danger',
              confirmText: '<strong>Revert writes a new commit that undoes ' + esc(hash) + '.</strong>',
              onRun: () => run(`git revert ${hash}`, () => Revert(hash), { refresh: true }) },
        ],
    });
}

/* ─── DOCS PAGE ───────────────────────────────────── */

const DOC_ITEMS = ['Init', 'Stage', 'Status', 'Commit', 'Push', 'Log', 'Revert', 'Clone', 'Branch', 'Pull',
    'Reflog', 'GitIgnore', 'Remote', 'Diff', 'Reset', 'Fetch', 'Stash', 'Merge', 'Tag',
    'Cherry-pick', 'Rebase', 'Clean', 'Show', 'Ls-files', 'Worktree', 'Shortlog', 'Blame', 'Magic Sync'];

const DOC_TEXT = {
    Init: 'git init creates a new empty Git repository in your current folder.\n\nCommand:\n  git init\n\nAfter running this, Git creates a hidden .git folder\nthat stores all version history and settings.',
    Stage: 'The staging area (index) is where Git stores changes\nyou want to include in your next commit.\n\nCommands:\n  git add file.txt    (stage a single file)\n  git add .           (stage all files)\n  git add -u          (stage modified/deleted)',
    Status: 'git status shows the current state of your working\ndirectory and staging area.\n\nCommands:\n  git status\n  git status -s       (short format)\n  git status -b       (show branch info)',
    Commit: 'A commit is a snapshot of your project at a point in time.\n\nCommands:\n  git commit -m "message"\n  git commit -a -m "message"   (stage + commit)\n  git commit --amend            (amend last commit)',
    Push: 'Push sends your local commits to a remote repository.\n\nCommands:\n  git push\n  git push origin branchname\n  git push -u origin branchname',
    Log: 'Log shows the full history of commits.\n\nCommands:\n  git log --oneline\n  git log --graph --oneline --decorate --all\n  git log -p                  (show diffs)',
    Revert: 'Revert undoes a specific commit by creating a new one.\n\nCommands:\n  git revert <commit-id>\n  git revert --no-commit <old>..<new>',
    Clone: 'Clone creates a local copy of a remote repository.\n\nCommands:\n  git clone <url>\n  git clone <url> myproject\n  git clone --depth 1 <url>',
    Branch: 'A branch is a separate line of development.\n\nCommands:\n  git branch                      (list)\n  git branch feature-login        (create)\n  git switch feature-login        (switch)\n  git branch -d feature-login     (delete)\n  git branch -m new-name          (rename)',
    Pull: 'git pull brings remote changes into your current branch.\n\nCommands:\n  git pull\n  git pull origin main\n  git pull --rebase',
    Reflog: 'Reflog records every move of HEAD, so you can recover\nfrom almost any mistake.\n\nCommands:\n  git reflog\n  git reflog show branchname',
    GitIgnore: 'The .gitignore file lists paths Git should not track.\n\nCommon entries:\n  *.DS_Store\n  build/\n  dist/\n  *.log\n  .env',
    Remote: 'git remote manages connections to other repositories.\n\nCommands:\n  git remote -v              (list)\n  git remote add name url    (add)\n  git remote remove name     (remove)',
    Diff: 'Shows differences between file versions.\n\nCommands:\n  git diff                    (unstaged)\n  git diff --cached          (staged)\n  git diff HEAD              (all changes)\n  git diff --stat            (summary)',
    Reset: 'Moves HEAD and controls commit history.\n\nModes:\n  --soft   undo commit, keep staged\n  --mixed  undo commit, unstage (default)\n  --hard   delete everything (destructive)',
    Fetch: 'Downloads changes from remote without merging.\n\nCommands:\n  git fetch origin\n  git fetch --all',
    Stash: 'Temporarily shelves changes for a clean directory.\n\nCommands:\n  git stash\n  git stash list\n  git stash pop\n  git stash apply',
    Merge: 'Joins two development histories together.\n\nCommands:\n  git merge feature-x\n  git merge --no-ff feature-x',
    Tag: 'Marks specific points as important (releases).\n\nCommands:\n  git tag v1.0\n  git push origin v1.0',
    'Cherry-pick': 'Applies changes from existing commits to current branch.\n\nCommands:\n  git cherry-pick <commit-hash>',
    Rebase: 'Reapplies commits on top of another base tip.\n\nCommands:\n  git rebase main\n  git rebase -i HEAD~3\n  git rebase --continue\n  git rebase --abort',
    Clean: 'Removes untracked files from working tree.\n\nCommands:\n  git clean -n    (preview)\n  git clean -f    (remove files)\n  git clean -fd   (remove files + dirs)',
    Show: 'Shows details about a Git object.\n\nCommands:\n  git show HEAD\n  git show <hash>\n  git show --stat',
    'Ls-files': 'Shows files in the index and working tree.\n\nCommands:\n  git ls-files\n  git ls-files --cached\n  git ls-files --others',
    Worktree: 'Manage multiple working trees.\n\nCommands:\n  git worktree list\n  git worktree add <path> <branch>\n  git worktree remove <name>',
    Shortlog: 'Summarizes git log grouped by author.\n\nCommands:\n  git shortlog\n  git shortlog -s\n  git shortlog -n',
    Blame: 'Shows what revision/author last modified each line.\n\nCommands:\n  git blame <file>\n  git blame -L 10,20 <file>',
    'Magic Sync': 'One-click sync workflow: stash, fetch, pull --rebase,\nand stash pop — in that order, safely.\n\nEquivalent to:\n  git stash push -u\n  git fetch origin\n  git pull --rebase\n  git stash pop',
};

let currentDoc = null;

function renderDocs(root) {
    root.innerHTML = `
        <header class="page-head">
            <div>
                <p class="micro page-kicker">${icon('docs', 12)} Field manual</p>
                <h1 class="page-title">Docs</h1>
                <p class="page-sub">Plain-language reference for every command GitScope runs.</p>
            </div>
        </header>
        <div class="docs-layout">
            <div class="docs-index">
                <label class="field-label" for="docsSearch">Search</label>
                <input type="text" id="docsSearch" placeholder="filter commands…" autocomplete="off" />
                <div class="doc-pills" id="docPills" role="tablist" aria-label="Documentation topics"></div>
            </div>
            <article class="docs-article" id="docArticle" aria-live="polite">
                <p class="docs-empty">Select a command from the index to read its entry.</p>
            </article>
        </div>`;

    const pills = root.querySelector('#docPills');
    const article = root.querySelector('#docArticle');

    const drawPills = (filter) => {
        const q = (filter || '').toLowerCase();
        pills.innerHTML = DOC_ITEMS.filter((t) => t.toLowerCase().includes(q)).map((t) => `
            <button type="button" class="doc-pill" role="tab" data-doc="${esc(t)}"
                aria-selected="${t === currentDoc}">${esc(t)}
                <span class="doc-arrow">${icon('chevRight', 12)}</span></button>`).join('') ||
            '<p class="docs-empty">No matches.</p>';
        hydrateIcons(pills);
        pills.querySelectorAll('[data-doc]').forEach((b) =>
            b.addEventListener('click', () => pickDoc(b.dataset.doc)));
    };

    const pickDoc = (title) => {
        currentDoc = title;
        pills.querySelectorAll('[data-doc]').forEach((b) =>
            b.setAttribute('aria-selected', String(b.dataset.doc === title)));
        const body = DOC_TEXT[title] || 'No documentation available.';
        const highlighted = esc(body).replace(/^(  git .*|git .*)$/gm, '<span class="docs-cmd">$1</span>');
        article.innerHTML = `<h2>${esc(title)}</h2><div class="docs-rule"></div><pre>${highlighted}</pre>`;
    };

    drawPills('');
    root.querySelector('#docsSearch').addEventListener('input', (e) => drawPills(e.target.value));
    if (currentDoc) pickDoc(currentDoc);
}

/* ─── ABOUT PAGE ──────────────────────────────────── */

function renderAbout(root) {
    root.innerHTML = `
        <div class="about">
            <div class="about-mark"><img src="./gitscope-logo.png" alt="GitScope logo" /></div>
            <h1 class="about-statement">Visual Git.<br /><em>Zero friction.</em></h1>
            <p class="about-lede">
                A modern, lightweight Git desktop client built with Go and Wails v2 —
                everything you see is GitScope's own backend, not a shell wrapper.
            </p>
            <div class="about-ledger">
                <div class="ledger-row"><span class="ledger-key">Version</span><span class="ledger-val">2.0.1</span></div>
                <div class="ledger-row"><span class="ledger-key">Stack</span><span class="ledger-val">Go · Wails v2 · vanilla JS</span></div>
                <div class="ledger-row"><span class="ledger-key">Backend</span><span class="ledger-val">${Object.keys(COMMANDS).length} commands · ${CATEGORY_ORDER.length} groups</span></div>
                <div class="ledger-row"><span class="ledger-key">Git</span><span class="ledger-val">${
                    state.gitAvailable
                        ? '<span class="ok">available on PATH</span>'
                        : '<span class="bad">not found on PATH</span>'}</span></div>
                <div class="ledger-row"><span class="ledger-key">Repository</span><span class="ledger-val">${
                    state.repoPath ? esc(state.repoPath) : 'none selected'}</span></div>
                <div class="ledger-row"><span class="ledger-key">Branches</span><span class="ledger-val">${
                    state.repoPath && state.branches.length ? esc(state.branches.join(', ')) : '—'}</span></div>
                <div class="ledger-row"><span class="ledger-key">Shortcuts</span><span class="ledger-val">Ctrl/⌘ K · 1–5 · / · Esc</span></div>
                <div class="ledger-row"><span class="ledger-key">License</span><span class="ledger-val">MIT</span></div>
            </div>
            <div class="about-links">
                <a href="https://github.com/Aswanidev-vs/GitScope" target="_blank" rel="noopener noreferrer">
                    ${icon('github', 15)} View on GitHub
                </a>
            </div>
        </div>`;
}

/* ─── COMMAND PALETTE ─────────────────────────────── */

let paletteOpen = false;
let paletteIndex = 0;
let paletteItems = [];
let paletteReturnFocus = null;

function paletteSource() {
    const nav = PAGES.map((p) => ({
        group: 'Go to', label: p.label, icon: p.icon, kbd: p.num, run: () => navigate(p.page || p.id),
    }));
    const cmds = Object.keys(COMMANDS).map((k) => ({
        group: 'Commands', label: COMMANDS[k].label, hint: COMMANDS[k].cat,
        icon: COMMANDS[k].icon, run: () => exec(k),
    }));
    const util = [
        { group: 'Actions', label: 'Select repository', icon: 'select', run: pickFolder },
        { group: 'Actions', label: 'Refresh repository status', icon: 'refresh', run: loadRepo },
        { group: 'Actions', label: 'Clear console', icon: 'clear', run: clearConsole },
        { group: 'Actions', label: 'Copy repository path', icon: 'copy',
          run: () => { if (state.repoPath) { copyText(state.repoPath); toast('Path copied.', 'success'); } } },
    ];
    return nav.concat(cmds, util);
}

function openPalette() {
    if (paletteOpen || dialogOpen()) return;
    paletteOpen = true;
    paletteIndex = 0;
    paletteReturnFocus = document.activeElement;
    $('#paletteRoot').hidden = false;
    const input = $('#paletteInput');
    input.value = '';
    drawPalette('');
    input.focus();
}

function closePalette() {
    if (!paletteOpen) return;
    paletteOpen = false;
    $('#paletteRoot').hidden = true;
    if (paletteReturnFocus && document.contains(paletteReturnFocus)) paletteReturnFocus.focus();
    paletteReturnFocus = null;
}

function drawPalette(query) {
    const q = (query || '').toLowerCase().trim();
    paletteItems = paletteSource().filter((it) =>
        !q || it.label.toLowerCase().includes(q) ||
        (it.hint || '').toLowerCase().includes(q) || it.group.toLowerCase().includes(q));
    if (paletteIndex >= paletteItems.length) paletteIndex = 0;
    const list = $('#paletteList');
    if (!paletteItems.length) {
        list.innerHTML = '<li class="palette-empty">Nothing matches “' + esc(query) + '”.</li>';
        return;
    }
    let html = '';
    let lastGroup = '';
    paletteItems.forEach((it, i) => {
        if (it.group !== lastGroup) { html += `<li class="palette-group" role="presentation">${it.group}</li>`; lastGroup = it.group; }
        html += `<li role="option" aria-selected="${i === paletteIndex}"><button type="button"
            class="palette-item${i === paletteIndex ? ' is-active' : ''}" data-i="${i}">
            ${icon(it.icon, 15)}<span>${esc(it.label)}</span>
            ${it.hint ? `<span class="pi-kbd">${esc(it.hint)}</span>` : it.kbd ? `<span class="pi-kbd">${it.kbd}</span>` : ''}
        </button></li>`;
    });
    list.innerHTML = html;
    list.querySelectorAll('[data-i]').forEach((b) => {
        b.addEventListener('click', () => {
            const it = paletteItems[Number(b.dataset.i)];
            closePalette();
            if (it) it.run();
        });
    });
    list.querySelector('.palette-item.is-active')?.scrollIntoView({ block: 'nearest' });
}

function movePalette(delta) {
    if (!paletteItems.length) return;
    paletteIndex = (paletteIndex + delta + paletteItems.length) % paletteItems.length;
    drawPalette($('#paletteInput').value);
}

/* ─── COPY HELPER ─────────────────────────────────── */

function copyText(text) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).catch(() => fallbackCopy(text));
    } else fallbackCopy(text);
    function fallbackCopy(t) {
        const ta = document.createElement('textarea');
        ta.value = t;
        ta.style.position = 'fixed';
        ta.style.opacity = '0';
        document.body.appendChild(ta);
        ta.select();
        try { document.execCommand('copy'); } catch (_) { /* noop */ }
        ta.remove();
    }
}

/* ─── CONSOLE RESIZE ──────────────────────────────── */

function initConsoleResize() {
    const handle = $('#consoleResize');
    const panel = $('#consolePanel');
    if (!handle || !panel) return;
    let startY = 0;
    let startH = 0;
    let dragging = false;

    const onMove = (e) => {
        if (!dragging) return;
        const delta = startY - e.clientY;
        const h = Math.max(64, Math.min(window.innerHeight * 0.55, startH + delta));
        panel.style.height = h + 'px';
        panel.style.setProperty('--console-height', h + 'px');
    };
    const onUp = () => {
        dragging = false;
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('mouseup', onUp);
        document.body.style.cursor = '';
        document.body.style.userSelect = '';
    };
    handle.addEventListener('mousedown', (e) => {
        e.preventDefault();
        dragging = true;
        startY = e.clientY;
        startH = panel.offsetHeight;
        document.body.style.cursor = 'ns-resize';
        document.body.style.userSelect = 'none';
        document.addEventListener('mousemove', onMove);
        document.addEventListener('mouseup', onUp);
    });
    handle.addEventListener('keydown', (e) => {
        const step = e.shiftKey ? 48 : 16;
        if (e.key === 'ArrowUp') {
            e.preventDefault();
            panel.style.height = Math.min(window.innerHeight * 0.55, panel.offsetHeight + step) + 'px';
        }
        if (e.key === 'ArrowDown') {
            e.preventDefault();
            panel.style.height = Math.max(64, panel.offsetHeight - step) + 'px';
        }
    });
}

/* ─── GLOBAL KEYBOARD ─────────────────────────────── */

function initKeyboard() {
    document.addEventListener('keydown', (e) => {
        const mod = e.ctrlKey || e.metaKey;

        if (mod && e.key.toLowerCase() === 'k') {
            e.preventDefault();
            if (dialogOpen()) return;
            paletteOpen ? closePalette() : openPalette();
            return;
        }

        if (paletteOpen) {
            if (e.key === 'Escape') { e.preventDefault(); closePalette(); }
            else if (e.key === 'ArrowDown') { e.preventDefault(); movePalette(1); }
            else if (e.key === 'ArrowUp') { e.preventDefault(); movePalette(-1); }
            else if (e.key === 'Enter') {
                e.preventDefault();
                const it = paletteItems[paletteIndex];
                if (it) { closePalette(); it.run(); }
            } else if (e.key === 'Tab') {
                e.preventDefault();
                const active = $('#paletteList .palette-item.is-active');
                if (document.activeElement === $('#paletteInput') && active) active.focus();
                else $('#paletteInput').focus();
            }
            return;
        }

        if (dialogOpen()) return;

        const typing = /^(INPUT|TEXTAREA|SELECT)$/.test(document.activeElement?.tagName || '');
        if (!typing && !mod && /^[1-5]$/.test(e.key)) {
            const p = PAGES[Number(e.key) - 1];
            if (p) { e.preventDefault(); navigate(p.id); }
        }
        if (!typing && e.key === '/') {
            e.preventDefault();
            const filter = $('#historyFilter') || $('#docsSearch');
            if (filter) filter.focus();
            else openPalette();
        }
        if (e.key === 'Escape' && typing) document.activeElement.blur();
    });
}

/* ─── BOOT ────────────────────────────────────────── */

function init() {
    renderShell();
    hydrateIcons(document);

    $$('.rail-btn[data-page]').forEach((b) => b.addEventListener('click', () => navigate(b.dataset.page)));
    $('#railPalette')?.addEventListener('click', openPalette);
    $('#paletteBtn')?.addEventListener('click', openPalette);
    $('#clearConsoleBtn')?.addEventListener('click', clearConsole);
    $('#magicSyncBtn')?.addEventListener('click', () => exec('magic'));
    $('#copyPathBtn')?.addEventListener('click', () => {
        if (!state.repoPath) return;
        copyText(state.repoPath);
        toast('Path copied to clipboard.', 'success');
    });
    $('#paletteInput')?.addEventListener('input', (e) => {
        paletteIndex = 0;
        drawPalette(e.target.value);
    });
    $('#paletteRoot')?.addEventListener('mousedown', (e) => {
        if (e.target.id === 'paletteRoot') closePalette();
    });

    initConsoleResize();
    initKeyboard();

    navigate('repo');

    loadRepo().then(() => {
        syncChrome();
        if (!state.gitAvailable) log('git was not found on PATH — install Git and press Refresh.', 'error');
        else if (!state.repoPath) log('Select a repository to get started.', 'info');
        else log('Opened ' + state.repoPath, 'success');
        startPolling();
    });
}

document.addEventListener('DOMContentLoaded', init);
