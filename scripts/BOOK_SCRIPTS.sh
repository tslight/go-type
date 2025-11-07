#!/bin/bash
# Quick reference for book management scripts

# === DOWNLOAD BOOKS ===
# Downloads from Gutendex API with full boilerplate intact
# Usage: python3 scripts/download_books.py
# Output: Books in internal/textgen/books/ + updated manifest.json

# === STRIP BOILERPLATE ===
# Removes PG headers/footers, [Illustration] lines, and consolidates empty lines
# Usage:
#   python3 scripts/strip_boilerplate.py          # Apply stripping
#   python3 scripts/strip_boilerplate.py --dry-run # Preview changes

# === COMPLETE WORKFLOW ===
python3 scripts/download_books.py      # Get fresh books with full headers
python3 scripts/strip_boilerplate.py   # Remove all boilerplate
make all                                # Rebuild Go app

# === CUSTOMIZE BOILERPLATE REMOVAL ===
# Edit this function in scripts/strip_boilerplate.py:
#
# def strip_gutenberg_boilerplate(content):
#     """Remove Project Gutenberg header and footer from text content."""
#     # Customize the stripping logic here
#     ...

# === INSPECT BOOKS ===
head -20 internal/textgen/books/alice*.txt              # See current state
python3 scripts/strip_boilerplate.py --dry-run | head  # Preview changes

# === RELEVANT FILES ===
# Scripts (scripts/ directory):
#   - download_books.py (download only, no stripping)
#   - strip_boilerplate.py (boilerplate removal)
#   - download_books.sh (utility script)
#   - BOOK_SCRIPTS.sh (this quick reference)
#
# Config:
#   - internal/textgen/books/manifest.json (tracks all books)
#   - internal/textgen/books/*.txt (downloaded book files)
#
# Documentation (docs/ directory):
#   - BOOK_WORKFLOW.md (detailed guide)
#   - SCRIPT_SEPARATION.md (changes and customization)
#   - This file: scripts/BOOK_SCRIPTS.sh (quick reference)
