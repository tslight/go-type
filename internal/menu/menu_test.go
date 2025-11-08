package menu

import (
	"strings"
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
	if initCmd := m.Init(); initCmd != nil {
		_ = initCmd()
	}
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

func TestMenuModel_GlobalStatsView(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// Trigger global stats via 'I'
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'I'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	v := m.View()
	if v == "" || !strings.Contains(v, "GLOBAL STATISTICS") {
		t.Fatalf("expected global statistics view, got %q", v)
	}
	// Exit global stats view with 'q'
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.showingGlobal {
		t.Fatalf("expected showingGlobal to be false after exit")
	}
}

func TestMenuModel_GlobalStatsEscExit(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	// Enter global stats view
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'I'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// Exit with ESC
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.showingGlobal {
		t.Fatalf("expected global stats view closed after ESC")
	}
}

func TestMenuModel_PageScrolling(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	if len(m.items) < m.viewport.Height+2 { // ensure enough items to scroll at least a page
		// Not enough items embedded; still exercise keys without failure
		_ = m.View()
		_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
		return
	}
	startIndex := m.selectedIndex
	startOffset := m.viewport.YOffset
	// Page forward
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.selectedIndex <= startIndex {
		t.Fatalf("expected selectedIndex to increase after 'f'")
	}
	if m.viewport.YOffset <= startOffset {
		t.Fatalf("expected viewport.YOffset to increase after 'f'")
	}
	// Page backward
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if m.selectedIndex >= startIndex+m.viewport.Height {
		t.Fatalf("expected selectedIndex to move back after 'b'")
	}
}

func TestMenuModel_PageScrolling_PgKeys(t *testing.T) {
	m := NewMenuModel(newTestManager(), 80, 24)
	if len(m.items) == 0 {
		t.Skip("no items to test paging")
	}
	// Ensure enough items; if not, still exercise without assertions that depend on movement.
	startIdx := m.selectedIndex
	startOff := m.viewport.YOffset
	// PageDown
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	if len(m.items) >= m.viewport.Height+2 { // only assert movement when sufficient items
		if m.selectedIndex <= startIdx {
			t.Fatalf("expected selectedIndex to increase after PgDown")
		}
		if m.viewport.YOffset <= startOff {
			t.Fatalf("expected viewport.YOffset to increase after PgDown")
		}
	}
	// PageUp
	if mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgUp}); mm != nil {
		if cast, ok := mm.(*MenuModel); ok {
			m = cast
		}
	}
	// After paging up, selectedIndex should not exceed bounds; loose check
	if m.selectedIndex < 0 || m.selectedIndex >= len(m.items) {
		t.Fatalf("selectedIndex out of bounds after PgUp")
	}
}
