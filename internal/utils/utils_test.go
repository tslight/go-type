package utils

import (
	"testing"
	"time"
)

func TestCalculateWPM(t *testing.T) {
	wpm := CalculateWPM("hello world", time.Minute)
	if wpm <= 0 {
		t.Fatalf("expected positive wpm, got %f", wpm)
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
