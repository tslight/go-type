#!/usr/bin/env python3
"""
Download top 100 most popular books from Project Gutenberg using Gutendex API.

Gutendex (https://gutendex.com) provides a JSON API with the most popular
Project Gutenberg books sorted by download count.
"""

import json
import urllib.request
import urllib.error
import os
import sys
import time
from pathlib import Path

BOOKS_DIR = "internal/textgen/books"


def normalize_title(title):
    """Normalize title for filename."""
    normalized = title.lower()
    normalized = normalized.replace("'", "")
    normalized = normalized.split(':')[0]  # Remove subtitle
    normalized = ' '.join(normalized.split())  # Normalize spaces
    normalized = normalized.replace(' ', '-')
    normalized = ''.join(c for c in normalized if c.isalnum() or c == '-')
    normalized = '-'.join(filter(None, normalized.split('-')))
    return normalized


def download_books(max_books=100, max_pages=5):
    """Download books from Gutendex API."""
    Path(BOOKS_DIR).mkdir(parents=True, exist_ok=True)

    downloaded = 0
    processed = 0
    page = 1

    print("📚 Fetching top 100 most popular Project Gutenberg books from Gutendex API...")
    print("")

    while downloaded < max_books and page <= max_pages:
        print(f"📄 Fetching page {page} from Gutendex...")

        try:
            url = f"https://gutendex.com/books?page={page}"
            req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0'})
            with urllib.request.urlopen(req, timeout=10) as response:
                data = json.loads(response.read().decode('utf-8'))

            books = data.get('results', [])
            if not books:
                print("No more books available")
                break

            for book in books:
                if downloaded >= max_books:
                    break

                book_id = book.get('id')
                title = book.get('title', '')

                if not book_id or not title:
                    continue

                normalized = normalize_title(title)
                filename = f"{BOOKS_DIR}/{book_id}-{normalized}.txt"
                processed += 1

                # Skip if already exists
                if os.path.exists(filename):
                    print(f"✓ Already have ({downloaded}/{max_books}): {book_id} - {title}")
                    downloaded += 1
                    continue

                print(f"📥 Downloading ({downloaded}/{max_books}): {book_id} - {title}")

                downloaded_successfully = False

                # Try both download URLs
                for url_template in [
                    f"https://www.gutenberg.org/cache/epub/{book_id}/pg{book_id}.txt",
                    f"https://www.gutenberg.org/files/{book_id}/{book_id}-0.txt"
                ]:
                    try:
                        req = urllib.request.Request(url_template, headers={'User-Agent': 'Mozilla/5.0'})
                        with urllib.request.urlopen(req, timeout=5) as response:
                            content = response.read().decode('utf-8', errors='ignore')

                            # Remove PG headers/footers
                            lines = content.split('\n')
                            start_idx = next(
                                (i for i, line in enumerate(lines) if '***START' in line), 0
                            )
                            end_idx = next(
                                (i for i, line in enumerate(lines) if '***END' in line), len(lines)
                            )
                            content = '\n'.join(lines[start_idx:end_idx])

                            with open(filename, 'w', encoding='utf-8') as f:
                                f.write(content)

                            size = os.path.getsize(filename) / 1024
                            alt_msg = " (alt URL)" if "files" in url_template else ""
                            print(f"✅ Downloaded{alt_msg}: {book_id} - {title} ({size:.0f}K)")
                            downloaded += 1
                            downloaded_successfully = True
                            break
                    except (urllib.error.URLError, urllib.error.HTTPError, Exception):
                        continue

                if not downloaded_successfully:
                    print(f"✗ Failed: {book_id} - {title}")
                    if os.path.exists(filename):
                        os.remove(filename)

                time.sleep(0.5)

            page += 1

        except Exception as e:
            print(f"Error fetching page {page}: {e}")
            break

    # Count final results
    final_count = len(os.listdir(BOOKS_DIR)) if os.path.exists(BOOKS_DIR) else 0

    print("")
    print("✅ Download complete!")
    print("📊 Summary:")
    print(f"  • Processed: {processed} books from Gutendex API")
    print(f"  • Successfully downloaded: {downloaded} books")
    print(f"  • Total in library: {final_count} books")
    print(f"  • Location: {BOOKS_DIR}")
    print("")
    print("📚 Sample books:")

    if os.path.exists(BOOKS_DIR):
        books_list = sorted(os.listdir(BOOKS_DIR))
        for book in books_list[:15]:
            print(f"  {book}")
        if len(books_list) > 15:
            print(f"  ... and {len(books_list) - 15} more")

    return downloaded, final_count


if __name__ == "__main__":
    try:
        download_books()
    except KeyboardInterrupt:
        print("\n\n⚠️  Download interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Error: {e}", file=sys.stderr)
        sys.exit(1)
