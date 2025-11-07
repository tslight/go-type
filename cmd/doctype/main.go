package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/godocgen"
	"github.com/tobe/go-type/internal/textgen"
	"github.com/tobe/go-type/pkg/cli"
)

// doctype - typing practice app based on Go documentation
var Version = "unknown"

func main() {
	config := cli.AppConfig{
		Name:            "doctype",
		Version:         Version,
		ListDescription: "List available Go documentation modules",
		ListItems: func() ([]string, error) {
			return godocgen.GetDocumentationNames(), nil
		},
		Configure: []func() error{
			func() error { return godocgen.ConfigureStateFile("doctype") },
		},
		SelectAndLoad: selectDoc,
	}

	if err := cli.RunApp(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func selectDoc(width, height int) (*cli.Selection, error) {
	docNames := godocgen.GetDocumentationNames()
	menuModel := cli.NewDocMenuModel(docNames, width, height)
	program := tea.NewProgram(menuModel)

	if _, err := program.Run(); err != nil {
		return nil, err
	}

	namePtr := menuModel.SelectedDocName()
	if namePtr == nil {
		return nil, nil
	}
	selectedDocName := *namePtr

	text, err := godocgen.GetDocumentation(selectedDocName)
	if err != nil {
		return nil, err
	}

	provider, err := cli.NewDocStateProvider(selectedDocName, len(text))
	if err != nil {
		return nil, err
	}

	selection := &cli.Selection{
		Text:     text,
		Book:     &textgen.Book{ID: 0, Name: selectedDocName},
		Provider: provider,
	}

	return selection, nil
}
