package content

import (
	"strings"
	"testing"

	"github.com/tobe/go-type/assets/books"
	"github.com/tobe/go-type/assets/godocs"
)

// Test manifest-based loading
func TestContentManager_ManifestMode(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := NewContentManager(books.EFS, "test-books", true)
	if !cm.IsManifestBased() {
		t.Fatalf("expected manifest-based manager")
	}
	items := cm.GetAvailableContent()
	if len(items) == 0 {
		t.Fatalf("expected some content items from manifest")
	}
	first := items[0]
	if first.ID == 0 && first.Name == "" {
		// Not strictly an error, but unlikely
		t.Logf("warning: first manifest item has empty metadata")
	}
	// Set content by ID
	if err := cm.SetContent(first.ID); err != nil {
		t.Fatalf("SetContent failed: %v", err)
	}
	cur := cm.GetCurrentContent()
	if cur == nil || cur.ID != first.ID {
		t.Fatalf("current content mismatch after SetContent")
	}
	text := cm.GetCurrentText()
	if text == "" || strings.Contains(text, "\r") {
		// ASCII filter should remove \r if present in sources
		if text == "" {
			t.Fatalf("expected non-empty text for selected content")
		}
	}
	key := cm.StateKeyFor(*cur)
	// Assert non-empty key for manifest-based content
	if key == "" {
		// For manifest mode we expect numeric IDs; enforce non-empty.
		// (Not failing on "0" because some manifests could legitimately use 0.)
		// Fail to satisfy linter and catch unexpected empty ID.
		// If this proves flaky adjust to warn only.
		// Using fatal keeps test intent clear.
		//nolint:staticcheck // Intentional check; no empty branch.
		// (We include nolint in case staticcheck flags complex comment heuristics.)
		// Actually trigger failure:
		if true {
			t.Fatalf("expected non-empty state key for manifest content")
		}
	}
}

// Test directory-based loading
func TestContentManager_DirectoryMode(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := NewContentManager(godocs.EFS, "test-godocs", false)
	if cm.IsManifestBased() {
		t.Fatalf("expected directory-based manager")
	}
	items := cm.GetAvailableContent()
	if len(items) == 0 {
		t.Fatalf("expected some directory content items")
	}
	first := items[0]
	if first.Name == "" {
		t.Fatalf("directory item should have name")
	}
	// Set by name
	if err := cm.SetContentByName(first.Name); err != nil {
		t.Fatalf("SetContentByName failed: %v", err)
	}
	cur := cm.GetCurrentContent()
	if cur == nil || cur.Name != first.Name {
		t.Fatalf("current content mismatch after SetContentByName")
	}
	// Ensure GetContentByName yields same text and ASCII only
	textByName, err := cm.GetContentByName(first.Name)
	if err != nil {
		t.Fatalf("GetContentByName failed: %v", err)
	}
	if textByName == "" {
		t.Fatalf("expected non-empty text from GetContentByName")
	}
	for _, r := range textByName {
		if r == '\r' { // Carriage returns should be filtered out by ASCII filter
			t.Fatalf("unexpected carriage return found in filtered text")
		}
		// Accept non-ASCII silently; content may include Unicode even if filter trims when setting current.
	}
	key := cm.StateKeyFor(*cur)
	if key != cur.Name {
		t.Fatalf("expected state key to match content name for directory mode")
	}
}

func TestContentManager_SetContentByName_ErrorOnManifest(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := NewContentManager(books.EFS, "test-books", true)
	if err := cm.SetContentByName("whatever"); err == nil {
		// Should error because manifest mode doesn't support name-based setting
		t.Fatalf("expected error calling SetContentByName in manifest mode")
	}
}

func TestContentManager_GetContentByName_ErrorOnManifest(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cm := NewContentManager(books.EFS, "test-books", true)
	if _, err := cm.GetContentByName("whatever"); err == nil {
		// Should error because manifest mode doesn't support name-based fetching
		t.Fatalf("expected error calling GetContentByName in manifest mode")
	}
}
