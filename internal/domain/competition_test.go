package domain

import (
	"testing"
	"time"

	"yadro-impulse/internal/config"
)

func mustTime(s string) time.Time {
	t, _ := time.Parse("15:04:05.000", s)
	return t
}

func baseCfg() *config.Config {
	return &config.Config{
		Laps:        1,
		LapLen:      1000,
		PenaltyLen:  100,
		FiringLines: 1,
		Start:       mustTime("09:00:00.000"),
		StartDelta:  30 * time.Second,
	}
}

func TestCompetition_SimpleFinish(t *testing.T) {
	c := NewCompetition(baseCfg())

	seq := []Event{
		{Time: mustTime("09:00:00.000"), Type: EventRegistered, CompetitorID: 1},
		{Time: mustTime("09:00:10.000"), Type: EventStartTimeSet, CompetitorID: 1, ExtraParams: []string{"09:00:20.000"}},
		{Time: mustTime("09:00:20.000"), Type: EventOnStartLine, CompetitorID: 1},
		{Time: mustTime("09:00:22.000"), Type: EventStarted, CompetitorID: 1},
		{Time: mustTime("09:02:22.000"), Type: EventEndedLap, CompetitorID: 1},
	}

	for _, e := range seq {
		if _, err := c.ProcessEvent(e); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(c.Competitors) != 1 {
		t.Fatalf("competitor missing")
	}
	comp := c.Competitors[1]
	if comp.Status != StatusFinished {
		t.Fatalf("expected finished, got %v", comp.Status)
	}
	if comp.TotalTime() != 2*time.Minute {
		t.Fatalf("bad total time: %v", comp.TotalTime())
	}
}

func TestCompetition_PenaltyCalculation(t *testing.T) {
	cfg := baseCfg()
	c := NewCompetition(cfg)

	events := []Event{
		{Time: mustTime("09:00:00.000"), Type: EventRegistered, CompetitorID: 1},
		{Time: mustTime("09:00:10.000"), Type: EventStartTimeSet, CompetitorID: 1, ExtraParams: []string{"09:00:15.000"}},
		{Time: mustTime("09:00:15.000"), Type: EventOnStartLine, CompetitorID: 1},
		{Time: mustTime("09:00:16.000"), Type: EventStarted, CompetitorID: 1},
		{Time: mustTime("09:01:00.000"), Type: EventOnFiringRange, CompetitorID: 1, ExtraParams: []string{"1"}},
		{Time: mustTime("09:01:01.000"), Type: EventTargetHit, CompetitorID: 1, ExtraParams: []string{"1"}},
		{Time: mustTime("09:01:02.000"), Type: EventTargetHit, CompetitorID: 1, ExtraParams: []string{"2"}},
		{Time: mustTime("09:01:03.000"), Type: EventTargetHit, CompetitorID: 1, ExtraParams: []string{"3"}},
		{Time: mustTime("09:01:04.000"), Type: EventTargetHit, CompetitorID: 1, ExtraParams: []string{"4"}},
		{Time: mustTime("09:01:05.000"), Type: EventLeftFiringRange, CompetitorID: 1},
		{Time: mustTime("09:01:10.000"), Type: EventEnteredPenalty, CompetitorID: 1},
		{Time: mustTime("09:02:00.000"), Type: EventLeftPenalty, CompetitorID: 1},
		{Time: mustTime("09:04:00.000"), Type: EventEndedLap, CompetitorID: 1},
	}

	for _, e := range events {
		if _, err := c.ProcessEvent(e); err != nil {
			t.Fatalf("process error: %v", err)
		}
	}

	comp := c.Competitors[1]
	if comp.TotalHits() != 4 || comp.TotalShots() != 5 {
		t.Fatalf("hits/shots mismatch: %d/%d", comp.TotalHits(), comp.TotalShots())
	}
	if len(comp.Penalties) != 1 || comp.Penalties[0].TotalDistance != 100 {
		t.Fatalf("penalty distance incorrect: %+v", comp.Penalties)
	}
}

func TestCompetition_FiringLineOverflow(t *testing.T) {
	cfg := baseCfg()
	cfg.FiringLines = 1
	c := NewCompetition(cfg)

	ev1 := Event{Time: mustTime("09:00:00.000"), Type: EventRegistered, CompetitorID: 1}
	ev2 := Event{Time: mustTime("09:00:00.100"), Type: EventRegistered, CompetitorID: 2}
	_, _ = c.ProcessEvent(ev1)
	_, _ = c.ProcessEvent(ev2)

	onRange1 := Event{Time: mustTime("09:01:00.000"), Type: EventOnFiringRange, CompetitorID: 1, ExtraParams: []string{"1"}}
	onRange2 := Event{Time: mustTime("09:01:00.500"), Type: EventOnFiringRange, CompetitorID: 2, ExtraParams: []string{"1"}}
	if _, err := c.ProcessEvent(onRange1); err != nil {
		t.Fatalf("first entry failed: %v", err)
	}
	if _, err := c.ProcessEvent(onRange2); err == nil {
		t.Fatal("expected firing line overflow error, got nil")
	}
}

func TestCompetition_MidnightRollover(t *testing.T) {
	cfg := baseCfg()
	c := NewCompetition(cfg)

	seq := []Event{
		{Time: mustTime("23:59:50.000"), Type: EventRegistered, CompetitorID: 1},
		{Time: mustTime("23:59:55.000"), Type: EventStarted, CompetitorID: 1},
		{Time: mustTime("00:00:10.000"), Type: EventEndedLap, CompetitorID: 1},
	}
	for _, e := range seq {
		if _, err := c.ProcessEvent(e); err != nil {
			t.Fatalf("midnight error: %v", err)
		}
	}
	comp := c.Competitors[1]
	if comp.TotalTime() != 15*time.Second {
		t.Fatalf("expected 15s, got %v", comp.TotalTime())
	}
}
