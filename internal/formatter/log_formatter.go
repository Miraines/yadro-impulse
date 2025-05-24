package formatter

import (
	"fmt"
	"strings"

	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/util"
)

type LogFormatter struct{}

func NewLogFormatter() *LogFormatter {
	return &LogFormatter{}
}

func (f *LogFormatter) Format(e domain.Event) string {
	t := util.FormatTime(e.Time)

	switch e.Type {
	case domain.EventRegistered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", t, e.CompetitorID)

	case domain.EventStartTimeSet:
		return fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s",
			t, e.CompetitorID, e.ExtraParams[0])

	case domain.EventOnStartLine:
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", t, e.CompetitorID)

	case domain.EventStarted:
		return fmt.Sprintf("[%s] The competitor(%d) has started", t, e.CompetitorID)

	case domain.EventOnFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)",
			t, e.CompetitorID, e.ExtraParams[0])

	case domain.EventTargetHit:
		return fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)",
			t, e.ExtraParams[0], e.CompetitorID)

	case domain.EventLeftFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", t, e.CompetitorID)

	case domain.EventEnteredPenalty:
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", t, e.CompetitorID)

	case domain.EventLeftPenalty:
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", t, e.CompetitorID)

	case domain.EventEndedLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", t, e.CompetitorID)

	case domain.EventCantContinue:
		return fmt.Sprintf("[%s] The competitor(%d) can't continue: %s",
			t, e.CompetitorID, strings.Join(e.ExtraParams, " "))

	case domain.EventFinished:
		return fmt.Sprintf("[%s] The competitor(%d) has finished", t, e.CompetitorID)

	case domain.EventDisqualified:
		return fmt.Sprintf("[%s] The competitor(%d) is disqualified", t, e.CompetitorID)

	default:
		return fmt.Sprintf("[%s] Unknown event %d for competitor(%d)", t, e.Type, e.CompetitorID)
	}
}
