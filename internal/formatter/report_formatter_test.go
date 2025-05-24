package formatter

import (
	"strings"
	"testing"
	"time"

	"yadro-impulse/internal/config"
	"yadro-impulse/internal/domain"
)

func fakeComp() domain.Competitor {
	lap := domain.Lap{
		Number:    1,
		StartTime: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		EndTime:   time.Date(0, 1, 1, 9, 10, 0, 0, time.UTC),
		Distance:  3000,
	}
	return domain.Competitor{
		ID:           99,
		Status:       domain.StatusFinished,
		PlannedStart: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
		ActualStart:  time.Date(0, 1, 1, 9, 0, 5, 0, time.UTC),
		FinishTime:   time.Date(0, 1, 1, 9, 10, 0, 0, time.UTC),
		Laps:         []domain.Lap{lap},
	}
}

func TestReportFormatter_Output(t *testing.T) {
	cfg := &config.Config{Laps: 1}
	f := NewReportFormatter(cfg)
	lines := f.Format([]domain.Competitor{fakeComp()})

	if len(lines) != 1 || !strings.Contains(lines[0], "99") {
		t.Fatalf("unexpected line: %v", lines)
	}
	if !strings.Contains(lines[0], "5.000") {
		t.Fatalf("speed missing: %s", lines[0])
	}
}
