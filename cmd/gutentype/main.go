package main

import (
"fmt"
"os"

tea "github.com/charmbracelet/bubbletea"
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
		Configure:     []func() error{},
		SelectAndLoad: func(width, height int) (*cli.Selection, error) {
			return selectBook(manager, width, height)
		},
	}

	if err := cli.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func selectBook(manager *content.ContentManager, width, height int) (*cli.Selection, error) {
	menuModel := cli.NewMenuModel(manager, width, height)
	program := tea.NewProgram(menuModel)

	if _, err := program.Run(); err != nil {
		return nil, err
	}

	selected := menuModel.SelectedBook()
	if selected == nil {
		return nil, nil
	}

	if err := manager.SetContent(selected.ID); err != nil {
		return nil, err
	}

	text := manager.GetCurrentText()
	provider := cli.NewBookStateProvider(manager, selected.ID, len(text))

	return &cli.Selection{
		Text:     text,
		Book:     selected,
		Provider: provider,
	}, nil
}
