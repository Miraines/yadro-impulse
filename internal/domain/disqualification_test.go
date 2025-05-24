package domain

import "testing"

func TestAtEnd_NotStartedDisqualification(t *testing.T) {
	c := NewCompetition(baseCfg())

	events := []Event{
		{Time: mustTime("09:00:00.000"), Type: EventRegistered, CompetitorID: 7},
		{Time: mustTime("09:00:02.000"), Type: EventStartTimeSet, CompetitorID: 7, ExtraParams: []string{"09:00:05.000"}},
	}
	for _, e := range events {
		_, _ = c.ProcessEvent(e)
	}

	c.lastEventTime = mustTime("09:00:40.000")
	gen := c.AtEnd()
	if len(gen) != 1 || gen[0].Type != EventDisqualified {
		t.Fatalf("expected disqualification event, got %+v", gen)
	}
	if c.Competitors[7].Status != StatusNotStarted {
		t.Fatalf("bad status: %v", c.Competitors[7].Status)
	}
}
