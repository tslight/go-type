package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/tobe/go-type/internal/textgen"
	"golang.org/x/term"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
	colorGray  = "\033[90m"
)

func main() {
	wordCount := flag.Int("w", 22, "Number of words to include in the typing test")
	bookID := flag.Int("book", -1, "Book ID to use (see -list for available books)")
	listBooks := flag.Bool("list", false, "List available books and their IDs")
	flag.Parse()

	// Handle list books flag
	if *listBooks {
		books := textgen.GetAvailableBooks()
		fmt.Println("\nAvailable books:")
		for _, book := range books {
			fmt.Printf("  ID %3d: %s\n", book.ID, book.Name)
		}
		fmt.Println()
		os.Exit(0)
	}

	// Load specific book if requested
	if *bookID > 0 {
		if err := textgen.SetBook(*bookID); err != nil {
			fmt.Printf("Error loading book %d: %v\n", *bookID, err)
			os.Exit(1)
		}
	}

	text := textgen.GetParagraph(*wordCount)
	currentBook := textgen.GetCurrentBook()

	// Filter text to ASCII only to avoid UTF-8 encoding issues
	text = toASCII(text)

	fmt.Printf("\nOn your mark, get set, GO TYPE! (Source: %s)\n\n", currentBook.Name)

	// Save cursor position before printing text
	fmt.Print("\033[s")
	// Print the text in gray
	fmt.Print(colorGray + text + colorReset + "\n")

	// Enable raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Restore cursor to saved position
	fmt.Print("\033[u")

	var startTime time.Time
	userInput := ""
	buf := make([]byte, 1)
	testStarted := false

	for {
		n, _ := os.Stdin.Read(buf)
		if n == 0 {
			continue
		}

		char := buf[0]

		// Start timing on first character typed
		if !testStarted && char >= 32 && char < 127 {
			startTime = time.Now()
			testStarted = true
		}

		// Enter to submit
		if char == '\r' || char == '\n' {
			break
		}

		// Ctrl-C to cancel
		if char == 3 {
			term.Restore(int(os.Stdin.Fd()), oldState)
			fmt.Print("\n\033[0J")
			fmt.Print("Cancelled.\n")
			os.Exit(0)
		}

		// Backspace
		if char == 127 || char == 8 {
			if len(userInput) > 0 {
				userInput = userInput[:len(userInput)-1]
				// Move cursor back one position
				fmt.Print("\033[D")
				// Redraw the character that's now at this position
				if len(userInput) < len(text) {
					fmt.Print(colorGray + string(text[len(userInput)]) + colorReset)
					fmt.Print("\033[D")
				}
			}
		} else if char >= 32 && char < 127 {
			userInput += string(char)
			// Compare against the text
			if len(userInput) <= len(text) {
				expectedChar := text[len(userInput)-1]
				if expectedChar == char {
					fmt.Print(colorGreen + string(expectedChar) + colorReset)
				} else {
					// Wrong character - show the expected character in red
					fmt.Print(colorRed + string(expectedChar) + colorReset)
				}
			} else {
				// Extra character beyond the text - show in red
				fmt.Print(colorRed + "+" + colorReset)
			}
		}
	}

	endTime := time.Now()
	var duration time.Duration

	// Only calculate metrics if test was started (first character was typed)
	if testStarted {
		duration = endTime.Sub(startTime)
	} else {
		duration = 0
	}

	// Restore terminal
	term.Restore(int(os.Stdin.Fd()), oldState)

	wpm := calculateWPM(userInput, duration)
	accuracy := calculateAccuracy(text, userInput)
	errors := calculateErrors(text, userInput)

	fmt.Printf("\n\nDuration: %.2f seconds\n", duration.Seconds())
	fmt.Printf("WPM: %.2f\n", wpm)
	fmt.Printf("Accuracy: %.2f%%\n", accuracy)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Typed: %d/%d characters\n", len(userInput), len(text))
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

// toASCII filters out non-ASCII characters to avoid UTF-8 encoding issues
func toASCII(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		if s[i] < 128 {
			result = append(result, s[i])
		}
	}
	return string(result)
}
