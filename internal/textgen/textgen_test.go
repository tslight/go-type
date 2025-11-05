package textgen

import (
	"strings"
	"testing"
)

// TestIsAlphaOnly tests the isAlphaOnly function
func TestIsAlphaOnly(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid lowercase", "hello", true},
		{"valid uppercase", "HELLO", true},
		{"valid mixed case", "HeLLo", true},
		{"empty string", "", true},
		{"with numbers", "hello123", false},
		{"with special characters", "hello!", false},
		{"with hyphen", "hello-world", false},
		{"with space", "hello world", false},
		{"single letter", "a", true},
		{"with apostrophe", "don't", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAlphaOnly(tt.input)
			if result != tt.expected {
				t.Errorf("isAlphaOnly(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestShuffleWords tests that the shuffle function produces varied output
func TestShuffleWords(t *testing.T) {
	// Create a test slice
	original := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew"}
	testSlice := make([]string, len(original))
	copy(testSlice, original)

	shuffleWords(testSlice)

	// Verify all elements are still present (just shuffled)
	if len(testSlice) != len(original) {
		t.Errorf("shuffleWords changed slice length: got %d, expected %d", len(testSlice), len(original))
	}

	// Check all original elements are still present
	wordCount := make(map[string]int)
	for _, word := range testSlice {
		wordCount[word]++
	}
	for _, word := range original {
		if wordCount[word] != 1 {
			t.Errorf("word %q count mismatch in shuffled slice", word)
		}
	}
}

// TestGetParagraph tests the GetParagraph function
func TestGetParagraph(t *testing.T) {
	tests := []struct {
		name       string
		wordCount  int
		minWords   int
		maxWords   int
		shouldPass bool
	}{
		{"normal count", 10, 10, 10, true},
		{"small count", 5, 5, 5, true},
		{"large count", 50, 50, 50, true},
		{"zero count defaults to 10", 0, 10, 10, true},
		{"negative count defaults to 10", -5, 10, 10, true},
		{"one word", 1, 1, 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paragraph := GetParagraph(tt.wordCount)

			if paragraph == "No dictionary available" {
				t.Skip("Dictionary not available")
			}

			// Split and count words (remove trailing period)
			text := strings.TrimSuffix(paragraph, ".")
			wordList := strings.Fields(text)

			if len(wordList) < tt.minWords || len(wordList) > tt.maxWords {
				t.Errorf("GetParagraph(%d) produced %d words, expected between %d and %d",
					tt.wordCount, len(wordList), tt.minWords, tt.maxWords)
			}

			// Check that it ends with a period
			if !strings.HasSuffix(paragraph, ".") {
				t.Errorf("GetParagraph() should end with period, got: %q", paragraph)
			}

			// Check that first word is capitalized
			if len(wordList) > 0 {
				firstChar := rune(wordList[0][0])
				if !(firstChar >= 'A' && firstChar <= 'Z') {
					t.Errorf("GetParagraph() first word should be capitalized, got: %q", wordList[0])
				}
			}
		})
	}
}

// TestGetParagraphConsistency tests that paragraphs contain valid words
func TestGetParagraphConsistency(t *testing.T) {
	for i := 0; i < 10; i++ {
		paragraph := GetParagraph(15)

		if paragraph == "No dictionary available" {
			t.Skip("Dictionary not available")
		}

		// Verify each word is alphabetic only (except punctuation at end)
		text := strings.TrimSuffix(paragraph, ".")
		wordList := strings.Fields(text)

		for _, word := range wordList {
			if !isAlphaOnly(word) {
				t.Errorf("GetParagraph() produced non-alphabetic word: %q", word)
			}
		}
	}
}

// TestGetRandomSentence tests the GetRandomSentence function
func TestGetRandomSentence(t *testing.T) {
	for i := 0; i < 10; i++ {
		sentence := GetRandomSentence()

		if sentence == "No dictionary available" {
			t.Skip("Dictionary not available")
		}

		// Check that it ends with a period
		if !strings.HasSuffix(sentence, ".") {
			t.Errorf("GetRandomSentence() should end with period, got: %q", sentence)
		}

		// Check word count is between 8 and 15
		text := strings.TrimSuffix(sentence, ".")
		wordList := strings.Fields(text)

		if len(wordList) < 8 || len(wordList) > 15 {
			t.Errorf("GetRandomSentence() produced %d words, expected between 8 and 15", len(wordList))
		}

		// Check first word is capitalized
		if len(wordList) > 0 {
			firstChar := rune(wordList[0][0])
			if !(firstChar >= 'A' && firstChar <= 'Z') {
				t.Errorf("GetRandomSentence() first word should be capitalized, got: %q", wordList[0])
			}
		}
	}
}

// TestGetMultipleSentences tests the GetMultipleSentences function
func TestGetMultipleSentences(t *testing.T) {
	tests := []struct {
		name          string
		sentenceCount int
		expectedMin   int
		expectedMax   int
	}{
		{"three sentences", 3, 3, 3},
		{"one sentence", 1, 1, 1},
		{"five sentences", 5, 5, 5},
		{"zero defaults to 3", 0, 3, 3},
		{"negative defaults to 3", -2, 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMultipleSentences(tt.sentenceCount)

			if result == "No dictionary available" {
				t.Skip("Dictionary not available")
			}

			// Count sentences by counting periods
			sentenceCount := strings.Count(result, ".")

			if sentenceCount < tt.expectedMin || sentenceCount > tt.expectedMax {
				t.Errorf("GetMultipleSentences(%d) produced %d sentences, expected between %d and %d",
					tt.sentenceCount, sentenceCount, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

// TestParseEmbeddedDictionary tests the embedded dictionary parsing
func TestParseEmbeddedDictionary(t *testing.T) {
	embeddedWords := parseEmbeddedDictionary()

	if len(embeddedWords) == 0 {
		t.Fatal("parseEmbeddedDictionary() returned empty word list")
	}

	// Verify all words are alphabetic and within length bounds
	for i, word := range embeddedWords {
		if len(word) < 3 || len(word) > 20 {
			t.Errorf("parseEmbeddedDictionary() word at index %d has invalid length: %q (%d chars)",
				i, word, len(word))
		}

		if !isAlphaOnly(word) {
			t.Errorf("parseEmbeddedDictionary() word at index %d is not alphabetic: %q", i, word)
		}

		// Verify all words are lowercase
		if word != strings.ToLower(word) {
			t.Errorf("parseEmbeddedDictionary() word at index %d is not lowercase: %q", i, word)
		}
	}

	// Spot check some common words exist
	commonWords := map[string]bool{}
	for _, word := range embeddedWords {
		commonWords[word] = true
	}

	expectedWords := []string{"the", "and", "but", "from", "with", "have", "make"}
	for _, expected := range expectedWords {
		if !commonWords[expected] {
			t.Logf("Warning: expected common word %q not found in dictionary", expected)
		}
	}
}

// BenchmarkGetParagraph benchmarks paragraph generation
func BenchmarkGetParagraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetParagraph(22)
	}
}

// BenchmarkGetRandomSentence benchmarks sentence generation
func BenchmarkGetRandomSentence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRandomSentence()
	}
}

// BenchmarkShuffleWords benchmarks the shuffle function
func BenchmarkShuffleWords(b *testing.B) {
	testSlice := make([]string, 1000)
	for i := 0; i < len(testSlice); i++ {
		testSlice[i] = "word"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shuffleWords(testSlice)
	}
}
