package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tobe/go-type/internal/textgen"
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
		menuModel := NewMenuModel(80, 24)
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

	// Get full book text for paging
	text := textgen.GetFullText()

	// Create and run the Bubble Tea model for typing test
	m := NewModel(text, selectedBook, 80, 24)
	p := tea.NewProgram(m)

	_, err := p.Run()
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func calculateWPM(userInput string, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	wordCount := float64(len(userInput)) / 5.0
	minutes := duration.Minutes()
	return wordCount / minutes
}

func calculateAccuracy(text, userInput string) float64 {
	if len(text) == 0 {
		return 0
	}

	correct := 0
	minLen := len(text)
	if len(userInput) < minLen {
		minLen = len(userInput)
	}

	for i := 0; i < minLen; i++ {
		if text[i] == userInput[i] {
			correct++
		}
	}

	return float64(correct) / float64(len(text)) * 100
}

func calculateErrors(text, userInput string) int {
	errors := 0

	minLen := len(text)
	if len(userInput) < minLen {
		minLen = len(userInput)
	}

	for i := 0; i < minLen; i++ {
		if text[i] != userInput[i] {
			errors++
		}
	}

	if len(text) > len(userInput) {
		errors += len(text) - len(userInput)
	} else if len(userInput) > len(text) {
		errors += len(userInput) - len(text)
	}

	return errors
}
