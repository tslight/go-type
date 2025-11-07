package cli

import (
	"github.com/tobe/go-type/internal/content"
)

// StateProvider abstracts persistence for typing sessions
type StateProvider interface {
	GetSavedCharPos() int
	SaveProgress(charPos int) error
	RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error)
}

// ContentStateProvider implements StateProvider using a ContentManager
type ContentStateProvider struct {
	manager    *content.ContentManager
	contentID  string
	textLength int
	statsTitle string
}

// NewContentStateProvider creates a state provider for a content manager and specific content
func NewContentStateProvider(manager *content.ContentManager, contentID string, textLength int, statsTitle string) *ContentStateProvider {
	return &ContentStateProvider{
		manager:    manager,
		contentID:  contentID,
		textLength: textLength,
		statsTitle: statsTitle,
	}
}

func (p *ContentStateProvider) GetSavedCharPos() int {
	return p.manager.StateManager.GetCharPos(p.contentID)
}

func (p *ContentStateProvider) SaveProgress(charPos int) error {
	// Get content name for display
	contentName := p.contentID
	if current := p.manager.GetCurrentContent(); current != nil {
		contentName = current.Name
	}
	return p.manager.StateManager.SaveProgress(p.contentID, contentName, charPos, p.textLength, "")
}

func (p *ContentStateProvider) RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error) {
	// Get content name for display
	contentName := p.contentID
	if current := p.manager.GetCurrentContent(); current != nil {
		contentName = current.Name
	}

	if err := p.manager.StateManager.RecordSession(p.contentID, contentName, wpm, accuracy, errors, charTyped, duration); err != nil {
		return "", err
	}
	stats := p.manager.StateManager.GetStats(p.contentID)
	return p.manager.StateManager.FormatStats(stats, p.statsTitle), nil
}

// Note: Backward-compat helper constructors have been removed in favor of the
// single content-agnostic NewContentStateProvider.
