package textgen

import (
	"testing"
)

// TestSaveProgress tests progress saving functionality
func TestSaveProgress(t *testing.T) {
	tests := []struct {
		name      string
		charPos   int
		hash      string
		shouldErr bool
	}{
		{
			name:      "save position zero with empty hash",
			charPos:   0,
			hash:      "",
			shouldErr: false,
		},
		{
			name:      "save position with hash",
			charPos:   100,
			hash:      "test-hash-123",
			shouldErr: false,
		},
		{
			name:      "save large position",
			charPos:   1000000,
			hash:      "hash",
			shouldErr: false,
		},
		{
			name:      "save negative position",
			charPos:   -1,
			hash:      "hash",
			shouldErr: false, // May or may not be rejected depending on implementation
		},
		{
			name:      "save with very long hash",
			charPos:   42,
			hash:      string(make([]byte, 10000)),
			shouldErr: false,
		},
	}

	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set book first
			err := SetBook(books[0].ID)
			if err != nil {
				t.Fatalf("Failed to set book: %v", err)
			}

			// Test SaveProgress
			err = SaveProgress(tt.charPos, tt.hash)
			if (err != nil) != tt.shouldErr {
				t.Errorf("SaveProgress(%d, %q) error = %v, wantErr = %v", tt.charPos, tt.hash, err, tt.shouldErr)
			}
		})
	}
}

// TestGetProgress tests progress retrieval
func TestGetProgress(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	book := &books[0]

	// Set up: Set book and save progress
	err := SetBook(book.ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	testPos := 42
	testHash := "test-hash-value"
	err = SaveProgress(testPos, testHash)
	if err != nil {
		t.Fatalf("Failed to save progress: %v", err)
	}

	// Test GetProgress
	progress := GetProgress()

	if progress == nil {
		t.Error("GetProgress returned nil when progress was saved")
		return
	}

	if progress.CharacterPos != testPos {
		t.Errorf("Expected CharacterPos %d, got %d", testPos, progress.CharacterPos)
	}

	if progress.LastHash != testHash {
		t.Errorf("Expected LastHash %q, got %q", testHash, progress.LastHash)
	}

	if progress.BookID != book.ID {
		t.Errorf("Expected BookID %d, got %d", book.ID, progress.BookID)
	}

	if progress.BookName != book.Name {
		t.Errorf("Expected BookName %q, got %q", book.Name, progress.BookName)
	}
}

// TestGetProgress_NoSavedProgress tests GetProgress with no saved progress
func TestGetProgress_NoSavedProgress(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	// Clear any progress first
	ClearProgress()

	// SetBook to a book we know hasn't been used
	err := SetBook(books[0].ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	// GetProgress with no saved progress
	progress := GetProgress()

	// Should return nil or empty progress
	if progress != nil && progress.CharacterPos != 0 {
		t.Errorf("Expected empty progress when nothing saved, got %v", progress)
	}
}

// TestClearProgress_Functionality tests progress clearing functionality
func TestClearProgress_Functionality(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	book := &books[0]

	// Setup: Set book and save progress
	err := SetBook(book.ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	err = SaveProgress(100, "test-hash")
	if err != nil {
		t.Fatalf("Failed to save progress: %v", err)
	}

	// Verify progress was saved
	progress1 := GetProgress()
	if progress1 == nil || progress1.CharacterPos != 100 {
		t.Fatalf("Progress not saved correctly")
	}

	// Clear progress
	err = ClearProgress()
	if err != nil {
		t.Errorf("ClearProgress failed: %v", err)
	}

	// Verify progress was cleared
	progress2 := GetProgress()
	if progress2 != nil && progress2.CharacterPos != 0 {
		t.Errorf("ClearProgress: expected cleared progress, got %v", progress2)
	}
}

// TestClearProgress_MultipleTimes tests clearing progress multiple times
func TestClearProgress_MultipleTimes(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	err := SetBook(books[0].ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	// Clear multiple times without saving - should not error
	for i := 0; i < 3; i++ {
		err = ClearProgress()
		if err != nil {
			t.Errorf("ClearProgress iteration %d failed: %v", i, err)
		}
	}
}

// TestGetProgressForBook_BasicFunctionality tests retrieving progress for specific book
func TestGetProgressForBook_BasicFunctionality(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) < 2 {
		t.Skip("Need at least 2 books for this test")
	}

	book1 := &books[0]
	book2 := &books[1]

	// Setup: Save progress for book 1
	err := SetBook(book1.ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	err = SaveProgress(50, "hash1")
	if err != nil {
		t.Fatalf("Failed to save progress: %v", err)
	}

	// Test GetProgressForBook with correct book
	progress := GetProgressForBook(book1)

	if progress == nil {
		t.Error("GetProgressForBook returned nil for book with saved progress")
	} else {
		if progress.CharacterPos != 50 {
			t.Errorf("Expected CharacterPos 50, got %d", progress.CharacterPos)
		}
		if progress.LastHash != "hash1" {
			t.Errorf("Expected hash hash1, got %q", progress.LastHash)
		}
	}

	// Test GetProgressForBook with different book (should have no progress)
	progress2 := GetProgressForBook(book2)

	if progress2 != nil && progress2.CharacterPos != 0 {
		t.Errorf("Expected empty progress for unused book, got %v", progress2)
	}

	// Cleanup
	ClearProgress()
}

// TestGetProgressForBook_MultipleBooks tests managing progress for multiple books
func TestGetProgressForBook_MultipleBooks(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) < 2 {
		t.Skip("Need at least 2 books for this test")
	}

	book1 := &books[0]
	book2 := &books[1]

	// Save progress for book 1
	err := SetBook(book1.ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}
	err = SaveProgress(25, "hash-b1")
	if err != nil {
		t.Fatalf("Failed to save progress: %v", err)
	}

	// Save progress for book 2
	err = SetBook(book2.ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}
	err = SaveProgress(75, "hash-b2")
	if err != nil {
		t.Fatalf("Failed to save progress: %v", err)
	}

	// Verify both books have their respective progress
	p1 := GetProgressForBook(book1)
	p2 := GetProgressForBook(book2)

	if p1 == nil || p1.CharacterPos != 25 {
		t.Errorf("Book 1: expected CharacterPos 25, got %v", p1)
	}

	if p2 == nil || p2.CharacterPos != 75 {
		t.Errorf("Book 2: expected CharacterPos 75, got %v", p2)
	}

	// Cleanup
	ClearProgress()
}

// TestProgressPersistence tests that progress changes are persistent
func TestProgressPersistence(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	err := SetBook(books[0].ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	// Save initial progress
	err = SaveProgress(10, "initial")
	if err != nil {
		t.Fatalf("Failed to save initial progress: %v", err)
	}

	// Update progress
	err = SaveProgress(20, "updated")
	if err != nil {
		t.Fatalf("Failed to update progress: %v", err)
	}

	// Verify updated progress
	progress := GetProgress()
	if progress.CharacterPos != 20 {
		t.Errorf("Expected updated CharacterPos 20, got %d", progress.CharacterPos)
	}
	if progress.LastHash != "updated" {
		t.Errorf("Expected updated hash, got %q", progress.LastHash)
	}

	// Cleanup
	ClearProgress()
}

// TestProgressWithZeroValues tests progress with zero and edge values
func TestProgressWithZeroValues(t *testing.T) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		t.Skip("No books available for testing")
	}

	err := SetBook(books[0].ID)
	if err != nil {
		t.Fatalf("Failed to set book: %v", err)
	}

	// Save progress at position 0
	err = SaveProgress(0, "")
	if err != nil {
		t.Fatalf("Failed to save progress at 0: %v", err)
	}

	progress := GetProgress()
	if progress != nil && progress.CharacterPos != 0 {
		t.Errorf("Expected CharacterPos 0, got %d", progress.CharacterPos)
	}

	// Cleanup
	ClearProgress()
}

// BenchmarkSaveProgress benchmarks progress saving
func BenchmarkSaveProgress(b *testing.B) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		b.Skip("No books available")
	}

	err := SetBook(books[0].ID)
	if err != nil {
		b.Fatalf("Failed to set book: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SaveProgress(i, "test-hash")
	}
}

// BenchmarkGetProgress benchmarks progress retrieval
func BenchmarkGetProgress(b *testing.B) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		b.Skip("No books available")
	}

	err := SetBook(books[0].ID)
	if err != nil {
		b.Fatalf("Failed to set book: %v", err)
	}

	err = SaveProgress(42, "test-hash")
	if err != nil {
		b.Fatalf("Failed to save progress: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetProgress()
	}
}

