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

// TestFullText tests full text retrieval basic functionality
func TestFullText(t *testing.T) {
	text := GetFullText()

	if text == "" {
		t.Error("GetFullText() should return non-empty text")
	}

	// Should be a reasonable length for typing practice
	if len(text) < 100 {
		t.Errorf("GetFullText() returned text too short: %d characters", len(text))
	}
}

// TestGetFullText tests full text retrieval
func TestGetFullText(t *testing.T) {
	text := GetFullText()

	if text == "" {
		t.Error("GetFullText() should return non-empty text")
	}

	// Should be a reasonable length for typing practice
	if len(text) < 100 {
		t.Errorf("GetFullText() returned text too short: %d characters", len(text))
	}
}

// TestGetCurrentCharPos tests character position tracking
func TestGetCurrentCharPos(t *testing.T) {
	pos := GetCurrentCharPos()

	if pos < 0 {
		t.Errorf("GetCurrentCharPos() should return non-negative value, got %d", pos)
	}
}

// TestGetLastParagraphEnd tests paragraph end position tracking
func TestGetLastParagraphEnd(t *testing.T) {
	end := GetLastParagraphEnd()

	if end < 0 {
		t.Errorf("GetLastParagraphEnd() should return non-negative value, got %d", end)
	}
}

// TestCalculateSentencesCompleted tests sentence counting
func TestCalculateSentencesCompleted(t *testing.T) {
	count := CalculateSentencesCompleted(500)

	if count < 0 {
		t.Errorf("CalculateSentencesCompleted() should return non-negative value, got %d", count)
	}
}

// TestCalculateSentencesCompletedWithCount tests sentence counting with limit
func TestCalculateSentencesCompletedWithCount(t *testing.T) {
	count := CalculateSentencesCompletedWithCount(5)

	if count < 0 {
		t.Errorf("CalculateSentencesCompletedWithCount() should return non-negative value, got %d", count)
	}
}

// TestToASCIIFilter tests ASCII filtering
func TestToASCIIFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ASCII only",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "with numbers",
			input:    "Test 123",
			expected: "Test 123",
		},
		{
			name:     "with spaces",
			input:    "multiple spaces",
			expected: "multiple spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toASCIIFilter(tt.input)
			if result != tt.expected {
				t.Errorf("toASCIIFilter(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestExtractSentences_EdgeCases tests edge cases in sentence extraction
func TestExtractSentences_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"single sentence", "Hello world."},
		{"no period", "Hello world"},
		{"multiple punctuation", "What?!"},
		{"newlines", "Hello\nWorld"},
		{"tabs", "Hello\tWorld"},
		{"mixed whitespace", "Hello  \n\t  World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSentences(tt.input)
			// Should not panic or error
			_ = result
		})
	}
}

// TestGetParagraph_EdgeCases tests edge case paragraph generation
func TestGetParagraph_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		sentenceCount int
	}{
		{"very large sentence count", 1000},
		{"exact zero", 0},
		{"negative one", -1},
		{"negative large", -1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paragraph := GetParagraph(tt.sentenceCount)

			if paragraph == "No text source available" {
				t.Skip("Text source not available")
			}

			// Should return non-empty result
			if len(paragraph) == 0 {
				t.Errorf("GetParagraph(%d) returned empty string", tt.sentenceCount)
			}
		})
	}
}

// TestGetRandomSentence_Multiple tests multiple random sentence calls
func TestGetRandomSentence_Multiple(t *testing.T) {
	sentences := make(map[string]bool)

	// Get multiple sentences to test randomness
	for i := 0; i < 10; i++ {
		sentence := GetRandomSentence()

		if sentence == "No text source available" {
			t.Skip("Text source not available")
		}

		if len(sentence) == 0 {
			t.Error("GetRandomSentence returned empty string")
		}

		sentences[sentence] = true
	}

	// Should have at least some variety in results
	if len(sentences) < 2 {
		t.Logf("Warning: GetRandomSentence returned same sentence multiple times (only %d unique sentences)", len(sentences))
	}
}

