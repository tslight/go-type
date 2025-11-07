package state

import (
	"strings"
	"testing"
)

func TestContentState_SaveProgressAndGet(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-1"
	name := "First"
	if err := mgr.SaveProgress(id, name, 50, 200, ""); err != nil {
		t.Fatalf("SaveProgress failed: %v", err)
	}
	st := mgr.GetState(id)
	if st == nil {
		t.Fatalf("expected state to be saved and retrievable")
	}
	if st.PercentComplete <= 0.0 || st.PercentComplete >= 100.0 {
		t.Fatalf("unexpected percent complete: %v", st.PercentComplete)
	}
	if got := mgr.GetCharPos(id); got != 50 {
		t.Fatalf("GetCharPos mismatch: got %d", got)
	}
}

func TestContentState_RecordSessionAndStats(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-2"
	name := "Second"
	// Record a couple of sessions
	if err := mgr.RecordSession(id, name, 50.0, 95.0, 2, 300, 60); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	if err := mgr.RecordSession(id, name, 70.0, 97.0, 1, 500, 120); err != nil {
		t.Fatalf("RecordSession failed: %v", err)
	}
	stats := mgr.GetStats(id)
	if stats["sessions_completed"].(int) != 2 {
		t.Fatalf("expected 2 sessions, got %v", stats["sessions_completed"])
	}
	if stats["total_time"].(int) != 180 {
		t.Fatalf("expected total time 180, got %v", stats["total_time"])
	}
	// Ensure formatting is non-empty and contains the title
	formatted := mgr.FormatStats(stats, "TEST STATS")
	if formatted == "" || !strings.Contains(formatted, "TEST STATS") {
		t.Fatalf("unexpected formatted stats: %q", formatted)
	}
}

func TestContentState_ClearState(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	id := "item-3"
	if err := mgr.SaveProgress(id, "Name", 10, 100, ""); err != nil {
		t.Fatalf("SaveProgress failed: %v", err)
	}
	if err := mgr.ClearState(id); err != nil {
		t.Fatalf("ClearState failed: %v", err)
	}
	if st := mgr.GetState(id); st != nil {
		t.Fatalf("expected state to be cleared")
	}
}
