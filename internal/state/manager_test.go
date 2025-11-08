package state

import "testing"

func TestManager_AllStates_EmptyAndPopulated(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := NewManager[string, ContentState]("app", func(s *ContentState) (string, bool) {
		if s.ID == "" {
			return "", false
		}
		return s.ID, true
	}, nil)
	if len(m.AllStates()) != 0 {
		t.Fatalf("expected empty states initially")
	}
	if err := m.Set(ContentState{ID: "a", Name: "A"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := m.Set(ContentState{ID: "b", Name: "B"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	all := m.AllStates()
	if len(all) != 2 {
		t.Fatalf("expected 2 states, got %d", len(all))
	}
}

func TestManager_Configure_ResetAndError(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := NewManager[string, ContentState]("app", func(s *ContentState) (string, bool) {
		if s.ID == "" {
			return "", false
		}
		return s.ID, true
	}, nil)
	_ = m.Set(ContentState{ID: "x", Name: "X"})
	if err := m.Configure("new-app"); err != nil {
		t.Fatalf("configure: %v", err)
	}
	if len(m.AllStates()) != 0 {
		t.Fatalf("expected states to reset on configure")
	}
	if m.StateFilePath() == "" {
		t.Fatalf("expected state file path after configure")
	}

	// Configure with empty app name on manager with empty default should error
	m2 := NewManager[string, ContentState]("", func(s *ContentState) (string, bool) { return s.ID, s.ID != "" }, nil)
	if err := m2.Configure(""); err == nil {
		t.Fatalf("expected error configuring with empty app name")
	}
}

func TestManager_ComputeStateFilePath_Fallbacks(t *testing.T) {
	// Unset HOME to force fallback; t.Setenv with empty string
	t.Setenv("HOME", "")
	m := NewManager[string, ContentState]("app", func(s *ContentState) (string, bool) { return s.ID, s.ID != "" }, nil)
	path := m.StateFilePath()
	if path == "" {
		t.Fatalf("expected non-empty path even with empty HOME")
	}
	// Should include .app.json suffix
	if suffix := ".app.json"; len(path) < len(suffix) || path[len(path)-len(suffix):] != suffix {
		t.Fatalf("expected path to end with %s, got %s", suffix, path)
	}
}

func TestManager_WipeFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	m := NewManager("wipeapp", func(s *ContentState) (string, bool) { return s.ID, s.ID != "" }, nil)
	if err := m.Set(ContentState{ID: "x", Name: "X"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if len(m.AllStates()) != 1 {
		t.Fatalf("expected 1 state before wipe")
	}
	if err := m.WipeFile(); err != nil {
		t.Fatalf("wipe: %v", err)
	}
	if len(m.AllStates()) != 0 {
		t.Fatalf("expected states cleared after wipe")
	}
	// Second wipe should succeed silently even if file absent
	if err := m.WipeFile(); err != nil {
		t.Fatalf("second wipe should not error: %v", err)
	}
}

func TestBuildStateFileName_Error(t *testing.T) {
	if _, err := BuildStateFileName(" "); err == nil {
		t.Fatalf("expected error for empty/whitespace app name")
	}
	name, err := BuildStateFileName("MyApp")
	if err != nil || name != ".myapp.json" {
		t.Fatalf("unexpected name result: %s err=%v", name, err)
	}
}
