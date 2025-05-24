package util

import "testing"

func TestFormatDuration_Long(t *testing.T) {
	d, _ := ParseDuration("12:00:00.000")
	if s := FormatDuration(d); s != "12:00:00.000" {
		t.Fatalf("expected 12:00:00.000 got %s", s)
	}
}
