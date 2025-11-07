package main

import (
	"fmt"
	"os"

	"github.com/tobe/go-type/assets/books"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/pkg/cli"
)

var Version = "unknown"

func main() {
	manager := content.NewContentManager(books.EFS, "gutentype", true)

	config := cli.AppConfig{
		Name:            "gutentype",
		Version:         Version,
		ListDescription: "List available books and their titles",
		ListItems: func() ([]string, error) {
			contents := manager.GetAvailableContent()
			names := make([]string, 0, len(contents))
			for _, c := range contents {
				names = append(names, c.Name)
			}
			return names, nil
		},
		Configure: []func() error{},
		SelectAndLoad: func(width, height int) (*cli.Selection, error) {
			return cli.SelectContent(manager, width, height)
		},
	}

	if err := cli.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Removed duplicate selectBook; unified via cli.SelectContent.
