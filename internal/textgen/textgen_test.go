package textgen

import (
	"strings"
	"testing"
)

// TestExtractSentences tests sentence extraction from text
func TestExtractSentences(t *testing.T) {
	text := "This is the first sentence. This is the second sentence! And this is the third question?"
	sentences := extractSentences(text)

	if len(sentences) < 2 {
		t.Errorf("extractSentences should find multiple sentences, got %d", len(sentences))
	}

	// All extracted sentences should have length > 20
	for i, sentence := range sentences {
		if len(sentence) <= 20 {
			t.Errorf("sentence %d is too short: %q (len=%d)", i, sentence, len(sentence))
		}
	}
}

// TestGetParagraph tests paragraph generation
func TestGetParagraph(t *testing.T) {
	tests := []struct {
		name          string
		sentenceCount int
	}{
		{"1 sentence", 1},
		{"3 sentences", 3},
		{"5 sentences", 5},
		{"zero defaults to 3", 0},
		{"negative defaults to 3", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paragraph := GetParagraph(tt.sentenceCount)

			if paragraph == "No text source available" {
				t.Skip("Text source not available")
			}

			// Should return non-empty paragraph
			if len(paragraph) < 20 {
				t.Errorf("GetParagraph(%d) returned too short paragraph: %q", tt.sentenceCount, paragraph)
			}
		})
	}
}

// TestGetRandomSentence tests single sentence generation
func TestGetRandomSentence(t *testing.T) {
	for i := 0; i < 5; i++ {
		sentence := GetRandomSentence()

		if sentence == "No text source available" {
			t.Skip("Text source not available")
		}

		// Should be non-empty
		if len(sentence) == 0 {
			t.Errorf("GetRandomSentence() returned empty string")
		}

		// Should be reasonably long
		if len(sentence) < 20 {
			t.Errorf("GetRandomSentence() returned short sentence: %q", sentence)
		}
	}
}

// TestGetMultipleSentences tests multiple sentence generation
func TestGetMultipleSentences(t *testing.T) {
	result := GetMultipleSentences(3)

	if result == "No text source available" {
		t.Skip("Text source not available")
	}

	if len(result) == 0 {
		t.Errorf("GetMultipleSentences returned empty result")
	}
}

// TestGetAvailableBooks tests book listing
func TestGetAvailableBooks(t *testing.T) {
	books := GetAvailableBooks()

	if len(books) == 0 {
		t.Error("GetAvailableBooks should return at least one book")
	}

	// Should have Frankenstein
	hasFrankenstein := false
	for _, b := range books {
		if b.ID == 84 {
			hasFrankenstein = true
			break
		}
	}
	if !hasFrankenstein {
		t.Error("GetAvailableBooks should always include Frankenstein (ID 84)")
	}
}

// TestSetBook tests switching between books
func TestSetBook(t *testing.T) {
	// Should be able to set to Frankenstein
	err := SetBook(84)
	if err != nil {
		t.Errorf("SetBook(84) failed: %v", err)
	}

	// Current book should be updated
	current := GetCurrentBook()
	if current == nil {
		t.Error("GetCurrentBook() returned nil after SetBook")
	}
	if current.ID != 84 {
		t.Errorf("Expected book ID 84, got %d", current.ID)
	}
}

// TestBookNameExtraction tests that book names are correctly extracted from filenames
func TestBookNameExtraction(t *testing.T) {
	// Get the list of available books
	books := GetAvailableBooks()

	if len(books) == 0 {
		t.Fatal("Expected at least one book to be available")
	}

	// Verify each book has a name
	for _, book := range books {
		if len(book.Name) == 0 {
			t.Errorf("Book %d has empty name", book.ID)
		}

		// Verify name is properly formatted (title case with spaces)
		if strings.Contains(book.Name, "-") {
			t.Errorf("Book %d name should have spaces, not dashes: %q", book.ID, book.Name)
		}
	}

	// Verify Frankenstein is available
	hasFrankenstein := false
	for _, b := range books {
		if b.ID == 84 && strings.Contains(b.Name, "Frankenstein") {
			hasFrankenstein = true
			break
		}
	}
	if !hasFrankenstein {
		t.Error("Frankenstein (ID 84) should be in available books")
	}
}

// BenchmarkGetParagraph benchmarks paragraph generation
func BenchmarkGetParagraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetParagraph(3)
	}
}

// BenchmarkGetRandomSentence benchmarks sentence generation
func BenchmarkGetRandomSentence(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRandomSentence()
	}
}

// BenchmarkExtractSentences benchmarks sentence extraction
func BenchmarkExtractSentences(b *testing.B) {
	testText := "This is the first sentence. This is the second sentence! And this is a question? More text here. Even more sentences to process."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractSentences(testText)
	}
}
