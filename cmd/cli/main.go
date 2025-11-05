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
	flag.Parse()

	text := textgen.GetParagraph(*wordCount)

	fmt.Print("\nOn your mark, get set, GO TYPE:\n\n")

	// Save the cursor position before printing the text
	cursorSave := "\033[s"
	fmt.Print(cursorSave)

	// Print the text in gray
	fmt.Println(colorGray + text + colorReset)

	// Move cursor back to the beginning of the text line
	cursorRestore := "\033[u"
	fmt.Print(cursorRestore)

	// Enable raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
			// Move to the beginning of the text (cursor was saved earlier)
			fmt.Print("\033[u")
			// Clear from cursor to end of screen
			fmt.Print("\033[0J")
			fmt.Print("\nCancelled.\n")
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
			// Print the character from the text with appropriate color
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
