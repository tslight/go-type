package selection

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/internal/menu"
)

// Selection represents the content and persistence hooks for a typing session.
type Selection struct {
	Text     string
	Content  *content.Content
	Provider StateProvider
}

// StateProvider abstracts persistence for typing sessions
type StateProvider interface {
	GetSavedCharPos() int
	GetSavedInput() string
	SaveProgress(charPos int, lastInput string) error
	RecordSession(wpm, accuracy float64, errors, charTypedRaw, effectiveChars, duration int) (string, error)
	// ResetState clears all progress and stats for this content.
	ResetState() error
}

// contentStateProvider implements StateProvider using a ContentManager (kept here to avoid package cycles)
type contentStateProvider struct {
	manager    *content.ContentManager
	contentID  string
	textLength int
	statsTitle string
}

func newContentStateProvider(manager *content.ContentManager, contentID string, textLength int, statsTitle string) *contentStateProvider {
	return &contentStateProvider{manager: manager, contentID: contentID, textLength: textLength, statsTitle: statsTitle}
}

func (p *contentStateProvider) GetSavedCharPos() int {
	return p.manager.StateManager.GetCharPos(p.contentID)
}
func (p *contentStateProvider) GetSavedInput() string {
	return p.manager.StateManager.GetLastInput(p.contentID)
}
func (p *contentStateProvider) SaveProgress(charPos int, lastInput string) error {
	name := p.contentID
	if current := p.manager.GetCurrentContent(); current != nil {
		name = current.Name
	}
	return p.manager.StateManager.SaveProgress(p.contentID, name, charPos, p.textLength, "", lastInput)
}
func (p *contentStateProvider) RecordSession(wpm, accuracy float64, errors, charTypedRaw, effectiveChars, duration int) (string, error) {
	name := p.contentID
	if current := p.manager.GetCurrentContent(); current != nil {
		name = current.Name
	}
	if err := p.manager.StateManager.RecordSession(p.contentID, name, wpm, accuracy, errors, charTypedRaw, effectiveChars, duration); err != nil {
		return "", err
	}
	stats := p.manager.StateManager.GetStats(p.contentID)
	return p.manager.StateManager.FormatStats(stats, p.statsTitle), nil
}

// SetFlash implements a lightweight interface for runner to set a flash message on the manager.
func (p *contentStateProvider) SetFlash(msg string) {
	if p.manager != nil {
		p.manager.SetPendingFlash(msg)
	}
}

// ResetState clears all progress and stats for this content.
func (p *contentStateProvider) ResetState() error {
	if p.manager == nil {
		return nil
	}
	// Remove all state (progress and stats) for this contentID.
	return p.manager.StateManager.ClearState(p.contentID)
}

// SelectContent runs the interactive menu, loads chosen content, and builds a Selection.
// (nil, nil) is returned if the user aborts.
func SelectContent(manager *content.ContentManager, width, height int) (*Selection, error) {
	if manager == nil {
		return nil, nil
	}
	menuModel := menu.NewMenuModel(manager, width, height)
	if _, err := runMenuProgram(menuModel); err != nil {
		return nil, err
	}
	selected := menuModel.SelectedContent()
	if selected == nil {
		return nil, nil
	}

	// Load based on manager mode.
	if manager.IsManifestBased() {
		if err := manager.SetContent(selected.ID); err != nil {
			return nil, err
		}
	} else {
		if err := manager.SetContentByName(selected.Name); err != nil {
			return nil, err
		}
	}

	text := manager.GetCurrentText()

	// Keep terminology content-agnostic throughout the CLI/package layer.
	provider := newContentStateProvider(manager, manager.StateKeyFor(*selected), len(text), "CONTENT STATISTICS")
	return &Selection{Text: text, Content: manager.GetCurrentContent(), Provider: provider}, nil
}

// runMenuProgram is a hook to run the Bubble Tea program.
// It’s a var so tests can stub it and avoid interactive UI.
var runMenuProgram = func(m tea.Model) (tea.Model, error) {
	return tea.NewProgram(m).Run()
}
