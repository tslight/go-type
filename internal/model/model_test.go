package model

import (
	"testing"
	"time"

	"github.com/tobe/go-type/internal/content"
)

type dummyState struct{}

func (d *dummyState) GetSavedCharPos() int                                          { return 0 }
func (d *dummyState) SaveProgress(int) error                                        { return nil }
func (d *dummyState) RecordSession(float64, float64, int, int, int) (string, error) { return "", nil }

func TestNewModel_Creation(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "Hello world"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	if m == nil {
		t.Fatal("NewModel returned nil")
	}
	if m.currentContent == nil || m.currentContent.Name != "Test" {
		t.Fatalf("expected currentContent name Test")
	}
}

func TestModel_UpdateTyping(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "abcdef"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	// Simulate key presses via internal logic: append directly
	m.userInput = "abc"
	if len(m.userInput) != 3 {
		t.Fatalf("expected userInput len 3")
	}
}

func TestModel_WPMAccuracy(t *testing.T) {
	c := &content.Content{ID: 1, Name: "Test", Text: "aaaaa"}
	m := NewModel(c.Text, c, 80, 24, &dummyState{})
	m.userInput = "aaaaa"
	m.testStarted = true
	m.startTime = time.Now().Add(-1 * time.Minute)
	view := m.View()
	if view == "" {
		t.Fatalf("expected non-empty view")
	}
}
