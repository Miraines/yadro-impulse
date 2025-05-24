package domain

import (
	"time"

	"yadro-impulse/internal/util"
)

type PenaltyLap struct {
	StartTime     time.Time
	EndTime       time.Time
	TotalDistance int
	Reason        string
}

func (p PenaltyLap) IsComplete() bool {
	return !p.StartTime.IsZero() && !p.EndTime.IsZero()
}

func (p PenaltyLap) Duration() time.Duration {
	if !p.IsComplete() {
		return 0
	}
	return p.EndTime.Sub(p.StartTime)
}

func (p PenaltyLap) Speed() float64 {
	return util.CalculateSpeed(p.TotalDistance, p.Duration())
}
