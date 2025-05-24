package domain

import "time"

type EventType int

const (
	EventRegistered      EventType = 1
	EventStartTimeSet    EventType = 2
	EventOnStartLine     EventType = 3
	EventStarted         EventType = 4
	EventOnFiringRange   EventType = 5
	EventTargetHit       EventType = 6
	EventLeftFiringRange EventType = 7
	EventEnteredPenalty  EventType = 8
	EventLeftPenalty     EventType = 9
	EventEndedLap        EventType = 10
	EventCantContinue    EventType = 11

	EventDisqualified EventType = 32
	EventFinished     EventType = 33
)

type Event struct {
	Time         time.Time
	Type         EventType
	CompetitorID int
	ExtraParams  []string
	Generated    bool
}

func (e Event) IsIncoming() bool {
	return e.Type >= 1 && e.Type <= 11
}

func (e Event) IsOutgoing() bool {
	return e.Type == EventDisqualified || e.Type == EventFinished
}
