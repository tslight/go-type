#!/bin/bash
# Quick reference for book management scripts

# === DOWNLOAD BOOKS ===
# Downloads from Gutendex API with full boilerplate intact
# Usage: python3 download_books.py
# Output: Books in internal/textgen/books/ + updated manifest.json

# === STRIP BOILERPLATE ===
# Removes PG headers/footers and strips to book title
# Usage:
#   python3 strip_boilerplate.py          # Apply stripping
#   python3 strip_boilerplate.py --dry-run # Preview changes

# === COMPLETE WORKFLOW ===
python3 download_books.py      # Get fresh books with full headers
python3 strip_boilerplate.py   # Remove all boilerplate
make all                         # Rebuild Go app

# === CUSTOMIZE BOILERPLATE REMOVAL ===
# Edit this function in strip_boilerplate.py:
#
# def strip_gutenberg_boilerplate(content, book_title=None):
#     """Remove Project Gutenberg header and footer from text content."""
#     # Customize the title matching logic here
#     if book_title:
#         # Find and remove everything before the book title appears
#         ...

# === INSPECT BOOKS ===
head -20 internal/textgen/books/alice*.txt              # See current state
python3 strip_boilerplate.py --dry-run | head -20      # Preview changes

# === RELEVANT FILES ===
# Scripts:
#   - download_books.py (download only, no stripping)
#   - strip_boilerplate.py (boilerplate removal with title matching)
#
# Config:
#   - internal/textgen/books/manifest.json (tracks all books)
#   - internal/textgen/books/*.txt (downloaded book files)
#
# Documentation:
#   - BOOK_WORKFLOW.md (detailed guide)
#   - SCRIPT_SEPARATION.md (changes and customization)
#   - This file: BOOK_SCRIPTS.sh (quick reference)
