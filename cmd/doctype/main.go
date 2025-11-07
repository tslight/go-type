package main

import (
	"flag"
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
	list := flag.Bool("l", false, "List available books and their titles")
	listLong := flag.Bool("list", false, "List available books and their titles (long form)")
	version := flag.Bool("v", false, "Show application version")
	versionLong := flag.Bool("version", false, "Show application version (long form)")
	flag.Parse()

	if *version || *versionLong {
		fmt.Println(Version)
		return
	}

	// Handle list books flag
	if *list || *listLong {
		books := textgen.GetAvailableBooks()
		for _, book := range books {
			fmt.Println(book.Name)
		}
		os.Exit(0)
	}

	var selectedDocText string
	var selectedDocName string
	var selectedBook *textgen.Book
	var stateProvider cli.StateProvider

	// Show doc selection menu
	docNames := godocgen.GetDocumentationNames()
	menuModel := cli.NewDocMenuModel(docNames, 80, 24)
	p := tea.NewProgram(menuModel)

	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running menu: %v\n", err)
		os.Exit(1)
	}

	namePtr := menuModel.SelectedDocName()
	if namePtr == nil {
		os.Exit(0)
	}
	selectedDocName = *namePtr

	text, err := godocgen.GetDocumentation(selectedDocName)
	if err != nil {
		fmt.Printf("Error loading documentation %q: %v\n", selectedDocName, err)
		os.Exit(1)
	}
	selectedDocText = text

	selectedBook = &textgen.Book{
		ID:   0,
		Name: selectedDocName,
	}
	provider, err := cli.NewDocStateProvider(selectedDocName, len(selectedDocText))
	if err != nil {
		fmt.Printf("Error preparing state provider: %v\n", err)
		os.Exit(1)
	}
	stateProvider = provider

	// Create and run the Bubble Tea model for typing test
	m := cli.NewModel(selectedDocText, selectedBook, 80, 24, stateProvider)
	typingProgram := tea.NewProgram(m)

	_, err = typingProgram.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
