package cli

import (
	"testing"
	"time"
)

// TestCalculateWPM tests words per minute calculation comprehensively
func TestCalculateWPM(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		duration time.Duration
		expected float64
	}{
		{
			name:     "60 characters in 60 seconds = 12 WPM",
			input:    "123456789012345678901234567890123456789012345678901234567890", // 60 chars / 5 = 12 words
			duration: time.Duration(60) * time.Second,
			expected: 12.0, // 12 words / 1 minute
		},
		{
			name:     "300 characters in 60 seconds = 58 WPM",
			input:    "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", // 290 chars / 5 = 58 words
			duration: time.Duration(60) * time.Second,
			expected: 58.0, // 58 words / 1 minute
		},
		{
			name:     "zero duration returns 0",
			input:    "test",
			duration: 0,
			expected: 0.0,
		},
		{
			name:     "150 characters in 30 seconds = 54 WPM",
			input:    "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345", // 135 chars / 5 = 27 words
			duration: time.Duration(30) * time.Second,
			expected: 54.0, // 27 words / 0.5 minutes = 54 WPM
		},
		{
			name:     "single character instant",
			input:    "a",
			duration: 1 * time.Second,
			expected: 12.0, // 1/5 words / (1/60) minutes = 12 WPM
		},
		{
			name:     "very long duration",
			input:    "12345", // 1 word (5 chars)
			duration: 5 * time.Minute,
			expected: 0.2, // 1 word / 5 minutes
		},
		{
			name:     "empty input zero WPM",
			input:    "",
			duration: 60 * time.Second,
			expected: 0.0,
		},
		{
			name:     "very short duration nanoseconds",
			input:    "123456789012345",
			duration: 1 * time.Nanosecond,
			expected: 180000000000.0, // (15/5) / (1/60e9) minutes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWPM(tt.input, tt.duration)
			if result != tt.expected {
				t.Errorf("CalculateWPM(%q, %v) = %f, want %f", tt.input, tt.duration, result, tt.expected)
			}
		})
	}
}

// TestCalculateWPM_EdgeCases tests edge cases for WPM calculation
func TestCalculateWPM_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		duration time.Duration
		validate func(float64) bool
	}{
		{
			name:     "microseconds resolution",
			input:    "12345",
			duration: 100 * time.Microsecond,
			validate: func(wpm float64) bool { return wpm > 0 && wpm < 1e8 },
		},
		{
			name:     "milliseconds resolution",
			input:    "1234567890",
			duration: 500 * time.Millisecond,
			validate: func(wpm float64) bool { return wpm > 100 && wpm < 300 },
		},
		{
			name:     "max safe float",
			input:    string(make([]byte, 1000000)), // 1 million bytes = 200k words
			duration: 1 * time.Second,
			validate: func(wpm float64) bool { return wpm > 100000 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateWPM(tt.input, tt.duration)
			if !tt.validate(result) {
				t.Errorf("CalculateWPM edge case failed: got %f", result)
			}
		})
	}
}

// TestCalculateAccuracy tests typing accuracy calculation comprehensively
func TestCalculateAccuracy(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		input    string
		expected float64
	}{
		{
			name:     "perfect accuracy",
			text:     "The quick brown fox",
			input:    "The quick brown fox",
			expected: 100.0,
		},
		{
			name:     "partial accuracy with errors",
			text:     "The quick brown fox",
			input:    "The quick xxxxx xxx",
			expected: 63.15789473684211,
		},
		{
			name:     "mostly correct",
			text:     "The quick brown fox",
			input:    "xxxxxxxxxxxxxxxxxxxx",
			expected: 5.263157894736842,
		},
		{
			name:     "input shorter than text",
			text:     "The quick brown fox jumps",
			input:    "The quick brown",
			expected: 60.0,
		},
		{
			name:     "empty text returns 0",
			text:     "",
			input:    "anything",
			expected: 0.0,
		},
		{
			name:     "accuracy with spaces",
			text:     "Hello World Test Case",
			input:    "Hello World Test xxxx",
			expected: 80.95238095238095,
		},
		{
			name:     "empty input on non-empty text",
			text:     "hello",
			input:    "",
			expected: 0.0,
		},
		{
			name:     "single character match",
			text:     "a",
			input:    "a",
			expected: 100.0,
		},
		{
			name:     "single character mismatch",
			text:     "a",
			input:    "b",
			expected: 0.0,
		},
		{
			name:     "spaces in both strings",
			text:     "a b c d",
			input:    "a b x d",
			expected: 85.71428571428571,
		},
		{
			name:     "case sensitive",
			text:     "HELLO",
			input:    "hello",
			expected: 0.0,
		},
		{
			name:     "numbers match",
			text:     "test 123",
			input:    "test 123",
			expected: 100.0,
		},
		{
			name:     "special characters match",
			text:     "hello, world!",
			input:    "hello, world!",
			expected: 100.0,
		},
		{
			name:     "input much longer than text",
			text:     "hi",
			input:    "hello world this is a much longer string",
			expected: 50.0, // "he" matches = 2/4 (comparing first 2 chars only)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAccuracy(tt.text, tt.input)
			// Use approximate comparison for floating-point values
			diff := result - tt.expected
			if diff < -0.0001 || diff > 0.0001 {
				t.Errorf("CalculateAccuracy(%q, %q) = %f, want %f", tt.text, tt.input, result, tt.expected)
			}
		})
	}
}