// BenchmarkGetProgressForBook benchmarks book-specific progress retrieval
func BenchmarkGetProgressForBook(b *testing.B) {
	books := GetAvailableBooks()
	if len(books) == 0 {
		b.Skip("No books available")
	}

	book := &books[0]

	err := SetBook(book.ID)
	if err != nil {
		b.Fatalf("Failed to set book: %v", err)
	}

	err = SaveProgress(42, "test-hash")
	if err != nil {
		b.Fatalf("Failed to save progress: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetProgressForBook(book)
	}
}

// TestMigrateBookState tests backward compatibility migration
func TestMigrateBookState(t *testing.T) {
	tests := []struct {
		name     string
		inputPos int
		inputID  int
	}{
		{
			name:     "migrate zero position",
			inputPos: 0,
			inputID:  1,
		},
		{
			name:     "migrate non-zero position",
			inputPos: 500,
			inputID:  2,
		},
		{
			name:     "migrate large position",
			inputPos: 1000000,
			inputID:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BookState{
				BookID:       tt.inputID,
				CharacterPos: tt.inputPos,
			}

			// Call migrate - should be a no-op
			migrateBookState(bs)

			// Verify state is unchanged
			if bs.BookID != tt.inputID {
				t.Errorf("BookID changed after migration: got %d, want %d", bs.BookID, tt.inputID)
			}
			if bs.CharacterPos != tt.inputPos {
				t.Errorf("CharacterPos changed after migration: got %d, want %d", bs.CharacterPos, tt.inputPos)
			}
		})
	}
}

// TestGetProgress_NoCurrentBook tests GetProgress when no book is set
func TestGetProgress_NoCurrentBook(t *testing.T) {
	// This tests the coverage gap where currentBook is nil
	// We can't easily set currentBook to nil from tests, so this is a limitation
	// of the current architecture for testing entry points
	t.Skip("Requires mocking of package-level currentBook variable")
}

// TestClearProgress_NoCurrentBook tests ClearProgress when no book is set
func TestClearProgress_NoCurrentBook(t *testing.T) {
	t.Skip("Requires mocking of package-level currentBook variable")
}

// TestGetProgressForBook_NilBook tests GetProgressForBook with nil input
func TestGetProgressForBook_NilBook(t *testing.T) {
	result := GetProgressForBook(nil)
	if result != nil {
		t.Errorf("GetProgressForBook(nil) = %v, want nil", result)
	}
}
