package domain

import "time"

const RequiredShots = 5

type Shot struct {
	Target int
	Hit    bool
	Time   time.Time
}

type FiringRange struct {
	Number    int
	EntryTime time.Time
	ExitTime  time.Time
	Shots     []Shot
}

func (f FiringRange) IsComplete() bool {
	return !f.EntryTime.IsZero() && !f.ExitTime.IsZero()
}

func (f FiringRange) ShotsCount() int {
	if len(f.Shots) < RequiredShots {
		return RequiredShots
	}
	return len(f.Shots)
}

func (f FiringRange) HitsCount() int {
	count := 0
	for _, shot := range f.Shots {
		if shot.Hit {
			count++
		}
	}
	return count
}

func (f FiringRange) MissesCount() int {
	misses := RequiredShots - f.HitsCount()
	if misses < 0 {
		return 0
	}
	return misses
}
