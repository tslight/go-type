#!/bin/bash
# Convenience wrapper to run the Python download script
# Call the Python script directly instead
python3 "$(dirname "$0")/download_books.py" "$@"
