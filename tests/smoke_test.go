package tests

import (
	"testing"

	"yadro-impulse/internal/config"
	"yadro-impulse/internal/domain"
)

func TestEndToEnd_SingleAthlete(t *testing.T) {
	cfg := &config.Config{
		Laps: 1, LapLen: 2000, PenaltyLen: 50,
		FiringLines: 1, StartDelta: 10,
	}
	comp := domain.NewCompetition(cfg)

	seq := []domain.Event{
		{Type: 1, CompetitorID: 1},
		{Type: 4, CompetitorID: 1},
		{Type: 10, CompetitorID: 1},
	}
	for _, e := range seq {
		if _, err := comp.ProcessEvent(e); err != nil {
			t.Fatalf("process: %v", err)
		}
	}
	rep := comp.FinalReport()
	if len(rep) != 1 || rep[0].Status != domain.StatusFinished {
		t.Fatalf("bad final report: %+v", rep)
	}
}
