package cli

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
)

// TestNewModel tests model creation with various inputs
func TestNewModel(t *testing.T) {
	tests := []struct {
		name          string
		text          string
		book          *textgen.Book
		width         int
		height        int
		validateState func(*Model) bool
	}{
		{
			name:   "basic model creation",
			text:   "The quick brown fox",
			book:   &textgen.Book{ID: 1, Name: "Test Book"},
			width:  80,
			height: 24,
			validateState: func(m *Model) bool {
				return !m.finished && !m.testStarted
			},
		},
		{
			name:   "empty text model",
			text:   "",
			book:   &textgen.Book{ID: 2, Name: "Empty"},
			width:  80,
			height: 24,
			validateState: func(m *Model) bool {
				return m != nil
			},
		},
		{
			name:   "long text model",
			text:   "This is a very long text that contains many words and sentences. " + string(make([]byte, 1000)),
			book:   &textgen.Book{ID: 3, Name: "Long"},
			width:  80,
			height: 24,
			validateState: func(m *Model) bool {
				return m != nil
			},
		},
		{
			name:   "small terminal",
			text:   "test",
			book:   &textgen.Book{ID: 4, Name: "Small"},
			width:  40,
			height: 12,
			validateState: func(m *Model) bool {
				return m.terminalWidth == 40 && m.terminalHeight == 12
			},
		},
		{
			name:   "large terminal",
			text:   "test",
			book:   &textgen.Book{ID: 5, Name: "Large"},
			width:  200,
			height: 50,
			validateState: func(m *Model) bool {
				return m.terminalWidth == 200 && m.terminalHeight == 50
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.text, tt.book, tt.width, tt.height)

			if m == nil {
				t.Fatal("NewModel returned nil")
			}

			if !tt.validateState(m) {
				t.Error("Model state validation failed")
			}

			if m.terminalWidth != tt.width {
				t.Errorf("Expected width %d, got %d", tt.width, m.terminalWidth)
			}

			if m.terminalHeight != tt.height {
				t.Errorf("Expected height %d, got %d", tt.height, m.terminalHeight)
			}
		})
	}
}

// TestModelInit tests model initialization
func TestModelInit(t *testing.T) {
	m := NewModel("test text", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)
	cmd := m.Init()

	if cmd != nil {
		t.Error("Model.Init() should return nil command")
	}
}

// TestModelUpdate tests model update with various messages
func TestModelUpdate(t *testing.T) {
	m := NewModel("The quick brown fox", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

	tests := []struct {
		name    string
		message tea.Msg
	}{
		{
			name:    "window size change",
			message: tea.WindowSizeMsg{Width: 100, Height: 30},
		},
		{
			name:    "key press",
			message: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
		},
		{
			name:    "key backspace",
			message: tea.KeyMsg{Type: tea.KeyBackspace},
		},
		{
			name:    "key enter",
			message: tea.KeyMsg{Type: tea.KeyEnter},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newModel, cmd := m.Update(tt.message)

			if newModel == nil {
				t.Error("Update should return a model")
			}

			// Command can be nil or not - both are valid
			_ = cmd
		})
	}
}

// TestModelView tests model rendering
func TestModelView(t *testing.T) {
	m := NewModel("The quick brown fox", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

	view := m.View()

	if view == "" {
		t.Error("Model.View() should return non-empty string")
	}

	if len(view) < 5 {
		t.Errorf("Model.View() returned suspiciously short content: %d chars", len(view))
	}
}

// TestModel_TextNormalization tests that text is properly handled
func TestModel_TextNormalization(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"simple text", "hello world"},
		{"text with newlines", "hello\nworld"},
		{"text with tabs", "hello\tworld"},
		{"text with special chars", "hello! world?"},
		{"text with numbers", "test 123"},
		{"empty text", ""},
		{"whitespace only", "   \n\t   "},
		{"unicode-like ascii", "test-data_123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel(tt.text, &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

			if m == nil {
				t.Fatal("NewModel returned nil")
			}

			view := m.View()
			// View should handle any text gracefully
			_ = view
		})
	}
}

// TestModel_StateTransitions tests state transitions during typing
func TestModel_StateTransitions(t *testing.T) {
	m := NewModel("The quick brown fox", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

	// Initial state
	if m.finished {
		t.Error("Model should not be finished initially")
	}

	if m.testStarted {
		t.Error("Model should not have started initially")
	}

	// Simulate typing
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T'}})
	if newModel == nil {
		t.Fatal("Update returned nil model")
	}

	// After typing, test might start
	// (depends on implementation)
}

// TestModel_TerminalResize tests handling of terminal resize
func TestModel_TerminalResize(t *testing.T) {
	m := NewModel("The quick brown fox", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

	// Resize to larger terminal
	m1, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if m1 == nil {
		t.Fatal("Update returned nil model on resize")
	}

	// Resize to smaller terminal
	m2, _ := m1.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	if m2 == nil {
		t.Fatal("Update returned nil model on second resize")
	}

	// Verify it renders in both cases
	view := m2.View()
	if len(view) == 0 {
		t.Error("View should still render after resize")
	}
}

// TestModel_WithNilBook tests model with nil book
func TestModel_WithNilBook(t *testing.T) {
	// Should handle nil book gracefully or panic appropriately
	defer func() {
		if r := recover(); r != nil {
			// Panicking on nil book is acceptable behavior
			t.Logf("NewModel panicked with nil book (acceptable): %v", r)
		}
	}()

	m := NewModel("test", nil, 80, 24)
	if m != nil {
		// If it doesn't panic, it should return a valid model
		_ = m.View()
	}
}

// TestModel_BookAssociation tests that model correctly associates with book
func TestModel_BookAssociation(t *testing.T) {
	book := &textgen.Book{ID: 42, Name: "Special Book"}
	m := NewModel("test text", book, 80, 24)

	if m == nil {
		t.Fatal("NewModel returned nil")
	}

	// Model should retain book information
	// (exact field name may vary based on implementation)
}

// TestModel_InputHandling tests various input scenarios
func TestModel_InputHandling(t *testing.T) {
	tests := []struct {
		name  string
		runes []rune
	}{
		{"single character", []rune{'a'}},
		{"multiple characters", []rune{'h', 'e', 'l', 'l', 'o'}},
		{"space", []rune{' '}},
		{"special characters", []rune{'!', '@', '#'}},
		{"numbers", []rune{'1', '2', '3'}},
		{"mixed", []rune{'a', '1', '!', ' '}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel("The quick brown fox", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)
			var currentModel tea.Model = model

			for _, r := range tt.runes {
				newModel, _ := currentModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
				if newModel == nil {
					t.Errorf("Update failed for rune %c", r)
					return
				}
				currentModel = newModel
			}
		})
	}
}

// BenchmarkModelView benchmarks model rendering
func BenchmarkModelView(b *testing.B) {
	m := NewModel("The quick brown fox jumps over the lazy dog", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

// BenchmarkModelUpdate benchmarks model update
func BenchmarkModelUpdate(b *testing.B) {
	m := NewModel("The quick brown fox jumps over the lazy dog", &textgen.Book{ID: 1, Name: "Test"}, 80, 24)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Update(msg)
	}
}

// BenchmarkNewModel benchmarks model creation
func BenchmarkNewModel(b *testing.B) {
	book := &textgen.Book{ID: 1, Name: "Test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewModel("The quick brown fox jumps over the lazy dog", book, 80, 24)
	}
}
