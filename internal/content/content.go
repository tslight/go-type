package content

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tobe/go-type/internal/state"
)

// Content represents a piece of content (book, doc, etc.) for typing practice
type Content struct {
	ID   int
	Name string
	Text string
}

// ContentManager manages loading and state for any type of content
type ContentManager struct {
	fs               ReadableFS
	StateManager     *state.ContentStateManager
	availableContent []Content
	currentContent   *Content
	rng              *rand.Rand
	useManifest      bool // Whether to use manifest.json or directory listing
	lastSearchQuery  string
	lastSearchDir    int    // 1 for forward, -1 for backward
	pendingFlash     string // transient flash message consumed by next menu
}

// ReadableFS is the minimal filesystem interface ContentManager needs.
// It must support directory operations (fs.FS) and direct file reads (ReadFile).
type ReadableFS interface {
	fs.FS
	ReadFile(name string) ([]byte, error)
}

// NewContentManager creates a new content manager for the given embedded filesystem
// name is used for state file naming (e.g., "gutentype", "doctype")
// useManifest determines whether to load from manifest.json (true) or directory listing (false)
func NewContentManager(fileSystem ReadableFS, name string, useManifest bool) *ContentManager {
	cm := &ContentManager{
		fs:            fileSystem,
		StateManager:  state.NewContentStateManager(name),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		useManifest:   useManifest,
		lastSearchDir: 1,
	}

	cm.loadAvailableContent()

	// Set a random content as default
	if len(cm.availableContent) > 0 {
		randomContent := cm.availableContent[cm.rng.Intn(len(cm.availableContent))]
		_ = cm.SetContent(randomContent.ID)
	}

	return cm
}

// SetLastSearch stores the last search query and direction so the menu can restore it on reopen.
func (cm *ContentManager) SetLastSearch(query string, direction int) {
	cm.lastSearchQuery = query
	if direction != 1 && direction != -1 {
		direction = 1
	}
	cm.lastSearchDir = direction
}

// GetLastSearch returns the last search query and direction.
func (cm *ContentManager) GetLastSearch() (string, int) { return cm.lastSearchQuery, cm.lastSearchDir }

// SetPendingFlash sets a transient message to be shown by the next menu model instance.
func (cm *ContentManager) SetPendingFlash(msg string) { cm.pendingFlash = msg }

// ConsumePendingFlash returns the pending flash and clears it.
func (cm *ContentManager) ConsumePendingFlash() string {
	msg := cm.pendingFlash
	cm.pendingFlash = ""
	return msg
}

// IsManifestBased reports whether this manager loads content via a manifest.json
// (true for books) vs directory listing (false for docs).
func (cm *ContentManager) IsManifestBased() bool {
	return cm.useManifest
}

// StateKeyFor returns the state persistence key for a given content item.
// For manifest-based managers this is the numeric ID; for directory-based it is the name.
func (cm *ContentManager) StateKeyFor(c Content) string {
	if cm.useManifest {
		return strconv.Itoa(c.ID)
	}
	return c.Name
}

// loadAvailableContent loads the list of available content from the embedded filesystem
func (cm *ContentManager) loadAvailableContent() {
	if cm.useManifest {
		cm.loadFromManifest()
	} else {
		cm.loadFromDirectory()
	}

	// Sort alphabetically by name
	sort.Slice(cm.availableContent, func(i, j int) bool {
		return cm.availableContent[i].Name < cm.availableContent[j].Name
	})
}

// loadFromManifest loads content from a manifest.json file
func (cm *ContentManager) loadFromManifest() {
	manifestData, err := cm.fs.ReadFile("manifest.json")
	if err != nil {
		return
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return
	}

	// Support both "books" and "content" keys for flexibility
	var contentMap map[string]interface{}
	if booksMap, ok := manifest["books"].(map[string]interface{}); ok {
		contentMap = booksMap
	} else if items, ok := manifest["content"].(map[string]interface{}); ok {
		contentMap = items
	} else {
		return
	}

	for idStr, itemData := range contentMap {
		id := 0
		_, _ = fmt.Sscanf(idStr, "%d", &id)
		if id <= 0 {
			continue
		}

		content := Content{ID: id, Name: "Unknown"}

		if itemInfo, ok := itemData.(map[string]interface{}); ok {
			if title, ok := itemInfo["title"].(string); ok {
				content.Name = title
			}
		}

		cm.availableContent = append(cm.availableContent, content)
	}
}

// loadFromDirectory loads content from .txt files in the root directory
func (cm *ContentManager) loadFromDirectory() {
	entries, err := fs.ReadDir(cm.fs, ".")
	if err != nil {
		return
	}

	id := 0
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			// Remove .txt extension and convert hyphens to slashes for display
			name := strings.TrimSuffix(entry.Name(), ".txt")
			name = strings.ReplaceAll(name, "-", "/")

			cm.availableContent = append(cm.availableContent, Content{
				ID:   id,
				Name: name,
			})
			id++
		}
	}
}

