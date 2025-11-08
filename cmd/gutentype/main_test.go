package main

import (
	"bytes"
	"testing"

	"github.com/tobe/go-type/internal/runner"
	"github.com/tobe/go-type/internal/selection"
)

func TestMain_GutentypeVersionFlag(t *testing.T) {
	var buf bytes.Buffer
	cfg := runner.AppConfig{
		Name:          "gutentype",
		Version:       "gt-version",
		Args:          []string{"--version"},
		Stdout:        &buf,
		ListItems:     func() ([]string, error) { return nil, nil },
		SelectAndLoad: func(w, h int) (*selection.Selection, error) { return nil, nil },
	}
	if err := runner.RunApp(cfg); err != nil {
		t.Fatalf("RunApp returned error: %v", err)
	}
	if buf.String() != "gt-version\n" {
		t.Fatalf("expected version output, got %q", buf.String())
	}
}

func TestMain_GutentypeMissingHandlers(t *testing.T) {
	// Missing SelectAndLoad should error
	cfg := runner.AppConfig{Name: "gutentype", Version: "x", ListItems: func() ([]string, error) { return nil, nil }}
	if err := runner.RunApp(cfg); err == nil {
		t.Fatalf("expected error when SelectAndLoad is nil")
	}
	// Missing ListItems should error
	cfg2 := runner.AppConfig{Name: "gutentype", Version: "x", SelectAndLoad: func(w, h int) (*selection.Selection, error) { return nil, nil }}
	if err := runner.RunApp(cfg2); err == nil {
		t.Fatalf("expected error when ListItems is nil")
	}
}
