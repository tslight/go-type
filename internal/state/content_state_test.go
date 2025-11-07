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

func TestContentState_EmptyIDErrors(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr := NewContentStateManager("test-app-state")
	if err := mgr.SaveProgress("", "", 0, 0, ""); err == nil {
		t.Fatalf("expected error for empty ID in SaveProgress")
	}
	if err := mgr.RecordSession("", "", 0, 0, 0, 0, 0); err == nil {
		t.Fatalf("expected error for empty ID in RecordSession")
	}
}

func TestBuildStateFileName_Variants(t *testing.T) {
	name, err := BuildStateFileName("App")
	if err != nil || name == "" {
		t.Fatalf("unexpected error or empty name: %v %q", err, name)
	}
	if _, err := BuildStateFileName(""); err == nil {
		t.Fatalf("expected error for empty app name")
	}
}

func TestManager_ConfigureAndAllStates(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// Use the underlying generic manager through ContentStateManager
	mgr := NewContentStateManager("app")
	// Configure with new app name
	if err := mgr.Configure("another"); err != nil {
		t.Fatalf("configure: %v", err)
	}
	// Save two states
	_ = mgr.SaveProgress("id1", "n1", 1, 10, "")
	_ = mgr.SaveProgress("id2", "n2", 2, 10, "")
	// Access the underlying AllStates by retrieving and counting
	// (We can't call AllStates directly; assert through Get on both IDs)
	if mgr.GetState("id1") == nil || mgr.GetState("id2") == nil {
		t.Fatalf("expected both states present")
	}
	// StateFilePath is non-empty after configure since we wrote
	if mgr.store.StateFilePath() == "" {
		t.Fatalf("expected non-empty state file path")
	}
}

func TestManager_LoadStatesOnNewInstance(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	mgr1 := NewContentStateManager("app")
	_ = mgr1.SaveProgress("idA", "NameA", 3, 10, "")
	_ = mgr1.SaveProgress("idB", "NameB", 4, 10, "")
	// New manager should load from same state file
	mgr2 := NewContentStateManager("app")
	if mgr2.GetState("idA") == nil || mgr2.GetState("idB") == nil {
		t.Fatalf("expected states to be loaded by new manager instance")
	}
}
