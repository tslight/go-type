# Script Separation Complete ✅

## Changes Made

Your Python scripts for book management have been separated into two focused tools:

### 1. **`download_books.py`** - Download Only
- Fetches books from Gutendex API
- Saves raw content with **all boilerplate intact**
- Creates/updates manifest.json
- Supports incremental updates (skips already-downloaded books)
- **Removed**: `strip_gutenberg_boilerplate()` function

### 2. **`strip_boilerplate.py`** - Boilerplate Removal (NEW)
- Processes all books in the manifest
- Removes UTF-8 BOM
- Removes PG headers (before `*** START OF THE PROJECT GUTENBERG EBOOK ***`)
- Removes PG footers (after `*** END OF THE PROJECT GUTENBERG EBOOK ***`)
- **NEW**: Strips content up to first occurrence of the book title
- Skips lines already having book title boilerplate removed
- Supports dry-run mode (`--dry-run`)
- Updates manifest with new file sizes

## Why This Separation?

You can now:
1. **Adjust boilerplate removal independently** - Edit `strip_boilerplate.py` without touching download logic
2. **Re-strip books with new rules** - Change the stripping function and re-run without re-downloading
3. **Handle edge cases per book** - Customize stripping for specific books if needed
4. **Experiment safely** - Use `--dry-run` to preview changes before applying

## How to Customize Boilerplate Removal

Edit the `strip_gutenberg_boilerplate()` function in `strip_boilerplate.py` (lines 28-79):

**Current behavior**:
```python
# Strips PG markers, then removes everything before the book title
if book_title:
    title_idx = clean_content.find(book_title)
    # ... removes content before this point ...
```

**Example customizations**:

Strip to a keyword instead:
```python
# Remove everything before "Chapter 1"
chapter_idx = clean_content.upper().find("CHAPTER 1")
if chapter_idx != -1:
    clean_content = clean_content[chapter_idx:]
```

Remove specific sections:
```python
# Remove "Contents" sections
if "TABLE OF CONTENTS" in clean_content:
    toc_end = clean_content.find("\n\n\n", toc_start) + 3
    clean_content = clean_content[:toc_start] + clean_content[toc_end:]
```

Skip stripping for certain books:
```python
if book_title and "Poetry" not in book_title:
    # Apply custom logic only for non-poetry
    ...
```

## Usage

**Complete workflow**:
```bash
# 1. Download books (with full PG headers/footers)
python3 download_books.py

# 2. Strip boilerplate using current rules
python3 strip_boilerplate.py

# 3. Rebuild Go app
make all
```

**Experiment with new stripping rules**:
```bash
# Preview changes without modifying files
python3 strip_boilerplate.py --dry-run

# Once satisfied with rules, apply them
python3 strip_boilerplate.py
```

## Current Results

Dry-run preview shows:
- Most books: 1-3% boilerplate reduction
- Some books (like Romeo & Juliet): 65% reduction (lots of extra metadata/intro)
- Files still have PG markers if title wasn't found (fallback to marker-only removal)

## Files Modified/Created

- ✅ `download_books.py` - Removed boilerplate function, simplified to download only
- ✅ `strip_boilerplate.py` - New script for flexible boilerplate removal
- ✅ `BOOK_WORKFLOW.md` - Comprehensive guide on using both scripts
- ✅ `SCRIPT_SEPARATION.md` - This file

## Next Steps

1. **Review the stripping rules** - Check if current title-based stripping works for your books
2. **Test with a few books** - Run the dry-run, inspect a few book files manually
3. **Customize if needed** - Adjust `strip_gutenberg_boilerplate()` function as needed
4. **Re-strip all books** - Run `python3 strip_boilerplate.py` when satisfied

Example inspection:
```bash
# See first 50 lines before stripping
head -50 internal/textgen/books/alice*.txt

# After stripping, verify boilerplate is gone
python3 strip_boilerplate.py
head -50 internal/textgen/books/alice*.txt
```

## Build Status

✅ All tests passing
✅ Go app builds successfully
✅ Both scripts have valid Python syntax
✅ Ready to use!
