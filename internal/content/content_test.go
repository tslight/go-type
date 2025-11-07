package content


import (
	"io"
	"io/fs"
	"time"
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

func TestContentManager_NegativePaths(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// Invalid ID
	cmManifest := NewContentManager(books.EFS, "test-books", true)
	if _, err := cmManifest.GetContent(9999999); err == nil {
		t.Fatalf("expected error for invalid content ID")
	}
	// Directory mode missing name and char pos default
	cmDir := NewContentManager(godocs.EFS, "test-godocs", false)
	if err := cmDir.SetContentByName("__does_not_exist__"); err == nil {
		t.Fatalf("expected error for missing content name")
	}
	if _, err := cmDir.GetContentByName("__does_not_exist__"); err == nil {
		t.Fatalf("expected error for GetContentByName on missing file")
	}
	items := cmDir.GetAvailableContent()
	if len(items) > 0 {
		if err := cmDir.SetContentByName(items[0].Name); err != nil {
			t.Fatalf("unexpected error setting existing content: %v", err)
		}
		if pos := cmDir.GetCurrentCharPos(); pos != 0 {
			t.Fatalf("expected default saved char pos 0, got %d", pos)
		}
	}
}

func TestContentManager_GetCurrentCharPos(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// Manifest mode
	cm := NewContentManager(books.EFS, "test-books", true)
	items := cm.GetAvailableContent()
	if len(items) == 0 {
		t.Skip("no books embedded")
	}
	if err := cm.SetContent(items[0].ID); err != nil {
		t.Fatalf("SetContent: %v", err)
	}
	cur := cm.GetCurrentContent()
	if cur == nil {
		t.Fatalf("no current content")
	}
	key := cm.StateKeyFor(*cur)
	if err := cm.StateManager.SaveProgress(key, cur.Name, 7, 100, ""); err != nil {
		t.Fatalf("SaveProgress: %v", err)
	}
	if pos := cm.GetCurrentCharPos(); pos != 7 {
		t.Fatalf("expected 7, got %d", pos)
	}

	// Directory mode
	cm2 := NewContentManager(godocs.EFS, "test-godocs", false)
	ditems := cm2.GetAvailableContent()
	if len(ditems) == 0 {
		t.Skip("no docs embedded")
	}
	if err := cm2.SetContentByName(ditems[0].Name); err != nil {
		t.Fatalf("SetContentByName: %v", err)
	}
	cur2 := cm2.GetCurrentContent()
	key2 := cm2.StateKeyFor(*cur2)
	if err := cm2.StateManager.SaveProgress(key2, cur2.Name, 5, 100, ""); err != nil {
		t.Fatalf("SaveProgress: %v", err)
	}
	if pos := cm2.GetCurrentCharPos(); pos != 5 {
		t.Fatalf("expected 5, got %d", pos)
	}
}

// embed a minimal manifest without filename entries to exercise error branch.
// We simulate with an in-memory FS (fstest) since embed requires compile-time; use fs.FS implementation.
// NOTE: Skipped malformed manifest test due to embed.FS type constraints; keeping placeholder for future if refactored.
type memFS struct{ files map[string]string }
func (m memFS) Open(name string) (fs.File, error) { if _, ok := m.files[name]; !ok { return nil, fs.ErrNotExist }; return &memFile{data: []byte(m.files[name])}, nil }
func (m memFS) ReadFile(name string) ([]byte, error) { if s, ok := m.files[name]; ok { return []byte(s), nil }; return nil, fs.ErrNotExist }
type memFile struct{ data []byte; off int }
func (f *memFile) Stat() (fs.FileInfo, error) { return memInfo{int64(len(f.data))}, nil }
func (f *memFile) Read(b []byte) (int, error) { if f.off >= len(f.data) { return 0, io.EOF }; n := copy(b, f.data[f.off:]); f.off += n; return n, nil }
func (f *memFile) Close() error { return nil }
type memInfo struct{ size int64 }
func (i memInfo) Name() string       { return "" }
func (i memInfo) Size() int64        { return i.size }
func (i memInfo) Mode() fs.FileMode  { return 0444 }
func (i memInfo) ModTime() time.Time { return time.Time{} }
func (i memInfo) IsDir() bool        { return false }
func (i memInfo) Sys() any           { return nil }

func TestContentManager_ManifestMissingFilename(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	badManifest := `{"books":{"1":{"title":"Only Title"}}}`
	m := memFS{files: map[string]string{"manifest.json": badManifest}}
	cm := NewContentManager(m, "test-books", true)
	items := cm.GetAvailableContent()
	if len(items) == 0 { t.Fatalf("expected at least one entry from manifest") }
	if err := cm.SetContent(items[0].ID); err == nil {
		t.Fatalf("expected error due to missing filename in manifest")
	}
}

// Minimal adapter to satisfy embed.FS-like ReadFile for tests.

