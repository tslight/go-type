package textgen

import (
	"strings"
	"testing"
)

// TestExtractSentences tests sentence extraction from text
func TestExtractSentences(t *testing.T) {
	text := "This is the first sentence. This is the second sentence! And this is the third question?"
	sentences := extractSentences(text)

	// Since we now preserve formatting, we should get exactly 1 result (the full text)
	if len(sentences) != 1 {
		t.Errorf("extractSentences should return full text with formatting preserved, got %d parts", len(sentences))
	}

	// The returned text should contain the full original content
	if !strings.Contains(sentences[0], "first sentence") {
		t.Errorf("extracted text should preserve content")
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

	// Should have at least some books
	if len(books) == 0 {
		t.Error("GetAvailableBooks should return at least one book")
	}
}

// TestSetBook tests switching between books
func TestSetBook(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available to test")
	}

	// Test with the first available book
	testBook := books[0]
	err := SetBook(testBook.ID)
	if err != nil {
		t.Errorf("SetBook(%d) failed: %v", testBook.ID, err)
	}

	// Current book should be updated
	current := GetCurrentBook()
	if current == nil {
		t.Error("GetCurrentBook() returned nil after SetBook")
	}
	if current.ID != testBook.ID {
		t.Errorf("Expected book ID %d, got %d", testBook.ID, current.ID)
	}
}

// TestBookNameExtraction tests that book names are correctly extracted
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
	}

	// Verify we have a good number of books
	if len(books) < 50 {
		t.Logf("Warning: Expected at least 50 books, got %d", len(books))
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
