#!/usr/bin/env python3
"""
Download top 100 most popular books from Project Gutenberg using Gutendex API.

Gutendex (https://gutendex.com) provides a JSON API with the most popular
Project Gutenberg books sorted by download count.

This script:
1. Queries Gutendex API for the top 100 most popular books
2. Downloads books with clean filenames (title-only, no ID prefix)
3. Creates a manifest.json file tracking what was downloaded successfully
4. Does NOT strip boilerplate - use strip_boilerplate.py for that

Note: Run strip_boilerplate.py after downloading to remove PG headers/footers
and additional boilerplate from the book files.
"""

import json
import urllib.request
import urllib.error
import os
import sys
import time
import re
from pathlib import Path

BOOKS_DIR = "assets/books"
MANIFEST_FILE = f"{BOOKS_DIR}/manifest.json"


def normalize_title(title):
    """Normalize title for filename - keep it readable."""
    # Remove special characters but keep spaces and dashes
    title = re.sub(r'[^\w\s-]', '', title)
    # Replace multiple spaces/dashes with single dash
    title = re.sub(r'[\s-]+', '-', title.strip())
    # Remove trailing/leading dashes
    title = title.strip('-').lower()
    return title


def load_manifest():
    """Load existing manifest if it exists."""
    if os.path.exists(MANIFEST_FILE):
        try:
            with open(MANIFEST_FILE, 'r') as f:
                return json.load(f)
        except:
            pass
    return {"books": {}, "total": 0}


def save_manifest(manifest):
    """Save manifest to track downloaded books."""
    with open(MANIFEST_FILE, 'w') as f:
        json.dump(manifest, f, indent=2)


def download_books(max_books=100, max_pages=5):
    """Download books from Gutendex API."""
    Path(BOOKS_DIR).mkdir(parents=True, exist_ok=True)

    manifest = load_manifest()
    downloaded = 0
    processed = 0
    page = 1
    failed_books = []

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

                processed += 1

                # Skip if already in manifest (already successfully downloaded)
                if str(book_id) in manifest["books"]:
                    print(f"✓ Already downloaded: {book_id} - {title}")
                    downloaded += 1
                    continue

                print(f"📥 Downloading ({downloaded}/{max_books}): {book_id} - {title}")

                # Use clean title-only filename (no ID prefix)
                normalized = normalize_title(title)
                filename = f"{BOOKS_DIR}/{normalized}.txt"

                downloaded_successfully = False

                # Try both download URLs
                for url_template in [
                    f"https://www.gutenberg.org/cache/epub/{book_id}/pg{book_id}.txt",
                    f"https://www.gutenberg.org/files/{book_id}/{book_id}-0.txt"
                ]:
                    try:
                        req = urllib.request.Request(url_template, headers={'User-Agent': 'Mozilla/5.0'})
                        with urllib.request.urlopen(req, timeout=10) as response:
                            content = response.read()

                        # Verify we got actual content (at least 5KB)
                        if len(content) > 5000:
                            # Write file as-is (boilerplate will be stripped separately)
                            with open(filename, 'wb') as f:
                                f.write(content)

                            file_size_kb = len(content) / 1024
                            print(f"  ✓ Downloaded: {file_size_kb:.1f}KB")

                            # Record in manifest
                            manifest["books"][str(book_id)] = {
                                "title": title,
                                "filename": normalized + ".txt",
                                "size_kb": round(file_size_kb, 1)
                            }
                            manifest["total"] = len(manifest["books"])

                            downloaded_successfully = True
                            break
                    except (urllib.error.URLError, urllib.error.HTTPError, urllib.error.ContentTooShortError) as e:
                        continue
                    except Exception as e:
                        continue

                if not downloaded_successfully:
                    print(f"✗ Failed: {book_id} - {title}")
                    failed_books.append(f"{book_id} - {title}")
                    if os.path.exists(filename):
                        try:
                            os.remove(filename)
                        except:
                            pass
                else:
                    downloaded += 1

                time.sleep(0.3)

            page += 1

        except Exception as e:
            print(f"Error fetching page {page}: {e}")
            break

    # Save manifest
    save_manifest(manifest)

    # Count final results
    final_count = len([f for f in os.listdir(BOOKS_DIR) if f.endswith('.txt')]) if os.path.exists(BOOKS_DIR) else 0

    print("")
    print("✅ Download complete!")
    print("📊 Summary:")
    print(f"  • Processed: {processed} books from Gutendex API")
    print(f"  • Successfully downloaded: {downloaded} books")
    print(f"  • Total book files: {final_count}")
    print(f"  • Location: {BOOKS_DIR}")
    print(f"  • Manifest: {MANIFEST_FILE}")

    if failed_books:
        print("")
        print("⚠️  Failed downloads:")
        for book in failed_books[:10]:
            print(f"  - {book}")
        if len(failed_books) > 10:
            print(f"  ... and {len(failed_books) - 10} more")

    print("")
    print("📚 Sample downloaded books:")

    if os.path.exists(BOOKS_DIR):
        books_list = sorted([f for f in os.listdir(BOOKS_DIR) if f.endswith('.txt')])
        for book in books_list[:10]:
            print(f"  {book}")
        if len(books_list) > 10:
            print(f"  ... and {len(books_list) - 10} more")

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