// GetAvailableContent returns the list of available content
func (cm *ContentManager) GetAvailableContent() []Content {
	return cm.availableContent
}

// GetContent loads and returns the full text for a specific content item
func (cm *ContentManager) GetContent(id int) (string, error) {
	// Find the content item
	var content *Content
	for i := range cm.availableContent {
		if cm.availableContent[i].ID == id {
			content = &cm.availableContent[i]
			break
		}
	}

	if content == nil {
		return "", fmt.Errorf("content with ID %d not found", id)
	}

	return cm.loadContentText(content)
}

// GetContentByName loads content by name (for directory-based loading)
func (cm *ContentManager) GetContentByName(name string) (string, error) {
	if cm.useManifest {
		return "", fmt.Errorf("GetContentByName not supported for manifest-based content")
	}

	// Convert slashes to hyphens for filename
	filename := strings.ReplaceAll(name, "/", "-") + ".txt"

	data, err := cm.fs.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read content %s: %w", name, err)
	}

	return cm.filterToASCII(string(data)), nil
}

// loadContentText loads the text for a content item
func (cm *ContentManager) loadContentText(content *Content) (string, error) {
	var filename string

	if cm.useManifest {
		// Load from manifest to get filename
		manifestData, err := cm.fs.ReadFile("manifest.json")
		if err != nil {
			return "", fmt.Errorf("failed to read manifest: %w", err)
		}

		var manifest map[string]interface{}
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			return "", fmt.Errorf("failed to parse manifest: %w", err)
		}

		var contentMap map[string]interface{}
		if booksMap, ok := manifest["books"].(map[string]interface{}); ok {
			contentMap = booksMap
		} else if items, ok := manifest["content"].(map[string]interface{}); ok {
			contentMap = items
		}

		idStr := fmt.Sprintf("%d", content.ID)
		if itemInfo, ok := contentMap[idStr].(map[string]interface{}); ok {
			if fname, ok := itemInfo["filename"].(string); ok {
				filename = fname
			}
		}

		if filename == "" {
			return "", fmt.Errorf("filename not found in manifest for ID %d", content.ID)
		}
	} else {
		// Use name directly (convert slashes to hyphens)
		filename = strings.ReplaceAll(content.Name, "/", "-") + ".txt"
	}

	data, err := cm.fs.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return cm.filterToASCII(string(data)), nil
}

// SetContent sets the current content and loads it
func (cm *ContentManager) SetContent(id int) error {
	text, err := cm.GetContent(id)
	if err != nil {
		return err
	}

	// Find the content item
	for i := range cm.availableContent {
		if cm.availableContent[i].ID == id {
			cm.currentContent = &Content{
				ID:   cm.availableContent[i].ID,
				Name: cm.availableContent[i].Name,
				Text: text,
			}
			return nil
		}
	}

	return fmt.Errorf("content with ID %d not found", id)
}

// SetContentByName sets content by name (for directory-based content)
func (cm *ContentManager) SetContentByName(name string) error {
	if cm.useManifest {
		return fmt.Errorf("SetContentByName not supported for manifest-based content")
	}

	// Find the content by name
	for i := range cm.availableContent {
		if cm.availableContent[i].Name == name {
			text, err := cm.GetContentByName(name)
			if err != nil {
				return err
			}

			cm.currentContent = &Content{
				ID:   cm.availableContent[i].ID,
				Name: name,
				Text: text,
			}
			return nil
		}
	}

	return fmt.Errorf("content %s not found", name)
}

// GetCurrentContent returns the currently loaded content
func (cm *ContentManager) GetCurrentContent() *Content {
	return cm.currentContent
}

// GetCurrentText returns the text of the currently loaded content
func (cm *ContentManager) GetCurrentText() string {
	if cm.currentContent == nil {
		return "No content loaded"
	}
	return cm.currentContent.Text
}

// GetCurrentCharPos returns the saved character position for the current content
func (cm *ContentManager) GetCurrentCharPos() int {
	if cm.currentContent == nil {
		return 0
	}

	id := cm.getContentID()
	return cm.StateManager.GetCharPos(id)
}

// getContentID returns a string ID for the current content
func (cm *ContentManager) getContentID() string {
	if cm.currentContent == nil {
		return ""
	}

	if cm.useManifest {
		return strconv.Itoa(cm.currentContent.ID)
	}
	// For directory-based content, use the name as ID
	return cm.currentContent.Name
}

// filterToASCII filters text to ASCII only to avoid UTF-8 encoding issues
// Preserves newlines for paragraph formatting
func (cm *ContentManager) filterToASCII(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		// Keep newlines and ASCII characters
		if s[i] == '\n' || (s[i] < 128 && s[i] >= 32) || s[i] == '\t' {
			result = append(result, s[i])
		}
	}
	return string(result)
}
