package domain

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"yadro-impulse/internal/config"
	"yadro-impulse/internal/util"
	"yadro-impulse/pkg/errors"
)

type Competition struct {
	Config         *config.Config
	Competitors    map[int]*Competitor
	Events         []Event
	lastEventTime  time.Time
	currentFRCount int
	dayOffset      int
}

func NewCompetition(cfg *config.Config) *Competition {
	return &Competition{
		Config:      cfg,
		Competitors: make(map[int]*Competitor),
	}
}

func (c *Competition) ProcessEvent(e Event) ([]Event, error) {
	eventTime := e.Time.Add(time.Duration(c.dayOffset) * 24 * time.Hour)
	if !c.lastEventTime.IsZero() && eventTime.Before(c.lastEventTime) {
		c.dayOffset++
		eventTime = eventTime.Add(24 * time.Hour)
	}
	e.Time = eventTime
	c.lastEventTime = eventTime
	c.Events = append(c.Events, e)

	comp, ok := c.Competitors[e.CompetitorID]
	if !ok {
		if e.Type == EventRegistered {
			comp = &Competitor{ID: e.CompetitorID, Status: StatusRegistered}
			c.Competitors[e.CompetitorID] = comp
		} else {
			return nil, errors.NewCompetitorNotFoundError(e.CompetitorID)
		}
	}
	comp.AddEvent(e)

	var out []Event

	switch e.Type {
	case EventStartTimeSet:
		t, err := util.ParseTime(e.ExtraParams[0])
		if err != nil {
			return nil, errors.NewInvalidEventError("bad start time", err)
		}
		comp.PlannedStart = t
		comp.Status = StatusReady

	case EventOnStartLine:

	case EventStarted:
		comp.ActualStart = e.Time
		comp.Status = StatusInProgress
		comp.Laps = append(comp.Laps, Lap{
			Number:    1,
			StartTime: e.Time,
			Distance:  c.Config.LapLen,
		})

	case EventOnFiringRange:
		if c.currentFRCount >= c.Config.FiringLines {
			return nil, fmt.Errorf("firing lines overflow at %s: %d > %d",
				e.Time.Format("15:04:05.000"), c.currentFRCount+1, c.Config.FiringLines)
		}

		num, err := strconv.Atoi(e.ExtraParams[0])
		if err != nil || num < 1 || num > c.Config.FiringLines {
			return nil, fmt.Errorf("invalid firing range number %q (must be 1..%d)",
				e.ExtraParams[0], c.Config.FiringLines)
		}

		c.currentFRCount++
		fr := FiringRange{Number: num, EntryTime: e.Time}

		if len(comp.Laps) > 0 {
			lastLap := &comp.Laps[len(comp.Laps)-1]
			lastLap.FiringRanges = append(lastLap.FiringRanges, fr)
			comp.CurrentFR = &lastLap.FiringRanges[len(lastLap.FiringRanges)-1]
		}

	case EventTargetHit:
		tgt, err := strconv.Atoi(e.ExtraParams[0])
		if err != nil {
			return nil, fmt.Errorf("invalid target number %q: %w", e.ExtraParams[0], err)
		}
		if comp.CurrentFR != nil {
			comp.CurrentFR.Shots = append(comp.CurrentFR.Shots,
				Shot{Target: tgt, Hit: true, Time: e.Time})
		}

	case EventLeftFiringRange:
		if comp.CurrentFR != nil {
			comp.CurrentFR.ExitTime = e.Time
			misses := comp.CurrentFR.MissesCount()
			comp.PendingPenaltyDistance += misses * c.Config.PenaltyLen

			c.currentFRCount--
			comp.CurrentFR = nil
		}

	case EventEnteredPenalty:
		dist := comp.PendingPenaltyDistance
		if dist == 0 {
			dist = c.Config.PenaltyLen
		}
		pl := PenaltyLap{StartTime: e.Time, TotalDistance: dist}
		comp.Penalties = append(comp.Penalties, pl)
		comp.PendingPenaltyDistance = 0

	case EventLeftPenalty:
		if len(comp.Penalties) > 0 {
			idx := len(comp.Penalties) - 1
			comp.Penalties[idx].EndTime = e.Time
		}

	case EventEndedLap:
		if len(comp.Laps) > 0 {
			last := &comp.Laps[len(comp.Laps)-1]
			last.EndTime = e.Time

			if last.Number < c.Config.Laps {
				comp.Laps = append(comp.Laps, Lap{
					Number:    last.Number + 1,
					StartTime: e.Time,
					Distance:  c.Config.LapLen,
				})
			} else {
				comp.Status = StatusFinished
				comp.FinishTime = e.Time
				out = append(out, Event{
					Time:         e.Time,
					CompetitorID: comp.ID,
					Type:         EventFinished,
					Generated:    true,
				})
			}
		}

	case EventCantContinue:
		comp.Status = StatusNotFinished
		comp.Comment = strings.Join(e.ExtraParams, " ")
		comp.FinishTime = e.Time

	case EventFinished:
		comp.Status = StatusFinished
		comp.FinishTime = e.Time

	case EventDisqualified:
		comp.Status = StatusDisqualified
		comp.FinishTime = e.Time
	}

	return out, nil
}

func (c *Competition) AtEnd() []Event {
	var out []Event

	for _, comp := range c.Competitors {
		if comp.Status == StatusReady || comp.Status == StatusRegistered {
			deadline := comp.PlannedStart.Add(c.Config.StartDelta)
			if c.lastEventTime.After(deadline) {
				evt := Event{
					Time:         deadline,
					CompetitorID: comp.ID,
					Type:         EventDisqualified,
					Generated:    true,
				}
				out = append(out, evt)
				comp.AddEvent(evt)
				comp.Status = StatusNotStarted
				comp.FinishTime = deadline
			}
		}
	}
	return out
}

func (c Competition) FinalReport() []Competitor {
	var list []Competitor
	for _, comp := range c.Competitors {
		list = append(list, *comp)
	}

	sort.Slice(list, func(i, j int) bool {
		ci, cj := list[i], list[j]

		if ci.Status == StatusFinished && cj.Status == StatusFinished {
			return ci.TotalTime() < cj.TotalTime()
		}
		if ci.Status == StatusFinished {
			return true
		}
		if cj.Status == StatusFinished {
			return false
		}
		return ci.ID < cj.ID
	})

	return list
}
