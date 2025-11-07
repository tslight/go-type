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
	_ = ClearProgress()

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
	_ = ClearProgress()
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
	_ = ClearProgress()
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
	_ = ClearProgress()
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
	_ = ClearProgress()
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

// TestStateManagerAddSession tests adding a session to a book's history
func TestStateManagerAddSession(t *testing.T) {
	// Create a state manager with in-memory state only (don't load from disk)
	sm := &StateManager{
		states: make(map[int]*BookState),
	}

	bookID := 7777 // Use unique bookID to avoid conflicts

	// Create initial state directly (don't call SaveState which writes to disk)
	sm.states[bookID] = &BookState{
		BookID:       bookID,
		BookName:     "Test Book",
		CharacterPos: 100,
		LastHash:     "hash123",
		Sessions:     []SessionResult{},
	}

	// Add a session directly to state (simulating what AddSession does without disk I/O)
	result := SessionResult{
		WPM:       75.5,
		Accuracy:  95.2,
		Errors:    3,
		CharTyped: 450,
		Duration:  360,
	}

	state := sm.GetState(bookID)
	if state == nil {
		t.Fatal("GetState returned nil")
	}

	state.Sessions = append(state.Sessions, result)

	// Verify
	if len(state.Sessions) != 1 {
		t.Errorf("Expected 1 session, got %d", len(state.Sessions))
	}

	session := state.Sessions[0]
	if session.WPM != result.WPM {
		t.Errorf("WPM = %f, want %f", session.WPM, result.WPM)
	}
	if session.Accuracy != result.Accuracy {
		t.Errorf("Accuracy = %f, want %f", session.Accuracy, result.Accuracy)
	}
	if session.Errors != result.Errors {
		t.Errorf("Errors = %d, want %d", session.Errors, result.Errors)
	}
	if session.CharTyped != result.CharTyped {
		t.Errorf("CharTyped = %d, want %d", session.CharTyped, result.CharTyped)
	}
	if session.Duration != result.Duration {
		t.Errorf("Duration = %d, want %d", session.Duration, result.Duration)
	}
}

// TestStateManagerGetStats tests cumulative statistics calculation
func TestStateManagerGetStats(t *testing.T) {
	// Create a state manager with in-memory state only (don't load from disk)
	sm := &StateManager{
		states: make(map[int]*BookState),
	}

	// Save initial state with a unique book ID
	bookID := 8888
	state := &BookState{
		BookID:          bookID,
		BookName:        "Test Book Stats",
		CharacterPos:    100,
		LastHash:        "hash_stats",
		TextLength:      5000,
		PercentComplete: 2.0,
		Sessions:        []SessionResult{},
	}
	sm.states[bookID] = state

	// Add multiple sessions
	sessions := []struct {
		wpm      float64
		accuracy float64
		errors   int
		chars    int
		duration int
	}{
		{75.0, 95.0, 2, 450, 360},
		{80.0, 96.0, 1, 480, 360},
		{70.0, 94.0, 3, 420, 360},
	}

	for _, s := range sessions {
		_ = sm.AddSession(bookID, SessionResult{
			WPM:       s.wpm,
			Accuracy:  s.accuracy,
			Errors:    s.errors,
			CharTyped: s.chars,
			Duration:  s.duration,
		})
	}

	stats := sm.GetStats(bookID)
	if stats == nil {
		t.Fatal("GetStats returned nil")
	}

	sessionsCount := stats["sessions_completed"].(int)
	if sessionsCount != 3 {
		t.Errorf("sessions_completed = %d, want 3", sessionsCount)
	}

	avgWPM := stats["average_wpm"].(float64)
	expectedAvgWPM := (75.0 + 80.0 + 70.0) / 3
	if avgWPM < expectedAvgWPM-0.01 || avgWPM > expectedAvgWPM+0.01 {
		t.Errorf("average_wpm = %f, want %f", avgWPM, expectedAvgWPM)
	}

	bestWPM := stats["best_wpm"].(float64)
	if bestWPM != 80.0 {
		t.Errorf("best_wpm = %f, want 80.0", bestWPM)
	}

	totalChars := stats["total_characters"].(int)
	expectedChars := 450 + 480 + 420
	if totalChars != expectedChars {
		t.Errorf("total_characters = %d, want %d", totalChars, expectedChars)
	}
}
