# Book Library Management

## Overview

The go-type application includes 100 classic books from Project Gutenberg. The book management system uses a manifest file to track downloaded books and enables easy updates.

## Current Library

**Status**: ✅ 100 books downloaded and ready
**Location**: `internal/textgen/books/`
**Manifest**: `internal/textgen/books/manifest.json`
**Total Size**: ~100 MB of text content

## Manifest System

The `manifest.json` file is the source of truth for all downloaded books. It tracks:

- **book_id**: Project Gutenberg ID (used for re-downloading if needed)
- **title**: Full book title as listed in Project Gutenberg
- **filename**: Clean filename used in the books directory
- **size_kb**: File size in kilobytes

### Example Manifest Entry

```json
{
  "84": {
    "title": "Frankenstein; Or, The Modern Prometheus",
    "filename": "frankenstein-or-the-modern-prometheus.txt",
    "size_kb": 438.4
  }
}
```

## Updating the Library

### First Time Download

```bash
cd /Users/tobe/go-type
python3 download_books.py
make all
```

This will:
1. Query Gutendex API for the 100 most popular Project Gutenberg books
2. Download each book with clean title-only filenames
3. Strip Project Gutenberg boilerplate (headers/footers)
4. Create/update `manifest.json` with download metadata

### Incremental Updates

The script is designed for incremental updates. Simply run it again:

```bash
python3 download_books.py
```

The script will:
- ✅ Skip any books already in the manifest
- ✅ Only download new books
- ✅ Update the manifest with new entries
- ✅ Preserve existing books

## File Naming Convention

Books are stored with clean, human-readable filenames derived from their titles:

- `a-christmas-carol-in-prose-being-a-ghost-story-of-christmas.txt`
- `frankenstein-or-the-modern-prometheus.txt`
- `pride-and-prejudice.txt`
- `alices-adventures-in-wonderland.txt`

**Note**: No Project Gutenberg ID prefix - just the title!

## Project Gutenberg Boilerplate

The download script automatically strips Project Gutenberg headers and footers from all files:

**Removed**:
- Project Gutenberg licensing header (beginning)
- Project Gutenberg footer with encoding info (end)
- Markers: `*** START OF THE PROJECT GUTENBERG EBOOK ...`
- Markers: `*** END OF THE PROJECT GUTENBERG EBOOK ...`

This keeps files clean and focused on the actual book content.

## Adding to the App

The Go application automatically discovers all `.txt` files in `internal/textgen/books/`:

1. Books are loaded at runtime via Go's `embed` package
2. Menu displays titles sorted alphabetically
3. No hardcoding needed - just add files to the directory!

## If Something Goes Wrong

### Books Directory Got Corrupted

```bash
# Clear everything
rm internal/textgen/books/*.txt

# Re-download fresh copies
python3 download_books.py
make all
```

The manifest will help re-download with correct titles!

### Need to Re-download Specific Books

Edit `internal/textgen/books/manifest.json` to remove the book entry for that Project Gutenberg ID, then run:

```bash
python3 download_books.py
```

## Future Enhancements

Possible improvements to the system:

- [ ] Add language filtering (download only English books)
- [ ] Add download progress resumption (partial file recovery)
- [ ] Add book category/genre tagging in manifest
- [ ] Add version tracking for book updates
- [ ] Implement automatic daily/weekly updates
