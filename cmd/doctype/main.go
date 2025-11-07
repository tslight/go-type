package main

import (
	"fmt"
	"os"

	"github.com/tobe/go-type/assets/godocs"
	"github.com/tobe/go-type/internal/content"
	"github.com/tobe/go-type/internal/runner"
	"github.com/tobe/go-type/internal/selection"
)

var Version = "unknown"

func main() {
	manager := content.NewContentManager(godocs.EFS, "doctype", false)

	config := runner.AppConfig{
		Name:            "doctype",
		Version:         Version,
		ListDescription: "List available Go documentation modules",
		ListItems: func() ([]string, error) {
			contents := manager.GetAvailableContent()
			names := make([]string, len(contents))
			for i, c := range contents {
				names[i] = c.Name
			}
			return names, nil
		},
		Configure: []func() error{},
		SelectAndLoad: func(width, height int) (*selection.Selection, error) {
			return selection.SelectContent(manager, width, height)
		},
	}

	if err := runner.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
