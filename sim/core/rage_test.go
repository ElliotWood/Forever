package core

import "testing"

// `^` is XOR in Go, not a power, so the level term used to read 60^2 = 62 instead of 3600 and
// every warrior generated 16% too much rage.
func TestRageConversionAtLevel60(t *testing.T) {
	if got := GetRageConversion(60); got < 230.5 || got > 230.7 {
		t.Fatalf("rage conversion at 60 = %v, want 230.6", got)
	}
}
