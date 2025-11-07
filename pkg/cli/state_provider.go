package cli

import (
	"fmt"

	"github.com/tobe/go-type/internal/godocgen"
	"github.com/tobe/go-type/internal/textgen"
)

// StateProvider abstracts persistence for typing sessions
type StateProvider interface {
	GetSavedCharPos() int
	SaveProgress(charPos int) error
	RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error)
}

// TextgenStateProvider implements StateProvider for book-based typing
type TextgenStateProvider struct{}

// NewTextgenStateProvider creates a state provider backed by textgen
func NewTextgenStateProvider() *TextgenStateProvider {
	return &TextgenStateProvider{}
}

func (p *TextgenStateProvider) GetSavedCharPos() int {
	return textgen.GetCurrentCharPos()
}

func (p *TextgenStateProvider) SaveProgress(charPos int) error {
	return textgen.SaveProgress(charPos, "")
}

func (p *TextgenStateProvider) RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error) {
	if err := textgen.RecordSession(wpm, accuracy, errors, charTyped, duration); err != nil {
		return "", err
	}
	stats := textgen.GetCurrentBookStats()
	if stats == nil {
		return "", nil
	}
	return textgen.FormatBookStats(stats), nil
}

// DocStateProvider implements StateProvider for Go documentation typing
type DocStateProvider struct {
	docName    string
	textLength int
}

// NewDocStateProvider creates a new provider for a documentation module
func NewDocStateProvider(docName string, textLength int) (*DocStateProvider, error) {
	if docName == "" {
		return nil, fmt.Errorf("doc name cannot be empty")
	}
	return &DocStateProvider{docName: docName, textLength: textLength}, nil
}

func (p *DocStateProvider) GetSavedCharPos() int {
	return godocgen.GetSavedCharPos(p.docName)
}

func (p *DocStateProvider) SaveProgress(charPos int) error {
	return godocgen.SaveDocProgress(p.docName, charPos, p.textLength)
}

func (p *DocStateProvider) RecordSession(wpm, accuracy float64, errors, charTyped, duration int) (string, error) {
	if err := godocgen.RecordDocSession(p.docName, wpm, accuracy, errors, charTyped, duration); err != nil {
		return "", err
	}
	stats := godocgen.GetDocStats(p.docName)
	return godocgen.FormatDocStats(stats), nil
}
