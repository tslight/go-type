package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
	"github.com/tobe/go-type/pkg/cli"
)

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

	var selectedBook *textgen.Book

	// Always show book selection menu
	menuModel := cli.NewMenuModel(80, 24)
	p := tea.NewProgram(menuModel)

	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running menu: %v\n", err)
		os.Exit(1)
	}

	selectedBook = menuModel.SelectedBook()
	if selectedBook == nil {
		os.Exit(0)
	}

	if err := textgen.SetBook(selectedBook.ID); err != nil {
		fmt.Printf("Error loading book %q: %v\n", selectedBook.Name, err)
		os.Exit(1)
	}

	// Start with a reasonable chunk size for lazy loading
	// Don't load the whole book - we'll load more as needed
	text := textgen.GetFullText()

	stateProvider := cli.NewTextgenStateProvider()

	// Create and run the Bubble Tea model for typing test
	m := cli.NewModel(text, selectedBook, 80, 24, stateProvider)
	typingProgram := tea.NewProgram(m)

	_, err = typingProgram.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
