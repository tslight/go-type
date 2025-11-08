![CI Result](https://github.com/tslight/go-type/actions/workflows/build.yml/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/tslight/go-type)](https://goreportcard.com/report/github.com/tslight/go-type) [![Go Reference](https://pkg.go.dev/badge/github.com/tslight/go-type.svg)](https://pkg.go.dev/github.com/tslight/go-type) [![Latest Release](https://img.shields.io/github/v/release/tslight/go-type?logo=github&label=Release)](https://github.com/tslight/go-type/releases/latest) [![Release Date](https://img.shields.io/github/release-date/tslight/go-type?label=Released)](https://github.com/tslight/go-type/releases/latest) [![Downloads](https://img.shields.io/github/downloads/tslight/go-type/latest/total?label=Downloads)](https://github.com/tslight/go-type/releases/latest) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
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

## Install

```bash
git clone https://github.com/tslight/go-type.git
cd go-type
sudo make install # will install to /usr/local/bin by default
```

Or download one of the pre-compiled binary from the releases page.

### Usage

```bash
gutentype -h
Usage of gutentype:
  -l    List available books and their titles
  -list
        List available books and their titles (long form)
  -v    Show application version
  -version
        Show application version (long form)

doctype -h
Usage of doctype:
  -l    List available Go documentation modules
  -list
        List available Go documentation modules (long form)
  -v    Show application version
  -version
        Show application version (long form) Classic literature
```

### Keyboard Shortcuts

- `Ctrl+Q` / `Ctrl+S`: save results and exit
- `Esc`: save results and return to the menu
- `Ctrl+C`: exit immediately without saving
- `Ctrl+D`: toggle debug overlay
- `Ctrl+J` / `Ctrl+K`: scroll one line
- `Ctrl+F` / `Ctrl+B`: page down / up

## Contributing

Pull requests, issues, and suggestions are welcome. Licensed under the MIT License – see `LICENSE` for details.

Happy typing! 🎯
