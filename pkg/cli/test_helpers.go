package cli

import (
"github.com/tobe/go-type/assets/books"
"github.com/tobe/go-type/internal/content"
)

// NewTestContentManager creates a ContentManager with test data
func NewTestContentManager() *content.ContentManager {
	return content.NewContentManager(books.EFS, "test-gutentype", true)
}
