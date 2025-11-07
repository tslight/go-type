package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
	"github.com/tobe/go-type/pkg/cli"
)

// Color constants for tests
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorGray  = "\033[90m"
)

func main() {
	bookMenu := flag.Bool("b", false, "Show book selection menu")
	bookFlag := flag.Bool("book", false, "Show book selection menu (long form)")
	listBooks := flag.Bool("l", false, "List available books and their titles")
	bookList := flag.Bool("list", false, "List available books and their titles (long form)")
	flag.Parse()

	// Handle list books flag
	if *listBooks || *bookList {
		books := textgen.GetAvailableBooks()
		for _, book := range books {
			fmt.Println(book.Name)
		}
		os.Exit(0)
	}

	var selectedBook *textgen.Book

	// If -b or -book flag is set, show menu
	if *bookMenu || *bookFlag {
		// Show book selection menu
		menuModel := cli.NewMenuModel(80, 24)
		p := tea.NewProgram(menuModel)

		_, err := p.Run()
		if err != nil {
			fmt.Printf("Error running menu: %v\n", err)
			os.Exit(1)
		}

		selectedBook = menuModel.SelectedBook()
		if selectedBook == nil {
			// User quit without selecting
			os.Exit(0)
		}

		// Load the selected book
		if err := textgen.SetBook(selectedBook.ID); err != nil {
			fmt.Printf("Error loading book %q: %v\n", selectedBook.Name, err)
			os.Exit(1)
		}
	} else {
		// Default: pick a random book
		selectedBook = textgen.GetCurrentBook()
	}

	// Start with a reasonable chunk size for lazy loading
	// Don't load the whole book - we'll load more as needed
	text := textgen.GetFullText()

	// Create and run the Bubble Tea model for typing test
	m := cli.NewModel(text, selectedBook, 80, 24)
	p := tea.NewProgram(m)

	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