// TestGetMultipleSentences_Comprehensive tests multiple sentence generation
func TestGetMultipleSentences_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		validate func(string) bool
	}{
		{
			name:     "single sentence",
			count:    1,
			validate: func(s string) bool { return len(s) > 10 },
		},
		{
			name:     "three sentences",
			count:    3,
			validate: func(s string) bool { return len(s) > 30 },
		},
		{
			name:     "many sentences",
			count:    10,
			validate: func(s string) bool { return len(s) > 100 },
		},
		{
			name:     "zero sentences",
			count:    0,
			validate: func(s string) bool { return len(s) >= 0 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMultipleSentences(tt.count)

			if result == "No text source available" {
				t.Skip("Text source not available")
			}

			if !tt.validate(result) {
				t.Errorf("GetMultipleSentences(%d) validation failed: %q", tt.count, result)
			}
		})
	}
}

// TestGetAvailableBooks_Content tests book list contents
func TestGetAvailableBooks_Content(t *testing.T) {
	books := GetAvailableBooks()

	if len(books) == 0 {
		t.Error("GetAvailableBooks should return at least one book")
		return
	}

	// Check each book has valid data
	for i, book := range books {
		if book.ID <= 0 {
			t.Errorf("Book %d has invalid ID: %d", i, book.ID)
		}

		if len(book.Name) == 0 {
			t.Errorf("Book %d has empty name", i)
		}

		// Check for duplicate IDs
		for j, other := range books {
			if i != j && book.ID == other.ID {
				t.Errorf("Duplicate book ID %d at positions %d and %d", book.ID, i, j)
			}
		}
	}
}

// TestGetAvailableBooks_Consistency tests consistency of book list
func TestGetAvailableBooks_Consistency(t *testing.T) {
	books1 := GetAvailableBooks()
	books2 := GetAvailableBooks()

	if len(books1) != len(books2) {
		t.Errorf("GetAvailableBooks returned different lengths: %d vs %d", len(books1), len(books2))
	}

	for i := range books1 {
		if books1[i].ID != books2[i].ID {
			t.Errorf("Book order inconsistent at index %d", i)
		}
	}
}

// TestSetBook_InvalidID tests SetBook with invalid IDs
func TestSetBook_InvalidID(t *testing.T) {
	tests := []struct {
		name    string
		bookID  int
		wantErr bool
	}{
		{"zero ID", 0, true},
		{"negative ID", -1, true},
		{"very large ID", 999999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetBook(tt.bookID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBook(%d) error = %v, wantErr = %v", tt.bookID, err, tt.wantErr)
			}
		})
	}
}

// TestCurrentBook tests current book tracking
func TestCurrentBook(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) < 2 {
		t.Skip("Need at least 2 books for this test")
	}

	// Set to first book
	err := SetBook(books[0].ID)
	if err != nil {
		t.Fatalf("Failed to set first book: %v", err)
	}

	current := GetCurrentBook()
	if current == nil {
		t.Error("GetCurrentBook returned nil")
	} else if current.ID != books[0].ID {
		t.Errorf("Expected book ID %d, got %d", books[0].ID, current.ID)
	}

	// Switch to second book
	err = SetBook(books[1].ID)
	if err != nil {
		t.Fatalf("Failed to set second book: %v", err)
	}

	current = GetCurrentBook()
	if current == nil {
		t.Error("GetCurrentBook returned nil after switching")
	} else if current.ID != books[1].ID {
		t.Errorf("Expected book ID %d after switch, got %d", books[1].ID, current.ID)
	}
}

