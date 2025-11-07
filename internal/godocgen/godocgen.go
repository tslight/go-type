package godocgen

import (
	"fmt"
	"io/fs"
	"math/rand"
	"strings"

	"github.com/tobe/go-type/assets/godocs"
)

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
