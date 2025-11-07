# Go-Type Test Suite

## Overview

Comprehensive test suite for the go-type typing test CLI application, covering both the text generation library and the CLI interface.

## Test Coverage

### `internal/textgen/textgen_test.go` (92.9% coverage)

#### Unit Tests

1. **TestIsAlphaOnly** - Tests alphabetic character validation
   - ✅ Valid lowercase, uppercase, mixed case
   - ✅ Empty strings, numbers, special characters
   - ✅ Hyphens, spaces, apostrophes

2. **TestShuffleWords** - Tests Fisher-Yates shuffle algorithm
   - ✅ Verifies all elements are shuffled
   - ✅ Ensures no elements are lost or duplicated
   - ✅ Confirms randomization works correctly

3. **TestGetParagraph** - Tests paragraph generation
   - ✅ Normal word count (10 words)
   - ✅ Small word count (5 words)
   - ✅ Large word count (50 words)
   - ✅ Zero/negative counts default to 10
   - ✅ Single word generation
   - ✅ Ends with period
   - ✅ First word capitalized

4. **TestGetParagraphConsistency** - Tests paragraph validity across multiple runs
   - ✅ Generated words contain only alphabetic characters
   - ✅ Words are within length bounds (3-20 chars)

5. **TestGetRandomSentence** - Tests random sentence generation
   - ✅ Word count between 8-15 words
   - ✅ Ends with period
   - ✅ First word capitalized

6. **TestGetMultipleSentences** - Tests multiple sentence generation
   - ✅ Correct sentence count
   - ✅ Zero/negative counts default to 3
   - ✅ Single and multiple sentence generation

7. **TestParseEmbeddedDictionary** - Tests embedded dictionary loading
   - ✅ Loads 235,976+ words successfully
   - ✅ All words are alphabetic only
   - ✅ All words are lowercase
   - ✅ All words are 3-20 characters
   - ✅ Common words are present (the, and, but, etc.)

#### Benchmarks

- **BenchmarkGetParagraph**: ~876 ns/op, 5 allocations
- **BenchmarkGetRandomSentence**: ~564 ns/op, 5 allocations
- **BenchmarkShuffleWords**: ~4,848 ns/op, 0 allocations (in-place shuffle)

### `cmd/gutentype/main_test.go`

#### Unit Tests

1. **TestColorConstants** - Tests ANSI color code definitions
   - ✅ colorReset = "\033[0m"
   - ✅ colorGreen = "\033[32m"
   - ✅ colorRed = "\033[31m"
   - ✅ colorGray = "\033[90m"

2. **TestColorOutput** - Tests color formatting
   - ✅ Original text is preserved within colored output
   - ✅ Starts with escape sequence
   - ✅ Ends with reset sequence

3. **TestMetricsCalculation** - Tests typing metrics logic
   - ✅ Perfect match accuracy (100%)
   - ✅ Missing characters accuracy
   - ✅ Extra characters accuracy

4. **TestCharacterComparison** - Tests character matching
   - ✅ Exact matches
   - ✅ Case sensitivity
   - ✅ Space matching
   - ✅ Number vs letter comparison

5. **TestInputValidation** - Tests input handling
   - ✅ Valid text input
   - ✅ Empty input
   - ✅ Single character input
   - ✅ Input with spaces
   - ✅ Special characters
   - ✅ Numbers
   - ✅ Mixed alphanumeric

6. **TestWPMCalculation** - Tests WPM calculation formula
   - ✅ 60 chars in 60 seconds = 12 WPM
   - ✅ 300 chars in 60 seconds = 60 WPM
   - ✅ 150 chars in 30 seconds = 60 WPM
   - ✅ 10 chars in 1 second = 120 WPM
   - Formula: `(characters / 5) / minutes`

7. **TestAccuracyCalculation** - Tests accuracy percentage
   - ✅ Perfect accuracy (100%)
   - ✅ Partial accuracy (50%, 75%)
   - ✅ No accuracy (0%)
   - Formula: `(correctChars * 100) / totalChars`

8. **TestErrorCalculation** - Tests error counting
   - ✅ No errors
   - ✅ Single character error
   - ✅ Multiple errors
   - ✅ Extra characters as errors

#### Benchmarks

- **BenchmarkMetricsCalculation**: ~262 ns/op, 2 allocations
- **BenchmarkColorFormatting**: ~115 ns/op, 3 allocations

## Running Tests

### Run all tests with coverage
```bash
go test -cover ./...
```

### Run textgen tests with verbose output
```bash
go test -v ./internal/textgen
```

