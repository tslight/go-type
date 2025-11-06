package textgen

import (
	"embed"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

//go:embed books/*.txt
var booksFS embed.FS

// Book represents an available book
type Book struct {
	ID   int
	Name string
}

var (
	sentences      []string
	currentBook    *Book
	rng            = rand.New(rand.NewSource(time.Now().UnixNano()))
	availableBooks = []Book{}
)

// init initializes the text source on package load
func init() {
	// Load list of available books
	loadAvailableBooks()
	// Default to a random available book
	if len(availableBooks) > 0 {
		randomBook := availableBooks[rng.Intn(len(availableBooks))]
		loadBook(randomBook.ID)
	} else {
		loadFrankenstein()
	}
}

// loadAvailableBooks discovers embedded books
func loadAvailableBooks() {
	entries, err := booksFS.ReadDir("books")
	if err != nil {
		// Fall back to Frankenstein only
		availableBooks = []Book{{ID: 84, Name: "Frankenstein"}}
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			// Parse book ID and name from filename (ID-title-lowercase-with-dashes.txt)
			filename := strings.TrimSuffix(entry.Name(), ".txt")
			parts := strings.SplitN(filename, "-", 2)

			if len(parts) >= 1 {
				// Parse ID
				var id int
				_, err := fmt.Sscanf(parts[0], "%d", &id)
				if err != nil || id <= 0 {
					continue
				}

				// Extract and format name
				var name string
				if len(parts) > 1 {
					// Convert dashes back to spaces and title case
					name = strings.ReplaceAll(parts[1], "-", " ")
					name = titleCase(name)
				} else {
					name = fmt.Sprintf("Book %d", id)
				}

				availableBooks = append(availableBooks, Book{ID: id, Name: name})
			}
		}
	}

	// Always add Frankenstein if not already present
	hasFrankenstein := false
	for _, b := range availableBooks {
		if b.ID == 84 {
			hasFrankenstein = true
			break
		}
	}
	if !hasFrankenstein {
		availableBooks = append(availableBooks, Book{ID: 84, Name: "Frankenstein"})
	}
}

// titleCase converts a lowercase dash-separated string to title case
func titleCase(s string) string {
	// Replace dashes with spaces
	s = strings.ReplaceAll(s, "-", " ")
	// Capitalize the first letter of each word
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			// Capitalize first letter, keep rest as lowercase
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// loadBook loads sentences from a specific book
func loadBook(bookID int) error {
	// First, try to find the book in our available books list to get the correct name
	var bookName string
	for _, b := range availableBooks {
		if b.ID == bookID {
			bookName = b.Name
			break
		}
	}
	if bookName == "" {
		bookName = fmt.Sprintf("Book %d", bookID)
	}

	// Try to find the book file in embed.FS
	// Try multiple filename formats for compatibility
	var content []byte
	var err error

	// Try: <id>-<title>.txt format
	entries, _ := booksFS.ReadDir("books")
	for _, entry := range entries {
		if !entry.IsDir() {
			filename := entry.Name()
			if strings.HasPrefix(filename, fmt.Sprintf("%d-", bookID)) && strings.HasSuffix(filename, ".txt") {
				content, err = booksFS.ReadFile("books/" + filename)
				if err == nil && len(content) > 0 {
					break
				}
			}
		}
	}

	// If not found, fall back to Frankenstein
	if err != nil || len(content) == 0 {
		return loadFrankenstein()
	}

	sentences = extractSentences(string(content))
	if len(sentences) == 0 {
		return fmt.Errorf("no sentences found in book %d", bookID)
	}

	currentBook = &Book{ID: bookID, Name: bookName}
	return nil
}

// loadFrankenstein loads the Frankenstein book from the books directory
func loadFrankenstein() error {
	// Load Frankenstein from the books directory (ID 84)
	entries, err := booksFS.ReadDir("books")
	if err != nil {
		return fmt.Errorf("failed to read books directory")
	}

	// Find the Frankenstein file (starts with "84-")
	for _, entry := range entries {
		if !entry.IsDir() {
			filename := entry.Name()
			if strings.HasPrefix(filename, "84-") && strings.HasSuffix(filename, ".txt") {
				content, err := booksFS.ReadFile("books/" + filename)
				if err == nil && len(content) > 0 {
					sentences = extractSentences(string(content))
					if len(sentences) > 0 {
						currentBook = &Book{ID: 84, Name: "Frankenstein"}
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("failed to load Frankenstein from books directory")
}

// extractSentences extracts sentences from text
func extractSentences(text string) []string {
	// Remove extra whitespace and newlines
	text = strings.Join(strings.Fields(text), " ")

	// Remove Project Gutenberg headers/footers (basic cleanup)
	text = strings.Split(text, "***START")[len(strings.Split(text, "***START"))-1]
	text = strings.Split(text, "***END")[0]

	// Split on sentence boundaries
	// Matches: period, exclamation, or question mark followed by space or end of string
	re := regexp.MustCompile(`([.!?])\s+`)
	parts := re.Split(text, -1)

	var results []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Only include sentences longer than 20 characters
		if len(part) > 20 {
			results = append(results, part)
		}
	}

	return results
}

// GetParagraph returns a randomly generated paragraph of sentences
func GetParagraph(sentenceCount int) string {
	if len(sentences) == 0 {
		return "No text source available"
	}

	if sentenceCount < 1 {
		sentenceCount = 3
	}

	var selectedSentences []string
	for i := 0; i < sentenceCount; i++ {
		idx := rng.Intn(len(sentences))
		selectedSentences = append(selectedSentences, sentences[idx])
	}

	return strings.Join(selectedSentences, " ")
}

// GetRandomSentence returns a single random sentence
func GetRandomSentence() string {
	if len(sentences) == 0 {
		return "No text source available"
	}

	idx := rng.Intn(len(sentences))
	return sentences[idx]
}

// GetMultipleSentences returns multiple random sentences
func GetMultipleSentences(count int) string {
	if count < 1 {
		count = 3
	}

	var result []string
	for i := 0; i < count; i++ {
		result = append(result, GetRandomSentence())
	}

	return strings.Join(result, " ")
}

// GetAvailableBooks returns the list of available books
func GetAvailableBooks() []Book {
	return availableBooks
}

// GetCurrentBook returns the currently loaded book
func GetCurrentBook() *Book {
	return currentBook
}

// SetBook loads a different book by ID
func SetBook(bookID int) error {
	return loadBook(bookID)
}
