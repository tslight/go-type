package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
	"github.com/tobe/go-type/pkg/cli"
)

var Version = "unknown"

func main() {
	config := cli.AppConfig{
		Name:            "gutentype",
		Version:         Version,
		ListDescription: "List available books and their titles",
		ListItems: func() ([]string, error) {
			books := textgen.GetAvailableBooks()
			names := make([]string, 0, len(books))
			for _, book := range books {
				names = append(names, book.Name)
			}
			return names, nil
		},
		Configure: []func() error{
			func() error { return textgen.ConfigureStateFile("gutentype") },
		},
		SelectAndLoad: selectBook,
	}

	if err := cli.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func selectBook(width, height int) (*cli.Selection, error) {
	menuModel := cli.NewMenuModel(width, height)
	program := tea.NewProgram(menuModel)

	if _, err := program.Run(); err != nil {
		return nil, err
	}

	selected := menuModel.SelectedBook()
	if selected == nil {
		return nil, nil
	}

	if err := textgen.SetBook(selected.ID); err != nil {
		return nil, err
	}

	text := textgen.GetFullText()
	provider := cli.NewTextgenStateProvider()

	return &cli.Selection{
		Text:     text,
		Book:     selected,
		Provider: provider,
	}, nil
}
