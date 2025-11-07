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
	docMenu := flag.Bool("m", false, "Show documentation selection menu")
	docFlag := flag.Bool("menu", false, "Show documentation selection menu (long form)")
	listDocs := flag.Bool("l", false, "List available Go documentation modules")
	docList := flag.Bool("list", false, "List available Go documentation modules (long form)")
	version := flag.Bool("version", false, "Show application version")
	flag.Parse()

	if *version {
		fmt.Println(Version)
		return
	}

	// Handle list docs flag
	if *listDocs || *docList {
		docs := godocgen.GetDocumentationNames()
		for _, doc := range docs {
			fmt.Println(doc)
		}
		os.Exit(0)
	}

	var selectedDocText string
	var selectedDocName *string
	var selectedBook *textgen.Book

	// If -m or -menu flag is set, show menu
	if *docMenu || *docFlag {
		// Show doc selection menu
		docNames := godocgen.GetDocumentationNames()
		menuModel := cli.NewDocMenuModel(docNames, 80, 24)
		p := tea.NewProgram(menuModel)

		_, err := p.Run()
		if err != nil {
			fmt.Printf("Error running menu: %v\n", err)
			os.Exit(1)
		}

		selectedDocName = menuModel.SelectedDocName()
		if selectedDocName == nil {
			// User quit without selecting
			os.Exit(0)
		}

		// Load the selected doc
		text, err := godocgen.GetDocumentation(*selectedDocName)
		if err != nil {
			fmt.Printf("Error loading documentation %q: %v\n", *selectedDocName, err)
			os.Exit(1)
		}
		selectedDocText = text

		// Create a Book struct with the doc name for display/tracking
		selectedBook = &textgen.Book{
			ID:   0,
			Name: *selectedDocName,
		}
	} else {
		// Get a random documentation
		text, err := godocgen.GetRandomDocumentation()
		if err != nil {
			fmt.Printf("Error loading documentation: %v\n", err)
			os.Exit(1)
		}
		selectedDocText = text
		// Create a Book struct for tracking (random selection)
		selectedBook = &textgen.Book{
			ID:   0,
			Name: "Random Go Documentation",
		}
	}

	// Create and run the Bubble Tea model for typing test
	m := cli.NewModel(selectedDocText, selectedBook, 80, 24)
	p := tea.NewProgram(m)

	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
