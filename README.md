![CI Result](https://github.com/tslight/go-type/actions/workflows/build.yml/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/tslight/go-type)](https://goreportcard.com/report/github.com/tslight/go-type) [![Go Reference](https://pkg.go.dev/badge/github.com/tslight/go-type.svg)](https://pkg.go.dev/github.com/tslight/go-type)
# GO TYPE! 🚀

Terminal typing practice in Go, powered by Bubble Tea. The project ships two fully featured apps that share the same UI/UX and persistence – only the content source changes.

## Apps at a Glance
- `gutentype`: practice with embedded Project Gutenberg classics.
- `doctype`: practice with embedded Go standard library documentation.

Both apps include:
- Real-time color feedback (green ✅ / red ❌ / gray prompt overlay)
- WPM, accuracy, error, and character metrics
- Persistent progress & history per book/doc (saved in your home directory)
- Scrollable viewport with lazy loading for long texts

## Quick Start
```bash
git clone https://github.com/tslight/go-type.git
cd go-type
sudo make install # will install to /usr/local/bin by default
```

### Run the apps
```bash
# Classic literature
gutentype          # choose a Project Gutenberg book from the menu
gutentype -list    # list available Project Gutenberg books
gutentype -v       # print version (also available as -version)

# Go documentation
doctype            # choose a Go document from the menu
doctype -list      # list available Go documents
doctype -v         # print version (also available as -version)
```

### Shared keyboard controls
- `Ctrl+Q` / `Ctrl+S`: save results and exit
- `Ctrl+C` / `Ctrl+D`: exit without saving
- `Ctrl+J` / `Ctrl+K`: scroll one line
- `Ctrl+F` / `Ctrl+B`: page down / up

## Project Layout
- `cmd/gutentype`, `cmd/doctype`: CLI entry points
- `pkg/cli`: shared Bubble Tea model, menus, and state provider abstraction
- `internal/textgen`: book library, progress persistence, stats formatter
- `internal/godocgen`: Go doc library, progress persistence, stats formatter
- `assets/`: embedded books and documentation sources

## Development
```bash
make fmt       # gofmt
make lint      # staticcheck, vet, fmt check
make test      # go test ./...
make           # full build, lint, test suite, multiplatform binaries
```

## Contributing & License
Pull requests, issues, and suggestions are welcome. Licensed under the MIT License – see `LICENSE` for details.

Happy typing! 🎯
