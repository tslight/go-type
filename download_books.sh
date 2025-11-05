#!/bin/bash
# Download top 100 most popular books from Project Gutenberg
# Filenames will be formatted as: <id>-<title-lowercase-with-dashes>.txt

BOOKS_DIR="internal/textgen/books"
mkdir -p "$BOOKS_DIR"

# Most popular Project Gutenberg books
# Format: "ID|Title"
declare -a BOOKS=(
    "11|Alice's Adventures in Wonderland"
    "14|Through the Looking-Glass"
    "76|Adventures of Huckleberry Finn"
    "84|Frankenstein"
    "98|A Tale of Two Cities"
    "103|The Murders in the Rue Morgue"
    "120|Treasure Island"
    "158|Emma"
    "161|The Gettysburg Address and Other Speeches"
    "219|Heart of Darkness"
    "244|The Picture of Dorian Gray"
    "768|A Christmas Carol"
    "769|Crime and Punishment"
    "1023|The Complete Works of William Shakespeare"
    "1080|The Picture of Dorian Gray"
    "1228|Oliver Twist"
    "1232|The Prince"
    "1248|The Picture of Dorian Gray"
    "1342|Pride and Prejudice"
    "1513|Vanity Fair"
    "1514|The Adventures of Tom Sawyer"
    "1517|Uncle Tom's Cabin"
    "1524|Moby Dick"
    "1661|A Study in Scarlet"
    "1952|The Yellow Wallpaper"
    "2814|Wuthering Heights"
    "4280|Beowulf"
    "4363|The Odyssey"
    "5200|Metamorphosis"
    "6130|The Moonstone"
    "7084|Frankenstein: or, The Modern Prometheus"
    "25344|The Great Gatsby"
    "30254|Moby Dick"
    "43362|Twenty Thousand Leagues Under the Sea"
    "44488|The Hound of the Baskervilles"
    "46796|Sense and Sensibility"
    "67098|Dracula"
    "11457|Peter Pan"
    "74|Jane Eyre"
    "145|The Man in the Iron Mask"
    "203|Beowulf"
    "514|Little Women"
    "766|A Christmas Carol"
    "3207|Grimms' Fairy Tales"
    "5740|Aesop's Fables"
)

# Function to convert title to filename format: lowercase, spaces to dashes
normalize_title() {
    local title="$1"
    # Convert to lowercase
    title=$(echo "$title" | tr '[:upper:]' '[:lower:]')
    # Replace apostrophes with nothing
    title=$(echo "$title" | sed "s/'//g")
    # Replace colons and everything after with nothing (e.g., "Title: Subtitle" becomes "Title")
    title=$(echo "$title" | sed 's/:.*//')
    # Replace multiple spaces with single space first
    title=$(echo "$title" | sed 's/[[:space:]]\+/ /g')
    # Replace spaces with dashes
    title=$(echo "$title" | sed 's/ /-/g')
    # Remove any remaining special characters except dashes
    title=$(echo "$title" | sed 's/[^a-z0-9-]//g')
    # Remove multiple consecutive dashes
    title=$(echo "$title" | sed 's/-\{2,\}/-/g')
    # Remove trailing dashes
    title=$(echo "$title" | sed 's/-*$//')
    echo "$title"
}

echo "Downloading top Project Gutenberg books..."
COUNT=0
DOWNLOADED=0
SEEN_TITLES=""

for book_entry in "${BOOKS[@]}"; do
    if [ $DOWNLOADED -ge 100 ]; then
        break
    fi

    # Parse ID and title
    ID=$(echo "$book_entry" | cut -d'|' -f1)
    TITLE=$(echo "$book_entry" | cut -d'|' -f2)

    # Normalize title for filename
    NORMALIZED=$(normalize_title "$TITLE")
    FILENAME="$BOOKS_DIR/${ID}-${NORMALIZED}.txt"

    # Skip duplicates based on title
    if [[ "$SEEN_TITLES" == *"$NORMALIZED"* ]]; then
        echo "⏭️  Skipping duplicate: $ID - $TITLE"
        COUNT=$((COUNT + 1))
        continue
    fi
    SEEN_TITLES="$SEEN_TITLES $NORMALIZED"

    # Skip if already exists
    if [ -f "$FILENAME" ]; then
        echo "⏭️  Skipping (already exists): $ID - $TITLE"
        COUNT=$((COUNT + 1))
        DOWNLOADED=$((DOWNLOADED + 1))
        continue
    fi

    echo "📥 Downloading ($DOWNLOADED/100): $ID - $TITLE"

    # Download the text file
    curl -s "https://www.gutenberg.org/cache/epub/${ID}/pg${ID}.txt" -o "$FILENAME"

    # Check if download was successful
    if [ -s "$FILENAME" ]; then
        # Clean the file (remove Project Gutenberg headers/footers)
        sed -i '' '/\*\*\*START/,/\*\*\*END/d' "$FILENAME"

        SIZE=$(du -h "$FILENAME" | cut -f1)
        echo "✅ Downloaded: $ID - $TITLE ($SIZE)"
        COUNT=$((COUNT + 1))
        DOWNLOADED=$((DOWNLOADED + 1))
    else
        echo "❌ Failed to download: $ID - $TITLE"
        rm -f "$FILENAME"
    fi

    # Rate limit to be respectful to Gutenberg servers
    sleep 1
done

echo ""
echo "✅ Downloaded $DOWNLOADED unique books to $BOOKS_DIR"
echo "📝 Total processed: $COUNT entries"
echo ""
echo "Available books:"
ls -1 "$BOOKS_DIR" | head -20
echo "..."
