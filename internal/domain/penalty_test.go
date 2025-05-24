package domain

import "testing"

func TestPenaltyLap_Speed(t *testing.T) {
	p := PenaltyLap{TotalDistance: 150}
	p.StartTime = mustTime("11:00:00.000")
	p.EndTime = mustTime("11:00:30.000")
	if v := p.Speed(); v < 4.9 || v > 5.1 {
		t.Fatalf("unexpected speed %.3f", v)
	}
}
