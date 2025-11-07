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
	if m.SelectedContent() != nil { /* initial selection allowed to be nil */
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
