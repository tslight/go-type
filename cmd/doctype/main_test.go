package main

import (
	"bytes"
	"testing"

	"github.com/tobe/go-type/internal/runner"
	"github.com/tobe/go-type/internal/selection"
)

// TestMainVersionFlag ensures the -v flag prints version and exits without error.
func TestMain_DoctypeVersionFlag(t *testing.T) {
	var buf bytes.Buffer
	cfg := runner.AppConfig{
		Name:    "doctype",
		Version: "test-version",
		Args:    []string{"-v"},
		Stdout:  &buf,
		SelectAndLoad: func(w, h int) (*selection.Selection, error) {
			return nil, nil
		},
		ListItems: func() ([]string, error) { return nil, nil },
	}
	if err := runner.RunApp(cfg); err != nil {
		t.Fatalf("RunApp returned error: %v", err)
	}
	out := buf.String()
	if out != "test-version\n" {
		t.Fatalf("expected version output, got %q", out)
	}
}

// TestMain_DoctypeListFlag ensures listing works.
func TestMain_DoctypeListFlag(t *testing.T) {
	var buf bytes.Buffer
	cfg := runner.AppConfig{
		Name:    "doctype",
		Version: "ignored",
		Args:    []string{"-l"},
		Stdout:  &buf,
		ListItems: func() ([]string, error) {
			return []string{"one", "two"}, nil
		},
		SelectAndLoad: func(w, h int) (*selection.Selection, error) { return nil, nil },
	}
	if err := runner.RunApp(cfg); err != nil {
		t.Fatalf("RunApp list returned error: %v", err)
	}
	out := buf.String()
	if out != "one\ntwo\n" {
		t.Fatalf("unexpected list output: %q", out)
	}
}
