package utils

import (
	"strings"
	"testing"
	"time"
)

// BenchmarkCalculateWPM ensures the WPM calculation remains efficient under large inputs.
func BenchmarkCalculateWPM(b *testing.B) {
	// Build a large synthetic input (~100k characters)
	base := strings.Repeat("abcdefghijklmnopqrstuvwxyz ", 2000) // ~54k chars
	input := base + base                                        // ~108k
	dur := 3 * time.Minute
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateWPM(input, dur)
	}
}
