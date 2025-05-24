package domain

import "testing"

func TestFiringRange_NumberValidation(t *testing.T) {
	cfg := baseCfg()
	cfg.FiringLines = 2
	c := NewCompetition(cfg)

	_ = mustTime

	evs := []Event{
		{Time: mustTime("10:00:00.000"), Type: EventRegistered, CompetitorID: 1},
		{Time: mustTime("10:01:00.000"), Type: EventOnFiringRange, CompetitorID: 1, ExtraParams: []string{"3"}},
	}
	_, _ = c.ProcessEvent(evs[0])
	if _, err := c.ProcessEvent(evs[1]); err == nil {
		t.Fatal("expected validation error for firing range number")
	}
}
