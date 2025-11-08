package utils

import (
	"time"
	"unicode/utf8"
)

// CalculateWPM calculates words per minute based on user input and duration
func CalculateWPM(userInput string, duration time.Duration) float64 {
	// Guard against zero or near-zero duration to avoid infinities.
	if duration <= 0 {
		return 0
	}
	minutes := duration.Minutes()
	// Count runes (characters), not bytes, so multibyte characters don't inflate WPM.
	runeCount := utf8.RuneCountInString(userInput)
	wordCount := float64(runeCount) / 5.0
	return wordCount / minutes
}

// CalculateAccuracy calculates typing accuracy as a percentage
func CalculateAccuracy(text, userInput string) float64 {
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

// CalculateErrors counts the number of typing errors
func CalculateErrors(text, userInput string) int {
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
