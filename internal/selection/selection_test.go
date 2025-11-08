package selection

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/assets/books"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/internal/menu"
)

func TestSelectContent_NilManager(t *testing.T) {
	res, err := SelectContent(nil, 80, 24)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil selection when manager is nil")
	}
}

func TestSelectContent_ProgramStub(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := content.NewContentManager(books.EFS, "test-books", true)
	// pick the first available item and force selection by simulating 'enter' behavior
	// We can't directly set selectedContent from here, so we'll rely on SetContent after selection returns.
	called := false
	orig := runMenuProgram
	defer func() { runMenuProgram = orig }()
	runMenuProgram = func(m tea.Model) (tea.Model, error) {
		called = true
		// Do not change the model; this simulates user pressing 'q' or not selecting, so result is nil
		return m, nil
	}
	res, err := SelectContent(cm, 80, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected runMenuProgram to be called")
	}
	// Selection is allowed to be nil in this path; assert no panic.
	_ = res
}

func TestSelectContent_SelectAndProvider(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := content.NewContentManager(books.EFS, "test-books", true)
	called := false
	orig := runMenuProgram
	defer func() { runMenuProgram = orig }()
	runMenuProgram = func(m tea.Model) (tea.Model, error) {
		called = true
		if mm, ok := m.(*menu.MenuModel); ok {
			// simulate pressing enter to select the first item
			nm, _ := mm.Update(tea.KeyMsg{Type: tea.KeyEnter})
			return nm, nil
		}
		return m, nil
	}
	sel, err := SelectContent(cm, 80, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatalf("expected runMenuProgram to be called")
	}
	if sel == nil || sel.Content == nil || sel.Provider == nil {
		t.Fatalf("expected non-nil selection, content, and provider")
	}
	// Exercise provider methods
	if err := sel.Provider.SaveProgress(10); err != nil {
		t.Fatalf("SaveProgress error: %v", err)
	}
	if got := sel.Provider.GetSavedCharPos(); got != 10 {
		t.Fatalf("expected saved char pos 10, got %d", got)
	}
	stats, err := sel.Provider.RecordSession(50.0, 95.0, 1, 100, 90, 30)
	if err != nil {
		t.Fatalf("RecordSession error: %v", err)
	}
	if stats == "" || !strings.Contains(stats, "STATISTICS") {
		t.Fatalf("expected formatted statistics, got %q", stats)
	}
}
