# TODO — Modernize Wails-Port UI

## Step 1
Update `wails-port/frontend/src/style.css`:
- Add `:focus-visible` styles for standard controls (buttons, inputs, textarea, select)
- Improve modal content styling to cover `select` and consistent spacing

## Step 2
Update `wails-port/frontend/src/style.css`:
- Move/standardize docs detail `<pre>` styling by introducing a CSS class (e.g. `.docs-pre`)

## Step 3
Update `wails-port/frontend/src/main.js`:
- Remove inline styles for docs detail `<pre>` and replace with the new class

## Step 4
Quick manual verification
- Sidebar navigation still switches pages
- Dashboard modals render and controls are styled
- Docs detail view renders correctly
- Console output remains functional