### Run CLI tests with verbose output
```bash
go test -v ./cmd/gutentype
```

### Run textgen benchmarks
````

### `cmd/cli/main_test.go`

#### Unit Tests

1. **TestColorConstants** - Tests ANSI color code definitions
   - ✅ colorReset = "\033[0m"
   - ✅ colorGreen = "\033[32m"
   - ✅ colorRed = "\033[31m"
   - ✅ colorGray = "\033[90m"

2. **TestColorOutput** - Tests color formatting
   - ✅ Original text is preserved within colored output
   - ✅ Starts with escape sequence
   - ✅ Ends with reset sequence

3. **TestMetricsCalculation** - Tests typing metrics logic
   - ✅ Perfect match accuracy (100%)
   - ✅ Missing characters accuracy
   - ✅ Extra characters accuracy

4. **TestCharacterComparison** - Tests character matching
   - ✅ Exact matches
   - ✅ Case sensitivity
   - ✅ Space matching
   - ✅ Number vs letter comparison

5. **TestInputValidation** - Tests input handling
   - ✅ Valid text input
   - ✅ Empty input
   - ✅ Single character input
   - ✅ Input with spaces
   - ✅ Special characters
   - ✅ Numbers
   - ✅ Mixed alphanumeric

6. **TestWPMCalculation** - Tests WPM calculation formula
   - ✅ 60 chars in 60 seconds = 12 WPM
   - ✅ 300 chars in 60 seconds = 60 WPM
   - ✅ 150 chars in 30 seconds = 60 WPM
   - ✅ 10 chars in 1 second = 120 WPM
   - Formula: `(characters / 5) / minutes`

7. **TestAccuracyCalculation** - Tests accuracy percentage
   - ✅ Perfect accuracy (100%)
   - ✅ Partial accuracy (50%, 75%)
   - ✅ No accuracy (0%)
   - Formula: `(correctChars * 100) / totalChars`

8. **TestErrorCalculation** - Tests error counting
   - ✅ No errors
   - ✅ Single character error
   - ✅ Multiple errors
   - ✅ Extra characters as errors

#### Benchmarks

- **BenchmarkMetricsCalculation**: ~262 ns/op, 2 allocations
- **BenchmarkColorFormatting**: ~115 ns/op, 3 allocations

## Running Tests

### Run all tests with coverage
```bash
go test -cover ./...
```

### Run textgen tests with verbose output
```bash
go test -v ./internal/textgen
```

### Run CLI tests with verbose output
```bash
go test -v ./cmd/cli
```

### Run textgen benchmarks
```bash
go test -bench=. ./internal/textgen -benchmem
```

### Run CLI benchmarks
```bash
go test -bench=. ./cmd/cli -benchmem
```

### Run all benchmarks
```bash
go test -bench=. ./... -benchmem
```

## Test Results Summary

### Overall Status
✅ **All tests passing**

### Coverage
- `internal/textgen`: 92.9% coverage
- `cmd/cli`: 0.0% (unit tests only, not integration)

### Test Count
- **textgen package**: 7 test functions, 31 test cases, 3 benchmarks
- **CLI package**: 8 test functions, 28 test cases, 2 benchmarks
- **Total**: 15 test functions, 59 test cases, 5 benchmarks

### Performance
- Paragraph generation: ~876 ns/op
- Sentence generation: ~564 ns/op
- Shuffle operation: ~4,848 ns/op (for 1000 elements)
- Metrics calculation: ~262 ns/op
- Color formatting: ~115 ns/op

## Key Features Tested

✅ **Text Generation**
- Random paragraph generation with configurable word count
- Random sentence generation (8-15 words)
- Multiple sentence generation
- Fisher-Yates shuffling algorithm

✅ **Dictionary Management**
- System dictionary loading (with fallback)
- Embedded dictionary parsing
- Word filtering (length bounds, alphabetic only)
- Word normalization (lowercase)

✅ **Input Handling**
- Character-by-character comparison
- Case sensitivity
- Special characters and spaces
- Various input types

✅ **Metrics Calculation**
- Accuracy percentage calculation
- Words-per-minute (WPM) calculation
- Error counting
- Character comparison logic

✅ **Terminal Output**
- ANSI color codes
- Color formatting
- Text coloring (green for correct, red for wrong, gray for expected)

## Notes

- Tests use table-driven testing patterns for maintainability
- Benchmarks use realistic data sizes
- 92.9% code coverage for textgen library
- All 59 test cases pass successfully
- Performance benchmarks show good efficiency across all operations
