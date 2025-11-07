package main

import (
"fmt"
"os"

tea "github.com/charmbracelet/bubbletea"
"github.com/tobe/go-type/assets/godocs"
"github.com/tobe/go-type/internal/content"
"github.com/tobe/go-type/pkg/cli"
)

var Version = "unknown"

func main() {
	manager := content.NewContentManager(godocs.EFS, "doctype", false)
	
	config := cli.AppConfig{
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
		Configure:     []func() error{},
		SelectAndLoad: func(width, height int) (*cli.Selection, error) {
			return selectDoc(manager, width, height)
		},
	}

	if err := cli.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func selectDoc(manager *content.ContentManager, width, height int) (*cli.Selection, error) {
	menuModel := cli.NewDocMenuModel(manager, width, height)
	program := tea.NewProgram(menuModel)

	if _, err := program.Run(); err != nil {
		return nil, err
	}

	namePtr := menuModel.SelectedDocName()
	if namePtr == nil {
		return nil, nil
	}
	selectedDocName := *namePtr

	if err := manager.SetContentByName(selectedDocName); err != nil {
		return nil, err
	}

	text := manager.GetCurrentText()
	provider := cli.NewDocStateProvider(manager, selectedDocName, len(text))

	selection := &cli.Selection{
		Text:     text,
		Book:     &content.Content{ID: 0, Name: selectedDocName, Text: text},
		Provider: provider,
	}

	return selection, nil
}
