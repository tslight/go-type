# Recent Changes - Multi-Book Support (v2.0)

## Summary
Refactored the text generation system from single-source dictionary to multi-book system using Project Gutenberg classics. Users can now choose which book to type from.

## Key Changes

### textgen.go (Complete Rewrite)
- ✅ Removed word-based generation (dictionary system)
- ✅ Removed Go documentation source
- ✅ Added `Book` struct with ID and Name fields
- ✅ Added `embed.FS` directory embedding for `books/` folder
- ✅ Implemented `loadAvailableBooks()` - scans embedded directory
- ✅ Implemented `loadBook(id int)` - loads specific book from embed.FS
- ✅ Implemented `getBookName(id int)` - maps book IDs to display names
- ✅ Implemented `GetAvailableBooks()` - returns list of available books
- ✅ Implemented `GetCurrentBook()` - returns currently loaded book
- ✅ Implemented `SetBook(id int)` - switches books at runtime
- ✅ Kept sentence extraction and paragraph generation logic
- ✅ Fallback to embedded Frankenstein if directory unavailable

### cmd/gutentype/main.go
- ✅ Removed `-source` flag (was: 'book'/'godocs')
- ✅ Added `-book int` flag for book ID selection
- ✅ Added `-list` flag to show available books
- ✅ Updated header to show current book name

### internal/textgen/textgen_test.go
- ✅ Removed `TestSourceName` (old API)
- ✅ Added `TestGetAvailableBooks()` - verifies book listing
- ✅ Added `TestSetBook()` - verifies book switching
- ✅ Added `TestGetBookName()` - verifies name mapping
- ✅ All existing paragraph/sentence tests still pass

### README.md
- ✅ Updated features section for Project Gutenberg books
- ✅ Rewrote usage examples with book selection
- ✅ Added `-list` flag documentation
- ✅ Added book management section
- ✅ Updated architecture description

### Project Structure
- ✅ `internal/textgen/books/book_11.txt` - Alice's Adventures (170 KB)
- ✅ `internal/textgen/books/book_14.txt` - Through the Looking-Glass (1.9 MB)
- ✅ `internal/textgen/frankenstein_clean.txt` - Frankenstein (fallback, 416 KB)
- ✅ `download_books.sh` - Script to download/clean additional books

## Usage Examples

```bash
# List available books
./cli -list

# Type from Alice (ID 11)
./cli -book 11

# Type from Frankenstein with custom sentence count
./cli -w 50 -book 84

# Default (Frankenstein, 22 sentences)
./cli
```

## Testing Results
- ✅ All unit tests pass
- ✅ All benchmarks functional
- ✅ Code lints successfully
- ✅ Builds without errors
- ✅ Binary size: 4.9M

## Next Steps
- Run `./download_books.sh` to download more books from Project Gutenberg
- Rebuild with `make build` to include new books
- New books automatically appear in `-list` output

## Technical Details
- Uses Go's `embed.FS` for compile-time file embedding
- Sentence extraction via regex: `([.!?])\s+`
- Filters sentences < 20 characters
- Book discovery scans `books/` directory for `book_*.txt` files
- Graceful fallback to Frankenstein if book not found
