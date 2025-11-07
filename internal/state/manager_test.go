package state

import "testing"

func TestManager_AllStates_EmptyAndPopulated(t *testing.T) {
    t.Setenv("HOME", t.TempDir())
    m := NewManager[string, ContentState]("app", func(s *ContentState) (string, bool) {
        if s.ID == "" { return "", false }
        return s.ID, true
    }, nil)
    if len(m.AllStates()) != 0 {
        t.Fatalf("expected empty states initially")
    }
    if err := m.Set(ContentState{ID: "a", Name: "A"}); err != nil { t.Fatalf("set: %v", err) }
    if err := m.Set(ContentState{ID: "b", Name: "B"}); err != nil { t.Fatalf("set: %v", err) }
    all := m.AllStates()
    if len(all) != 2 { t.Fatalf("expected 2 states, got %d", len(all)) }
}

func TestManager_Configure_ResetAndError(t *testing.T) {
    t.Setenv("HOME", t.TempDir())
    m := NewManager[string, ContentState]("app", func(s *ContentState) (string, bool) {
        if s.ID == "" { return "", false }
        return s.ID, true
    }, nil)
    _ = m.Set(ContentState{ID: "x", Name: "X"})
    if err := m.Configure("new-app"); err != nil { t.Fatalf("configure: %v", err) }
    if len(m.AllStates()) != 0 { t.Fatalf("expected states to reset on configure") }
    if m.StateFilePath() == "" { t.Fatalf("expected state file path after configure") }

    // Configure with empty app name on manager with empty default should error
    m2 := NewManager[string, ContentState]("", func(s *ContentState) (string, bool) { return s.ID, s.ID != "" }, nil)
    if err := m2.Configure(""); err == nil { t.Fatalf("expected error configuring with empty app name") }
}
