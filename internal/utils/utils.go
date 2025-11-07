package utils

import "time"

// CalculateWPM calculates words per minute based on user input and duration
func CalculateWPM(userInput string, duration time.Duration) float64 {
	if duration.Seconds() == 0 {
		return 0
	}
	wordCount := float64(len(userInput)) / 5.0
	minutes := duration.Minutes()
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
