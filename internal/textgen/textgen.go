package textgen

import (
	"embed"
	"fmt"
	"math/rand"
	"sort"
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
	sentences        []string
	fullText         string // Full book text
	currentBook      *Book
	rng              = rand.New(rand.NewSource(time.Now().UnixNano()))
	availableBooks   = []Book{}
	stateManager     *StateManager
	currentCharPos   int    // Track character position in the full text (for pager-based resume)
	lastParagraphEnd int    // Track the exact end of the last paragraph returned
	rawBookContent   string // Store raw book content for lazy loading
)

// init initializes the text source on package load
func init() {
	// Load list of available books
	loadAvailableBooks()
	// Initialize state manager
	stateManager = NewStateManager()
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

	// Sort books alphabetically by name
	sort.Slice(availableBooks, func(i, j int) bool {
		return availableBooks[i].Name < availableBooks[j].Name
	})
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

	// Store raw content for lazy loading, don't process entire book yet
	rawBookContent = string(content)
	// Clean up the raw content
	rawBookContent = strings.Split(rawBookContent, "***START")[len(strings.Split(rawBookContent, "***START"))-1]
	rawBookContent = strings.Split(rawBookContent, "***END")[0]
	rawBookContent = strings.TrimSpace(rawBookContent)

	// Only load the first chunk for immediate display
	chunkSize := 50000 // ~50KB initial load
	if len(rawBookContent) > chunkSize {
		fullText = rawBookContent[:chunkSize]
	} else {
		fullText = rawBookContent
	}

	currentBook = &Book{ID: bookID, Name: bookName}

	// Load saved progress for this book
	if state := stateManager.GetState(bookID); state != nil {
		currentCharPos = state.CharacterPos
		// Make sure position doesn't exceed text length
		if currentCharPos > len(fullText) {
			currentCharPos = 0
		}
	} else {
		currentCharPos = 0
	}

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
	// DON'T collapse whitespace - preserve paragraph structure
	// Just remove Project Gutenberg headers/footers
	text = strings.Split(text, "***START")[len(strings.Split(text, "***START"))-1]
	text = strings.Split(text, "***END")[0]
	text = strings.TrimSpace(text)

	// Return the full text as a single "sentence" to preserve formatting
	// The full text with formatting will be used directly
	return []string{text}
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

// GetCurrentCharPos returns the current character position in the full text
func GetCurrentCharPos() int {
	return currentCharPos
}

// GetLastParagraphEnd returns the exact end position of the last paragraph returned
func GetLastParagraphEnd() int {
	return lastParagraphEnd
}

// SetBook loads a different book by ID
func SetBook(bookID int) error {
	return loadBook(bookID)
}

// GetFullText returns the complete text for the current book
// This loads text lazily - only loading what's been typed through
// Filters to ASCII only to match display
func GetFullText() string {
	if len(fullText) == 0 {
		return "No text source available"
	}
	// Ensure we have the full text loaded
	if len(fullText) < len(rawBookContent) {
		fullText = rawBookContent
	}
	return toASCIIFilter(fullText)
}

// toASCIIFilter filters out non-ASCII characters to avoid UTF-8 encoding issues
// Preserves newlines for paragraph formatting
func toASCIIFilter(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		// Keep newlines and ASCII characters
		if s[i] == '\n' || (s[i] < 128 && s[i] >= 32) || s[i] == '\t' {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// CalculateSentencesCompleted calculates how many characters have been typed
// This is deprecated - use character position directly now
func CalculateSentencesCompleted(paragraphLength int) int {
	return paragraphLength
}

// CalculateSentencesCompletedWithCount calculates progress in characters
func CalculateSentencesCompletedWithCount(actualSentenceCount int) int {
	// Return the new character position after completing the paragraph
	return currentCharPos + (actualSentenceCount * 50)
}

// SaveProgress saves the current typing progress for the current book
// charPos is the character position in the full text where user left off
func SaveProgress(charPos int, lastHash string) error {
	if currentBook == nil {
		return nil
	}
	currentCharPos = charPos
	return stateManager.SaveState(currentBook.ID, currentBook.Name, charPos, lastHash)
}

// GetProgress returns the saved progress for the current book
func GetProgress() *BookState {
	if currentBook == nil {
		return nil
	}
	return stateManager.GetState(currentBook.ID)
}

// ClearProgress clears the saved progress for the current book
func ClearProgress() error {
	if currentBook == nil {
		return nil
	}
	currentCharPos = 0
	return stateManager.ClearState(currentBook.ID)
}
