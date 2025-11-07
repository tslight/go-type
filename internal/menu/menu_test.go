package menu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/assets/books"
	"github.com/tobe/go-type/internal/content"
)

func newTestManager() *content.ContentManager {
	return content.NewContentManager(books.EFS, "test-gutentype", true)
}

func TestNewMenuModel_Basic(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	if m == nil {
		t.Fatal("NewMenuModel returned nil")
	}
	// Explicitly call Init to cover that path.
	if initCmd := m.Init(); initCmd != nil { _ = initCmd() }
	// SelectedContent should be nil initially; assert explicitly.
	if m.SelectedContent() != nil {
		t.Fatalf("expected no selection on initialization")
	}
	if v := m.View(); v == "" {
		t.Fatalf("expected non-empty view")
	}
	_ = m.View()
}

func TestMenuModel_HandleResize(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	nm, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	if nm == nil {
		t.Fatal("Update returned nil model")
	}
}

func TestMenuModel_SearchAndSelect(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// Enter search mode
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Type some query letters (use runes that likely exist like 'a')
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Execute search
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Navigate to next match if any
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Ensure view renders without panic
	if v := m.View(); v == "" {
		t.Fatalf("expected non-empty view after search")
	}
}

func TestMenuModel_ShowStatsView(t *testing.T) {
	cm := newTestManager()
	m := NewMenuModel(cm, 80, 24)
	// Basic sanity check the viewport dimensions are set
	if m.viewport.Width <= 0 || m.viewport.Height <= 0 {
		t.Fatalf("viewport not properly initialized")
	}
	// Toggle stats view
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	v := m.View()
	if v == "" {
		t.Fatalf("expected stats view to render")
	}
	// Exit stats view
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.showingStats {
		t.Fatalf("expected stats view to be closed after 'q'")
	}
}

func TestMenuModel_EnterSelect(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// ensure we have items
	if len(m.items) == 0 {
		t.Skip("no items embedded for selection test")
	}
	// trigger enter selection
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.SelectedContent() == nil {
		t.Fatalf("expected SelectedContent after enter")
	}
	if !m.done {
		t.Fatalf("expected done=true after selection")
	}
}

func TestMenuModel_BackwardSearchAndPrev(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// Enter backward search
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Type query 'e'
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Execute search
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Previous result navigation 'N'
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if v := m.View(); v == "" {
		t.Fatalf("expected non-empty view after backward search")
	}
}

func TestMenuModel_StatsEscExit(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// Enter stats view
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Exit with esc
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.showingStats {
		t.Fatalf("expected stats view closed after esc")
	}
}
