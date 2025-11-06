#!/usr/bin/env python3
"""
Strip Project Gutenberg boilerplate and other headers from book files.

This script processes downloaded book files and removes:
1. Project Gutenberg header (everything before *** START OF THE PROJECT GUTENBERG EBOOK ***)
2. Project Gutenberg footer (everything after *** END OF THE PROJECT GUTENBERG EBOOK ***)
3. Additional boilerplate before the book title appears
"""

import os
import sys
import json
from pathlib import Path

BOOKS_DIR = "internal/textgen/books"
MANIFEST_FILE = f"{BOOKS_DIR}/manifest.json"


def strip_gutenberg_boilerplate(content, book_title=None):
    """Remove Project Gutenberg header and footer from text content.

    Also removes boilerplate before the book title first appears.
    """
    # Decode if bytes
    if isinstance(content, bytes):
        try:
            content = content.decode('utf-8', errors='ignore')
        except:
            return content

    # Remove BOM if present
    if content.startswith('\ufeff'):
        content = content[1:]

    lines = content.split('\n')

    # Find and remove everything up to and including the START marker
    start_idx = 0
    for i, line in enumerate(lines):
        if '*** START OF THE PROJECT GUTENBERG EBOOK' in line:
            start_idx = i + 1
            break

    # Find and remove everything from the END marker onwards
    end_idx = len(lines)
    for i, line in enumerate(lines):
        if '*** END OF THE PROJECT GUTENBERG EBOOK' in line:
            end_idx = i
            break

    # Extract content between markers
    clean_content = '\n'.join(lines[start_idx:end_idx])

    # If we have a book title, remove everything before its first occurrence
    if book_title:
        # Try exact title first
        title_idx = clean_content.find(book_title)
        if title_idx != -1:
            # Find the start of the line containing the title
            line_start = clean_content.rfind('\n', 0, title_idx)
            if line_start != -1:
                clean_content = clean_content[line_start + 1:]
            else:
                clean_content = clean_content[title_idx:]
        else:
            # Try a simpler search for shorter title variants
            # Remove leading whitespace and take first 50 chars of title
            short_title = book_title.split(':')[0].strip()[:50]
            if len(short_title) > 10:
                title_idx = clean_content.upper().find(short_title.upper())
                if title_idx != -1:
                    line_start = clean_content.rfind('\n', 0, title_idx)
                    if line_start != -1:
                        clean_content = clean_content[line_start + 1:]
                    else:
                        clean_content = clean_content[title_idx:]

    # Remove leading/trailing whitespace
    clean_content = clean_content.strip()

    return clean_content


def strip_books(dry_run=False):
    """Strip boilerplate from all books in the manifest."""

    # Load manifest
    if not os.path.exists(MANIFEST_FILE):
        print(f"❌ Manifest not found at {MANIFEST_FILE}")
        return

    with open(MANIFEST_FILE, 'r') as f:
        manifest = json.load(f)

    books_map = manifest.get("books", {})
    total = len(books_map)
    processed = 0

    print(f"📚 Stripping boilerplate from {total} books...")
    print()

    for book_id_str, book_info in books_map.items():
        processed += 1
        title = book_info.get("title", "Unknown")
        filename = book_info.get("filename", "")

        if not filename:
            print(f"⚠️  ({processed}/{total}) Skipping {book_id_str} - no filename")
            continue

        filepath = os.path.join(BOOKS_DIR, filename)

        if not os.path.exists(filepath):
            print(f"⚠️  ({processed}/{total}) File not found: {filename}")
            continue

        # Read original file
        with open(filepath, 'rb') as f:
            original_content = f.read()

        original_size = len(original_content)

        # Strip boilerplate
        clean_content = strip_gutenberg_boilerplate(original_content, title)
        clean_size = len(clean_content)

        # Calculate reduction
        reduction_kb = (original_size - clean_size) / 1024
        reduction_pct = (reduction_kb * 1024 / original_size * 100) if original_size > 0 else 0

        if dry_run:
            print(f"  ({processed}/{total}) [{reduction_pct:.1f}% reduction] {filename}")
        else:
            # Write cleaned content
            with open(filepath, 'wb') as f:
                f.write(clean_content.encode('utf-8'))

            print(f"  ✓ ({processed}/{total}) [{reduction_pct:.1f}%] {filename}")

            # Update manifest with new size
            new_size_kb = clean_size / 1024
            book_info["size_kb"] = round(new_size_kb, 1)

    # Save updated manifest if not dry run
    if not dry_run:
        with open(MANIFEST_FILE, 'w') as f:
            json.dump(manifest, f, indent=2)
        print()
        print("✅ Boilerplate stripping complete!")
        print(f"📊 Manifest updated with new file sizes")
    else:
        print()
        print("✅ Dry run complete (no changes made)")


if __name__ == "__main__":
    dry_run = "--dry-run" in sys.argv or "-n" in sys.argv

    if dry_run:
        print("🔍 DRY RUN MODE - No changes will be made")
        print()

    strip_books(dry_run=dry_run)
