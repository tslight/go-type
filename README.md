 ![CI Result](https://github.com/tslight/go-type/actions/workflows/build.yml/badge.svg?event=push) [![Go Report Card](https://goreportcard.com/badge/github.com/tslight/go-type)](https://goreportcard.com/report/github.com/tslight/go-type) [![Go Reference](https://pkg.go.dev/badge/github.com/tslight/go-type.svg)](https://pkg.go.dev/github.com/tslight/go-type)
# GO TYPE! 🚀

A terminal-based typing speed test CLI application written in Go. Real-time character-by-character feedback with accuracy and WPM metrics.

## Features

✨ **Real-time Feedback**
- Characters turn **green** when typed correctly
- Characters turn **red** when typed incorrectly
- Expected text shown in **gray** as an overlay

⚡ **Performance Metrics**
- Words Per Minute (WPM) calculation
- Accuracy percentage
- Error count tracking
- Character completion tracking

📚 **Classic Literature from Project Gutenberg**
- Multiple classic books for variety
- Currently includes: Frankenstein, Alice's Adventures in Wonderland, Through the Looking-Glass, and more
- Embedded at compile time for offline use
- Choose which book to type from with the `-book` flag
- Expand collection by running the included download script

🛠️ **Developer Friendly**
- Comprehensive test suite (92.9% coverage)
- Makefile with convenient targets
- Performance benchmarks included
- Cross-platform support (Linux, macOS, Windows)

## Installation

### Prerequisites
- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/tslight/go-type.git
cd go-type

# Build the binary
make build
```

## Usage

### Basic Typing Test (22 sentences from Frankenstein)
```bash
./cli
# or
go run ./cmd/cli
```

### List Available Books
```bash
./cli -list
```

Output:
```
Available books:
  ID  11: Alice's Adventures in Wonderland
  ID  14: Through the Looking-Glass
  ID  84: Frankenstein
```

### Choose a Specific Book
```bash
./cli -book 11       # Type from Alice's Adventures
./cli -book 14       # Type from Through the Looking-Glass
./cli -book 84       # Type from Frankenstein (default)
```

### Custom Sentence Count
```bash
./cli -w 50              # 50 sentences (from Frankenstein)
./cli -w 100 -book 11    # 100 sentences from Alice
go run ./cmd/cli -w 5    # 5 sentences
```

### Available Flags
```bash
-w int     Number of sentences to generate (default 22)
-book int  Book ID to use (use -list to see available books)
-list      Show all available books and exit
```

## How to Play

1. **Start the test**: Run the CLI
2. **Read the prompt**: Text appears in gray
3. **Type**: As you type, characters overlay the gray text
   - ✅ Green = correct character
   - ❌ Red = wrong character (shows expected char)
   - ➕ Red plus sign = typing beyond the text
4. **Submit**: Press Enter to finish
5. **Review**: See your WPM, accuracy, and error count

### Keyboard Controls
- **Any character**: Type the test
- **Backspace**: Delete previous character
- **Enter**: Submit test
- **Ctrl+C**: Cancel test

## Architecture

### Core Components

**textgen Library** (`internal/textgen/`)
- Generates random paragraphs from embedded Project Gutenberg books
- Sentence-based extraction and randomization
- Multi-book support with book discovery
- Supports runtime book switching via `SetBook()`
- Embedded directory (`books/`) contains all available texts
- Fallback to Frankenstein if book not found

**CLI Application** (`cmd/cli/`)
- Raw terminal mode input handling
- Real-time character-by-character display
- ANSI color code support for terminal styling
- Metrics calculation and reporting
- Book selection with `-book` flag
- Book listing with `-list` flag

### Book Management

The application uses Go's `embed` package to include classic literature:

**Embedded Books** (`internal/textgen/books/`)
- Files follow naming convention: `<id>-<title-lowercase-with-dashes>.txt`
  - Example: `11-alices-adventures-in-wonderland.txt`
  - Example: `1342-pride-and-prejudice.txt`
- Each book is a complete classic from Project Gutenberg
- Compile-time embedding ensures offline availability
- No external dependencies or network access required

**Available Books**
Use `./cli -list` to see all available books. Currently includes 39+ titles:
- Alice's Adventures in Wonderland (ID 11)
- Through the Looking-Glass (ID 14)
- Pride and Prejudice (ID 1342)
- Frankenstein (ID 84)
- The Great Gatsby (ID 25344)
- Dracula (ID 67098)
- Crime and Punishment (ID 769)
- And many more classical works

**Expanding Your Library**

To add more books from Project Gutenberg:

```bash
# The download script automatically fetches and names books correctly
./download_books.sh

# This will:
# 1. Download from Project Gutenberg (respecting rate limits)
# 2. Remove PG headers/footers
# 3. Save as: <id>-<normalized-title>.txt
# 4. Skip duplicates automatically

# After downloading, rebuild to include new books
make build
```

The `download_books.sh` script:
- Downloads popular public domain books (with ID and title pairs)
- Automatically converts titles to filename format (lowercase, dashes for spaces)
- Cleans Project Gutenberg headers/footers
- Saves with consistent `<id>-<title>.txt` naming convention
- Detects and skips duplicate titles automatically
- Implements rate limiting (1 second between requests)

**Example - Adding more books manually:**
```bash
# If you want to add a specific book, you can edit the BOOKS array in download_books.sh
# Format: "ID|Title"
# Example: "74|Jane Eyre"

# Then run the script to download only new books
./download_books.sh
```

## Development

### Code Quality
```bash
make lint                # Format and lint code
make fmt                 # Format code only
make vet                 # Run go vet
make check               # Lint + test
make all                 # Full workflow (clean → lint → test → build)
```

### Adding Tests
Tests follow Go conventions:
- `*_test.go` files in the same package
- Table-driven test patterns for maintainability
- Comprehensive benchmarks for performance tracking

## Dependencies

- **golang.org/x/term** (v0.36.0) - Terminal control and raw mode
- **golang.org/x/sys** (v0.37.0) - System calls

## Metrics Explained

### Words Per Minute (WPM)
```
WPM = (Characters Typed ÷ 5) ÷ Minutes Elapsed
```
Uses the standard 5 characters per word formula.

### Accuracy
```
Accuracy = (Correctly Typed Characters ÷ Total Expected Characters) × 100
```
Measures the percentage of characters you typed correctly.

### Errors
Total number of character mismatches:
- Wrong characters typed
- Missing characters (not typed)
- Extra characters (typed beyond the text)

## Contributing

Contributions are welcome! Please feel free to:
- Report bugs
- Suggest features
- Submit pull requests
- Improve documentation

## License

MIT License - see LICENSE file for details

## Troubleshooting

### "Error: tty" on Windows
Use Windows Terminal, PowerShell, or Git Bash for proper terminal support.

### Dictionary Not Found
The application falls back to the embedded dictionary automatically. No action needed.

### Build Fails
Ensure you have Go 1.21+ installed:
```bash
go version
```

## Future Enhancements

Potential features for future versions:
- [ ] Web-based interface
- [ ] Statistics persistence
- [ ] Leaderboards
- [ ] Multiple languages
- [ ] Difficulty levels
- [ ] Theme customization
- [ ] Touch typing tutorials

## Author

Created with ❤️ as a typing speed test learning project.

---

**Ready to test your typing speed?** Run `make run` and start typing! 🎯
