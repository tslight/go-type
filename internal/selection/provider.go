package selection

import "github.com/tobe/go-type/internal/content"

// NewContentStateProvider is exported for use in other packages (e.g., menu)
func NewContentStateProvider(manager *content.ContentManager, contentID string, textLength int, statsTitle string) StateProvider {
	return newContentStateProvider(manager, contentID, textLength, statsTitle)
}