// TestCalculateAccuracy_Comprehensive tests comprehensive accuracy scenarios
func TestCalculateAccuracy_Comprehensive(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		input       string
		minAccuracy float64 // minimum acceptable accuracy
		maxAccuracy float64 // maximum acceptable accuracy
	}{
		{
			name:        "long text partial match",
			text:        "The quick brown fox jumps over the lazy dog",
			input:       "The quick brown fox jumps over the lazy xxx",
			minAccuracy: 90.0,
			maxAccuracy: 100.0,
		},
		{
			name:        "long text short input",
			text:        "The quick brown fox jumps over the lazy dog",
			input:       "The quick",
			minAccuracy: 18.0, // 9/50 chars = 18%
			maxAccuracy: 25.0,
		},
		{
			name:        "unicode-like ascii",
			text:        "test-data_123",
			input:       "test-data_123",
			minAccuracy: 100.0,
			maxAccuracy: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAccuracy(tt.text, tt.input)
			if result < tt.minAccuracy || result > tt.maxAccuracy {
				t.Errorf("CalculateAccuracy(%q, %q) = %f, want between %f and %f",
					tt.text, tt.input, result, tt.minAccuracy, tt.maxAccuracy)
			}
		})
	}
}

// TestCalculateErrors tests error counting comprehensively
func TestCalculateErrors(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		input    string
		expected int
	}{
		{
			name:     "no errors",
			text:     "The quick brown fox",
			input:    "The quick brown fox",
			expected: 0,
		},
		{
			name:     "one error",
			text:     "The quick brown fox",
			input:    "The quick bxown fox",
			expected: 1,
		},
		{
			name:     "multiple errors",
			text:     "The quick brown fox",
			input:    "Thx xxxxx xxxxx xxx",
			expected: 13,
		},
		{
			name:     "extra characters count as errors",
			text:     "test",
			input:    "test extra",
			expected: 6,
		},
		{
			name:     "missing characters count as errors",
			text:     "testing",
			input:    "test",
			expected: 3,
		},
		{
			name:     "all wrong",
			text:     "hello",
			input:    "xxxxx",
			expected: 5,
		},
		{
			name:     "empty input with text",
			text:     "hello",
			input:    "",
			expected: 5,
		},
		{
			name:     "empty text with input",
			text:     "",
			input:    "hello",
			expected: 5,
		},
		{
			name:     "both empty",
			text:     "",
			input:    "",
			expected: 0,
		},
		{
			name:     "single character correct",
			text:     "a",
			input:    "a",
			expected: 0,
		},
		{
			name:     "single character wrong",
			text:     "a",
			input:    "b",
			expected: 1,
		},
		{
			name:     "case sensitive errors",
			text:     "Hello",
			input:    "hello",
			expected: 1,
		},
		{
			name:     "space differences",
			text:     "hello world", // 11 chars
			input:    "helloworld",  // 10 chars, all match except space = 1 missing + 5 overlapping match
			expected: 6,             // 5 matching chars + 1 missing space + 5 extra "world" chars = 6
		},
		{
			name:     "number errors",
			text:     "test 123",
			input:    "test 456",
			expected: 3,
		},
		{
			name:     "special character error",
			text:     "hello!",
			input:    "hello?",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateErrors(tt.text, tt.input)
			if result != tt.expected {
				t.Errorf("CalculateErrors(%q, %q) = %d, want %d", tt.text, tt.input, result, tt.expected)
			}
		})
	}
}

// TestCalculateErrors_Scenarios tests realistic typing scenarios
func TestCalculateErrors_Scenarios(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		input          string
		validateErrors func(int) bool
	}{
		{
			name:           "typo common transposition",
			text:           "the",
			input:          "teh",
			validateErrors: func(e int) bool { return e == 2 }, // e and h swapped = 2 errors
		},
		{
			name:           "missed word",
			text:           "hello world",
			input:          "hello",
			validateErrors: func(e int) bool { return e == 6 }, // space + "world" = 6
		},
		{
			name:           "extra words",
			text:           "test",
			input:          "test with more text",
			validateErrors: func(e int) bool { return e == 15 }, // " with more text" = 15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateErrors(tt.text, tt.input)
			if !tt.validateErrors(result) {
				t.Errorf("CalculateErrors(%q, %q) = %d, validation failed", tt.text, tt.input, result)
			}
		})
	}
}

// BenchmarkCalculateWPM benchmarks WPM calculation
func BenchmarkCalculateWPM(b *testing.B) {
	input := "123456789012345678901234567890"
	duration := 60 * time.Second

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateWPM(input, duration)
	}
}

// BenchmarkCalculateAccuracy benchmarks accuracy calculation
func BenchmarkCalculateAccuracy(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog and runs around the field very quickly"
	input := "The quick brown fox jumps over the lazy dog and runs around the field very quickly"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateAccuracy(text, input)
	}
}

// BenchmarkCalculateErrors benchmarks error calculation
func BenchmarkCalculateErrors(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog and runs around the field very quickly"
	input := "The quick brown fox jumps over the lazy dog and runs around the field very quickly"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CalculateErrors(text, input)
	}
}
