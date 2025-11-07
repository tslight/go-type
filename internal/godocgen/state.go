package godocgen

import (
	"github.com/tobe/go-type/internal/statestore"
)

var docManager = statestore.NewContentStateManager("doctype")

// ConfigureStateFile overrides the state file name based on the provided app name.
func ConfigureStateFile(appName string) error {
	return docManager.Configure(appName)
}

// GetDocState returns the state for a doc, if any
func GetDocState(docName string) *statestore.ContentState {
	return docManager.GetState(docName)
}

// GetSavedCharPos returns the saved character position for a doc
func GetSavedCharPos(docName string) int {
	return docManager.GetCharPos(docName)
}

// SaveDocProgress stores the current progress for a doc
func SaveDocProgress(docName string, charPos int, textLength int) error {
	return docManager.SaveProgress(docName, docName, charPos, textLength, "")
}

// RecordDocSession appends a session result for a doc
func RecordDocSession(docName string, wpm, accuracy float64, errors, charTyped, duration int) error {
	return docManager.RecordSession(docName, docName, wpm, accuracy, errors, charTyped, duration)
}

// GetDocStats returns aggregated statistics for a doc
func GetDocStats(docName string) map[string]interface{} {
	return docManager.GetStats(docName)
}

// FormatDocStats returns a formatted stats string for display
func FormatDocStats(stats map[string]interface{}) string {
	return docManager.FormatStats(stats, "DOCUMENT STATISTICS")
}
