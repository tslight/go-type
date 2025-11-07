package godocgen

import (
	"fmt"
	"io/fs"
	"math/rand"
	"strings"

	"github.com/tobe/go-type/assets/godocs"
)

// Doc represents a Go documentation module for typing practice
type Doc struct {
	ID   int    // Unique identifier for the doc
	Name string // Display name (e.g., "net/http")
	Text string // Full documentation text
}

var currentDoc *Doc

// GetAvailableDocumentation returns a list of available Go documentation modules
func GetAvailableDocumentation() ([]string, error) {
	var docs []string
	entries, err := fs.ReadDir(godocs.EFS, ".")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			// Remove .txt extension and convert hyphens back to slashes for display
			name := strings.TrimSuffix(entry.Name(), ".txt")
			name = strings.ReplaceAll(name, "-", "/")
			docs = append(docs, name)
		}
	}

	return docs, nil
}

// GetDocumentation returns the full text of a specific Go documentation module
func GetDocumentation(name string) (string, error) {
	// Convert slashes to hyphens for filename
	filename := strings.ReplaceAll(name, "/", "-") + ".txt"

	data, err := fs.ReadFile(godocs.EFS, filename)
	if err != nil {
		return "", fmt.Errorf("failed to read documentation for %s: %w", name, err)
	}

	return string(data), nil
}

// GetRandomDocumentation returns a random Go documentation module
func GetRandomDocumentation() (string, error) {
	docs, err := GetAvailableDocumentation()
	if err != nil {
		return "", err
	}

	if len(docs) == 0 {
		return "", fmt.Errorf("no documentation available")
	}

	randomDoc := docs[rand.Intn(len(docs))]
	return GetDocumentation(randomDoc)
}

// GetDocumentationNames returns friendly names of all documentation
func GetDocumentationNames() []string {
	docs, err := GetAvailableDocumentation()
	if err != nil {
		return []string{}
	}
	return docs
}

// GetAvailableDocs returns all docs as Doc structs with metadata
func GetAvailableDocs() ([]Doc, error) {
	names, err := GetAvailableDocumentation()
	if err != nil {
		return nil, err
	}

	var docs []Doc
	for i, name := range names {
		text, err := GetDocumentation(name)
		if err != nil {
			continue // Skip docs that fail to load
		}
		docs = append(docs, Doc{
			ID:   i,
			Name: name,
			Text: text,
		})
	}

	return docs, nil
}

// SetDoc sets the current documentation and returns it
func SetDoc(docID int) (*Doc, error) {
	docs, err := GetAvailableDocs()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		if doc.ID == docID {
			currentDoc = &doc
			return &doc, nil
		}
	}

	return nil, fmt.Errorf("documentation with ID %d not found", docID)
}

// GetCurrentDoc returns the currently selected documentation
func GetCurrentDoc() *Doc {
	return currentDoc
}

// GetDocText returns the text for a documentation by ID
func GetDocText(docID int) (string, error) {
	docs, err := GetAvailableDocs()
	if err != nil {
		return "", err
	}

	for _, doc := range docs {
		if doc.ID == docID {
			return doc.Text, nil
		}
	}

	return "", fmt.Errorf("documentation with ID %d not found", docID)
}
