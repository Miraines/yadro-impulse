package domain

import (
	"time"
)

type CompetitorStatus int

const (
	StatusRegistered CompetitorStatus = iota
	StatusReady
	StatusStarted
	StatusInProgress
	StatusFinished
	StatusNotStarted
	StatusNotFinished
	StatusDisqualified
)

func (s CompetitorStatus) String() string {
	return [...]string{
		"Registered",
		"Ready",
		"Started",
		"InProgress",
		"Finished",
		"NotStarted",
		"NotFinished",
		"Disqualified",
	}[s]
}

type Competitor struct {
	ID           int
	Status       CompetitorStatus
	PlannedStart time.Time
	ActualStart  time.Time
	FinishTime   time.Time

	Laps                   []Lap
	Penalties              []PenaltyLap
	Events                 []Event
	Comment                string
	CurrentFR              *FiringRange
	PendingPenaltyDistance int
}

func (c *Competitor) AddEvent(e Event) {
	c.Events = append(c.Events, e)
}

func (c Competitor) TotalTime() time.Duration {
	if c.ActualStart.IsZero() || c.FinishTime.IsZero() {
		return 0
	}
	return c.FinishTime.Sub(c.ActualStart)
}

func (c Competitor) StartDelay() time.Duration {
	if c.PlannedStart.IsZero() || c.ActualStart.IsZero() {
		return 0
	}
	return c.ActualStart.Sub(c.PlannedStart)
}

func (c Competitor) CompletedLaps() []Lap {
	var out []Lap
	for _, lap := range c.Laps {
		if lap.IsComplete() {
			out = append(out, lap)
		}
	}
	return out
}

func (c Competitor) TotalHits() int {
	hits := 0
	for _, lap := range c.Laps {
		for _, fr := range lap.FiringRanges {
			hits += fr.HitsCount()
		}
	}
	return hits
}

func (c Competitor) TotalShots() int {
	shots := 0
	for _, lap := range c.Laps {
		for _, fr := range lap.FiringRanges {
			shots += fr.ShotsCount()
		}
	}
	return shots
}

func (c Competitor) TotalPenaltyTime() time.Duration {
	var sum time.Duration
	for _, p := range c.Penalties {
		if p.IsComplete() {
			sum += p.Duration()
		}
	}
	return sum
}

func (c Competitor) AvgPenaltySpeed() float64 {
	totalDist := 0
	var totalTime time.Duration
	for _, p := range c.Penalties {
		if p.IsComplete() {
			totalDist += p.TotalDistance
			totalTime += p.Duration()
		}
	}
	if totalTime <= 0 {
		return 0
	}
	return float64(totalDist) / totalTime.Seconds()
}
