package domain

import "testing"

func TestProcessEvent_MidnightHandled(t *testing.T) {
	c := NewCompetition(baseCfg())
	evt1 := Event{Time: mustTime("23:59:59.900"), Type: EventRegistered, CompetitorID: 1}
	evt2 := Event{Time: mustTime("00:00:00.100"), Type: EventStarted, CompetitorID: 1}

	if _, err := c.ProcessEvent(evt1); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if _, err := c.ProcessEvent(evt2); err != nil {
		t.Fatalf("midnight rollover incorrectly treated as error: %v", err)
	}
	if c.dayOffset != 1 {
		t.Fatalf("dayOffset not incremented, got %d", c.dayOffset)
	}
}
