package textgen

import (
	"bufio"
	_ "embed"
	"math/rand"
	"os"
	"strings"
	"time"
)

//go:embed dictionary.txt
var dictionaryContent string

var (
	words []string
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// init attempts to load words from system dictionary first, then falls back to embedded dictionary
func init() {
	// Try to load from system dictionary first
	words = loadDictionary()
	// If system dictionary not available, use embedded dictionary
	if len(words) == 0 {
		words = parseEmbeddedDictionary()
	}
}

// loadDictionary attempts to load words from the system dictionary
func loadDictionary() []string {
	// Common dictionary file locations across different OSes
	dictionaryPaths := []string{
		"/usr/share/dict/words",                                             // Linux (Debian/Ubuntu)
		"/usr/dict/words",                                                   // Older Unix systems
		"/usr/share/dict/american-english",                                  // macOS
		"/opt/homebrew/share/dict/words",                                    // macOS Homebrew (Apple Silicon)
		"C:\\Program Files\\GNU Aspell\\dict\\en_US.dict",                   // Windows - Aspell
		"C:\\Users\\AppData\\Local\\Programs\\Git\\usr\\share\\dict\\words", // Windows - Git Bash
	}

	for _, path := range dictionaryPaths {
		if loaded := loadFromFile(path); len(loaded) > 0 {
			return loaded
		}
	}

	return nil
}

// loadFromFile reads words from a dictionary file
func loadFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var loadedWords []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		// Filter for reasonable length words (avoid very short or very long)
		if len(word) >= 3 && len(word) <= 20 && isAlphaOnly(word) {
			loadedWords = append(loadedWords, strings.ToLower(word))
		}
	}

	if len(loadedWords) > 100 {
		// Shuffle the words to avoid sequential selection from sorted dictionary
		shuffleWords(loadedWords)
	}

	return loadedWords
}

// parseEmbeddedDictionary parses the embedded dictionary file
func parseEmbeddedDictionary() []string {
	var loadedWords []string
	scanner := bufio.NewScanner(strings.NewReader(dictionaryContent))

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		// Filter for reasonable length words (avoid very short or very long)
		if len(word) >= 3 && len(word) <= 20 && isAlphaOnly(word) {
			loadedWords = append(loadedWords, strings.ToLower(word))
		}
	}

	if len(loadedWords) > 100 {
		// Shuffle the words to avoid sequential selection from sorted dictionary
		shuffleWords(loadedWords)
	}

	return loadedWords
}

// shuffleWords uses Fisher-Yates algorithm to shuffle the word slice
func shuffleWords(words []string) {
	for i := len(words) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		words[i], words[j] = words[j], words[i]
	}
}

// isAlphaOnly checks if a string contains only alphabetic characters
func isAlphaOnly(s string) bool {
	for _, ch := range s {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
			return false
		}
	}
	return true
}

// GetParagraph returns a randomly generated paragraph of words
func GetParagraph(wordCount int) string {
	if len(words) == 0 {
		return "No dictionary available"
	}

	if wordCount < 1 {
		wordCount = 10
	}

	selectedWords := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		selectedWords[i] = words[rng.Intn(len(words))]
	}

	// Capitalize first word
	if len(selectedWords) > 0 {
		selectedWords[0] = strings.ToUpper(selectedWords[0][:1]) + selectedWords[0][1:]
	}

	// Join with spaces and add period at the end
	paragraph := strings.Join(selectedWords, " ") + "."
	return paragraph
}

// GetRandomSentence returns a single sentence with a random number of words
func GetRandomSentence() string {
	// Random sentence length between 8-15 words
	wordCount := 8 + rng.Intn(8)
	return GetParagraph(wordCount)
}

// GetMultipleSentences returns a paragraph of multiple sentences
func GetMultipleSentences(sentenceCount int) string {
	if sentenceCount < 1 {
		sentenceCount = 3
	}

	var sentences []string
	for i := 0; i < sentenceCount; i++ {
		sentences = append(sentences, GetRandomSentence())
	}

	return strings.Join(sentences, " ")
}