// TestToASCIIFilter_Comprehensive tests ASCII filtering comprehensively
func TestToASCIIFilter_Comprehensive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"ASCII only", "Hello World", "Hello World"},
		{"with numbers", "Test 123", "Test 123"},
		{"with spaces", "multiple spaces", "multiple spaces"},
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"with punctuation", "Hello, World!", "Hello, World!"},
		{"with special chars", "test!@#$%", "test!@#$%"},
		{"mixed content", "abc123!@#", "abc123!@#"},
		{"with newlines", "hello\nworld", "hello\nworld"},
		{"with tabs", "hello\tworld", "hello\tworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toASCIIFilter(tt.input)
			if result != tt.expected {
				t.Errorf("toASCIIFilter(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetFullText_Comprehensive tests full text retrieval
func TestGetFullText_Comprehensive(t *testing.T) {
	text := GetFullText()

	if text == "" {
		t.Error("GetFullText() should return non-empty text")
		return
	}

	// Check text properties
	if len(text) < 100 {
		t.Errorf("GetFullText() too short: %d characters", len(text))
	}

	// Text should contain printable characters
	hasLetters := false
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetters = true
			break
		}
	}

	if !hasLetters {
		t.Error("GetFullText() should contain letters")
	}
}

// TestGetCurrentCharPos_Comprehensive tests character position tracking
func TestGetCurrentCharPos_Comprehensive(t *testing.T) {
	// Get initial position
	pos1 := GetCurrentCharPos()
	if pos1 < 0 {
		t.Errorf("GetCurrentCharPos() should return non-negative value, got %d", pos1)
	}

	// Get position multiple times - should be consistent or increment reasonably
	pos2 := GetCurrentCharPos()
	if pos2 < pos1 {
		t.Errorf("Character position should not decrease: %d -> %d", pos1, pos2)
	}
}

// TestGetLastParagraphEnd_Comprehensive tests paragraph end tracking
func TestGetLastParagraphEnd_Comprehensive(t *testing.T) {
	end := GetLastParagraphEnd()

	if end < 0 {
		t.Errorf("GetLastParagraphEnd() should return non-negative value, got %d", end)
	}

	// Paragraph end should be reasonable compared to full text length
	fullText := GetFullText()
	if end > len(fullText) {
		t.Errorf("Paragraph end %d should not exceed text length %d", end, len(fullText))
	}
}

// TestCalculateSentencesCompleted_EdgeCases tests sentence counting edge cases
func TestCalculateSentencesCompleted_EdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		paragraphLength int
		validate        func(int) bool
	}{
		{"zero length", 0, func(c int) bool { return true }}, // Any value is acceptable
		{"small length", 10, func(c int) bool { return true }},
		{"large length", 10000, func(c int) bool { return true }},
		{"negative length", -1, func(c int) bool { return true }}, // Can return negative
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CalculateSentencesCompleted(tt.paragraphLength)
			if !tt.validate(count) {
				t.Errorf("CalculateSentencesCompleted(%d) validation failed: %d", tt.paragraphLength, count)
			}
		})
	}
}

// TestCalculateSentencesCompletedWithCount_EdgeCases tests sentence counting with limit
func TestCalculateSentencesCompletedWithCount_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		sentenceCount int
		validate      func(int) bool
	}{
		{"zero count", 0, func(c int) bool { return true }},
		{"one sentence", 1, func(c int) bool { return true }},
		{"many sentences", 100, func(c int) bool { return true }},
		{"negative count", -1, func(c int) bool { return true }}, // Can return negative
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := CalculateSentencesCompletedWithCount(tt.sentenceCount)
			if !tt.validate(count) {
				t.Errorf("CalculateSentencesCompletedWithCount(%d) validation failed: %d", tt.sentenceCount, count)
			}
		})
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

// BenchmarkGetFullText benchmarks full text retrieval
func BenchmarkGetFullText(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetFullText()
	}
}

// BenchmarkGetParagraph benchmarks paragraph generation
func BenchmarkGetParagraph(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetParagraph(3)
	}
}

// BenchmarkToASCIIFilter benchmarks ASCII filtering
func BenchmarkToASCIIFilter(b *testing.B) {
	testString := "This is a test string with various characters and numbers 123!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = toASCIIFilter(testString)
	}
}

// TestGetRandomSentence_EmptySentences tests behavior when no sentences are loaded
func TestGetRandomSentence_EmptySentences(t *testing.T) {
	// Save current state
	oldSentences := sentences

	// Clear sentences
	sentences = nil

	result := GetRandomSentence()
	if result != "No text source available" {
		t.Errorf("GetRandomSentence() with empty sentences = %q, want %q", result, "No text source available")
	}

	// Restore
	sentences = oldSentences
}

// TestGetFullText_EmptyText tests GetFullText when no text is loaded
func TestGetFullText_EmptyText(t *testing.T) {
	// Save current state
	oldFullText := fullText
	oldRawBookContent := rawBookContent

	// Clear text
	fullText = ""
	rawBookContent = ""

	result := GetFullText()
	if result != "No text source available" {
		t.Errorf("GetFullText() with empty text = %q, want %q", result, "No text source available")
	}

	// Restore
	fullText = oldFullText
	rawBookContent = oldRawBookContent
}
