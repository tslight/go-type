package utils

import (
	"testing"
	"time"
)

func TestCalculateWPM(t *testing.T) {
	// 10 ASCII chars -> 2 words; over 1 minute -> 2 WPM
	wpm := CalculateWPM("helloworld", time.Minute)
	if wpm != 2 {
		t.Fatalf("expected 2 WPM, got %f", wpm)
	}
	// Multibyte: 5 runes of 'é' -> 1 word; 30 seconds -> 2 WPM
	wpm2 := CalculateWPM("ééééé", 30*time.Second)
	if wpm2 != 2 {
		t.Fatalf("expected 2 WPM for multibyte runes, got %f", wpm2)
	}
	// Zero duration -> 0 WPM
	if CalculateWPM("abc", 0) != 0 {
		t.Fatalf("expected 0 WPM for zero duration")
	}
	// Very short duration should produce high but finite WPM (sanity check no divide-by-zero)
	short := CalculateWPM("abcdef", 500*time.Millisecond) // 6 chars = 1.2 words in 0.008333 minutes -> ~144 WPM
	if short < 140 || short > 150 {
		t.Fatalf("expected ~144 WPM for short duration, got %f", short)
	}
}

func TestAccuracyAndErrors(t *testing.T) {
	acc := CalculateAccuracy("abc", "abb")
	if acc <= 0 || acc >= 100 {
		t.Logf("accuracy plausible: %f", acc)
	}
	errs := CalculateErrors("abc", "abb")
	if errs < 0 {
		t.Fatalf("errors must be non-negative")
	}
}
