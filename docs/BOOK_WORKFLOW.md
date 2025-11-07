# Book Download & Processing Workflow

The book system is now split into two separate scripts for flexibility in boilerplate removal:

## 1. Download Books

**Script**: `scripts/download_books.py`

Downloads books from the Gutendex API and saves them with all original content intact (including PG headers/footers).

```bash
python3 scripts/download_books.py
```

**What it does**:
- Fetches the top 100 most popular Project Gutenberg books from Gutendex API
- Creates clean title-only filenames (no ID prefixes)
- Saves raw content as-is (no stripping yet)
- Creates/updates `manifest.json` with book metadata
- Skips books already in the manifest (incremental updates)

**Output**:
- Book files in `internal/textgen/books/`
- Manifest at `internal/textgen/books/manifest.json`

## 2. Strip Boilerplate

**Script**: `scripts/strip_boilerplate.py`

Removes Project Gutenberg headers, footers, and other boilerplate from downloaded files.

```bash
# Process all books in manifest
python3 scripts/strip_boilerplate.py

# Dry run - show what would be changed
python3 scripts/strip_boilerplate.py --dry-run
```

**What it does**:
- Removes UTF-8 BOM if present
- Strips PG header (everything before `*** START OF THE PROJECT GUTENBERG EBOOK ***`)
- Strips PG footer (everything after `*** END OF THE PROJECT GUTENBERG EBOOK ***`)
- Removes lines containing only `[Illustration]`
- Reduces multiple consecutive empty lines to single empty lines
- Updates `manifest.json` with new file sizes

**Customization**:

To adjust how boilerplate is stripped, edit the `strip_gutenberg_boilerplate()` function in `scripts/strip_boilerplate.py`:

```python
def strip_gutenberg_boilerplate(content, book_title=None):
    # ... existing code ...

    # Currently strips to book title. Customize this section:
    if book_title:
        # Try exact title first
        title_idx = clean_content.find(book_title)
        # ... handle title matching ...
```

Examples of adjustments:
- Strip to a specific phrase instead of title
- Remove entire sections (like illustrations, tables of contents)
- Handle language-specific boilerplate
- Strip multiple markers

## Complete Workflow

1. **First time setup**:
   ```bash
   python3 scripts/download_books.py     # Downloads 100 books with headers intact
   python3 scripts/strip_boilerplate.py  # Removes all boilerplate
   make all                               # Build Go app
   ```

2. **Add more books later**:
   ```bash
   python3 scripts/download_books.py     # Fetches any new books not in manifest
   python3 scripts/strip_boilerplate.py  # Strips boilerplate from new files
   make all                               # Rebuild
   ```

3. **Adjust stripping rules**:
   ```bash
   python3 scripts/strip_boilerplate.py --dry-run  # Preview changes
   # Edit scripts/strip_boilerplate.py as needed
   python3 scripts/strip_boilerplate.py            # Apply updated rules
   ```

## File Size Behavior

Original files from Gutendex include full PG boilerplate (usually 10-30KB overhead).

After `strip_boilerplate.py`:
- Example: Alice's Adventures → 174KB → 144KB (30KB removed)
- Most books: 15-25% size reduction after stripping

The manifest is automatically updated with new file sizes after stripping.

## Troubleshooting

**"Manifest not found"**:
- Run `download_books.py` first to create it

**Books not stripped**:
- Ensure manifest has the books listed
- Check file paths in manifest match actual files in `internal/textgen/books/`

**Title not found for stripping**:
- The script handles this gracefully - it still removes PG markers
- Some books may need manual title matching adjustments in the stripping function

## Technical Details

- **Manifest**: `internal/textgen/books/manifest.json` - JSON with book metadata (ID, title, filename, size)
- **Books directory**: `internal/textgen/books/` - Contains all downloaded .txt files
- **Download source**: Gutendex API (https://gutendex.com) - Top 100 most popular PG books
- **Go integration**: Embedded via `//go:embed books/*.txt books/manifest.json`
