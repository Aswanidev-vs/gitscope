# GitScope - Wails Desktop Client

A modern Git GUI client built with Go + Wails (previously Fyne).

## Architecture

- **Backend**: Go — all Git operations via the `internal/git` package
- **Frontend**: Vanilla JS/HTML/CSS — built with Vite
- **UI Framework**: Wails v2 (WebView2 on Windows)

## Commands

```bash
# Development (hot reload)
wails dev

# Production build
wails build

# Or build manually
cd frontend && npm run build && cd ..
go build -o gitscope_wails.exe ./gitscope-wails/
```

## Structure

| Path | Purpose |
|------|---------|
| `app.go` | Wails App struct with all bound methods |
| `main.go` | Entry point |
| `frontend/` | Web UI source |
| `frontend/src/main.js` | App logic, dialog helpers, console |
| `frontend/src/style.css` | Complete app styles (dark theme) |

## Ported from Fyne

Originally built with the Fyne toolkit. Ported to Wails for:
- Better UI flexibility (HTML/CSS)
- Native system integration
- Smaller binaries (no OpenGL dependency)
