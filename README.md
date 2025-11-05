# go-type 🚀

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

📚 **Text Generation**
- Generates random paragraphs
- System dictionary support
- 235,000+ word embedded dictionary fallback
- Configurable paragraph length (default: 22 words)

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

### Basic Typing Test (22 words)
```bash
./cli
# or
go run ./cmd/cli
```

### Custom Word Count
```bash
./cli -w 50          # 50-word paragraph
go run ./cmd/cli -w 100  # 100-word paragraph
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
- Generates random paragraphs from a word dictionary
- Uses Fisher-Yates shuffle for randomization
- Supports system dictionary + embedded fallback
- Provides multiple generation modes (paragraph, sentence, multi-sentence)

**CLI Application** (`cmd/cli/`)
- Raw terminal mode input handling
- Real-time character-by-character display
- ANSI color code support for terminal styling
- Metrics calculation and reporting

### Dictionary Management

The application uses a two-tier dictionary system:

1. **Primary**: System dictionary (if available)
   - /usr/share/dict/words (Linux)
   - /usr/share/dict/american-english (macOS)
   - C:\Program Files\GNU Aspell\dict\en_US.dict (Windows)

2. **Fallback**: Embedded dictionary (235,976 words)
   - Automatically used if system dictionary not found
   - Ensures consistency across all platforms
   - Embedded at compile time using Go's embed package

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
