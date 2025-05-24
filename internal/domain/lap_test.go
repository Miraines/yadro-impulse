package domain

import "testing"

func TestLapSpeed(t *testing.T) {
	l := Lap{Distance: 3000}
	l.StartTime = mustTime("10:00:00.000")
	l.EndTime = mustTime("10:05:00.000")
	if spd := l.Speed(); spd < 9.9 || spd > 10.1 {
		t.Fatalf("unexpected speed: %.3f", spd)
	}
}
