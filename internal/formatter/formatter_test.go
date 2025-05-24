package formatter

import (
	"strings"
	"testing"
	"time"

	"yadro-impulse/internal/domain"
)

func TestLogFormatter(t *testing.T) {
	f := NewLogFormatter()
	ev := domain.Event{Time: time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC), Type: domain.EventRegistered, CompetitorID: 42}
	out := f.Format(ev)
	if !strings.Contains(out, "competitor(42) registered") {
		t.Fatalf("unexpected log line: %s", out)
	}
}
