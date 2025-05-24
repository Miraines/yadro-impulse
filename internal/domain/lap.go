package domain

import (
	"time"
	"yadro-impulse/internal/util"
)

type Lap struct {
	Number       int
	StartTime    time.Time
	EndTime      time.Time
	Distance     int
	FiringRanges []FiringRange
}

func (l Lap) IsComplete() bool {
	return !l.StartTime.IsZero() && !l.EndTime.IsZero()
}

func (l Lap) Duration() time.Duration {
	if !l.IsComplete() {
		return 0
	}
	return l.EndTime.Sub(l.StartTime)
}

func (l Lap) Speed() float64 {
	return util.CalculateSpeed(l.Distance, l.Duration())
}
