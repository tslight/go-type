package runner

import (
	"bytes"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/internal/selection"
)

func TestRunApp_VersionFlags(t *testing.T) {
	var out bytes.Buffer
	err := RunApp(AppConfig{
		Name:      "app",
		Version:   "1.2.3",
		Args:      []string{"-v"},
		Stdout:    &out,
		Stderr:    &out,
		ListItems: func() ([]string, error) { return nil, nil },
		SelectAndLoad: func(int, int) (*selection.Selection, error) {
			return nil, nil
		},
	})
	if err != nil {
		t.Fatalf("RunApp returned error on version flag: %v", err)
	}
	if got := out.String(); got != "1.2.3\n" {
		t.Fatalf("unexpected version output: %q", got)
	}
}

func TestRunApp_ListFlag(t *testing.T) {
	var out bytes.Buffer
	err := RunApp(AppConfig{
		Name:    "app",
		Version: "1.0.0",
		Args:    []string{"-l"},
		Stdout:  &out,
		Stderr:  &out,
		ListItems: func() ([]string, error) {
			return []string{"one", "two"}, nil
		},
		SelectAndLoad: func(int, int) (*selection.Selection, error) { return nil, nil },
	})
	if err != nil {
		t.Fatalf("RunApp returned error on list flag: %v", err)
	}
	want := "one\n two\n"
	// Allow either with or without leading space depending on formatting
	got := out.String()
	if got != "one\n"+"two\n" && got != want {
		t.Fatalf("unexpected list output: %q", got)
	}
}

func TestRunApp_HandlersRequired(t *testing.T) {
	if err := RunApp(AppConfig{}); err == nil {
		t.Fatalf("expected error when handlers missing")
	}
	// Missing ListItems
	if err := RunApp(AppConfig{SelectAndLoad: func(int, int) (*selection.Selection, error) { return nil, nil }}); err == nil {
		t.Fatalf("expected error when ListItems missing")
	}
}

func TestRunApp_AbortSelection(t *testing.T) {
	var out bytes.Buffer
	err := RunApp(AppConfig{
		Name:      "app",
		Version:   "0.0.1",
		Args:      []string{},
		Stdout:    &out,
		Stderr:    &out,
		ListItems: func() ([]string, error) { return []string{"x"}, nil },
		SelectAndLoad: func(int, int) (*selection.Selection, error) {
			return nil, nil // user aborted selection
		},
	})
	if err != nil {
		t.Fatalf("RunApp returned error when selection aborted: %v", err)
	}
}

func TestRunApp_ListItemsError(t *testing.T) {
	var out bytes.Buffer
	err := RunApp(AppConfig{
		Name:          "app",
		Version:       "0.0.1",
		Args:          []string{"--list"},
		Stdout:        &out,
		Stderr:        &out,
		ListItems:     func() ([]string, error) { return nil, errors.New("boom") },
		SelectAndLoad: func(int, int) (*selection.Selection, error) { return nil, nil },
	})
	if err == nil {
		t.Fatalf("expected error from ListItems to propagate")
	}
}

type fakeProvider struct{}

func (f *fakeProvider) GetSavedCharPos() int                                          { return 0 }
func (f *fakeProvider) SaveProgress(int) error                                        { return nil }
func (f *fakeProvider) RecordSession(float64, float64, int, int, int) (string, error) { return "", nil }

func TestRunApp_NormalRun(t *testing.T) {
	// Stub model program to avoid interactive TUI
	orig := runModelProgram
	defer func() { runModelProgram = orig }()
	runModelProgram = func(m tea.Model) (tea.Model, error) { return m, nil }
	// Build a fake selection
	sel := &selection.Selection{Text: "hello world", Content: &content.Content{ID: 1, Name: "Test", Text: "hello world"}, Provider: &fakeProvider{}}
	var out bytes.Buffer
	err := RunApp(AppConfig{
		Name:      "app",
		Version:   "0.0.1",
		Args:      []string{},
		Stdout:    &out,
		Stderr:    &out,
		ListItems: func() ([]string, error) { return []string{"x"}, nil },
		SelectAndLoad: func(int, int) (*selection.Selection, error) {
			return sel, nil
		},
	})
	if err != nil {
		t.Fatalf("expected normal run without error, got %v", err)
	}
}
