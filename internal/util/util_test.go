package util

import "testing"

func TestParseTime_Valid(t *testing.T) {
	got, err := ParseTime("12:34:56.789")
	if err != nil || got.Hour() != 12 || got.Nanosecond()/1e6 != 789 {
		t.Fatalf("bad parse: %v, %v", got, err)
	}
}

func TestParseTime_Bad(t *testing.T) {
	if _, err := ParseTime("99:00:00"); err == nil {
		t.Fatal("expected error for invalid hour")
	}
}

func TestFormatDuration_Roundtrip(t *testing.T) {
	d, _ := ParseDuration("01:02:03.500")
	str := FormatDuration(d)
	if str != "01:02:03.500" {
		t.Fatalf("expected 01:02:03.500, got %s", str)
	}
}

func TestEventLine(t *testing.T) {
	line := "[09:30:00.000] 4 7"
	_, et, id, _, err := ParseEventLine(line)
	if err != nil || et != 4 || id != 7 {
		t.Fatalf("parse failed: %v", err)
	}
}
