package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"yadro-impulse/pkg/errors"
)

type Config struct {
	Laps        int           `json:"laps"`
	LapLen      int           `json:"lapLen"`
	PenaltyLen  int           `json:"penaltyLen"`
	FiringLines int           `json:"firingLines"`
	Start       time.Time     `json:"start"`
	StartDelta  time.Duration `json:"startDelta"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type alias Config
	aux := struct {
		Start      string `json:"start"`
		StartDelta string `json:"startDelta"`
		*alias
	}{alias: (*alias)(c)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.NewInvalidConfigError("failed to decode config JSON", err)
	}

	// Validate positive values
	if c.Laps <= 0 {
		return errors.NewInvalidConfigError("laps must be positive", nil)
	}
	if c.LapLen <= 0 {
		return errors.NewInvalidConfigError("lapLen must be positive", nil)
	}
	if c.PenaltyLen <= 0 {
		return errors.NewInvalidConfigError("penaltyLen must be positive", nil)
	}
	if c.FiringLines <= 0 {
		return errors.NewInvalidConfigError("firingLines must be positive", nil)
	}

	startTime, err := time.Parse("15:04:05", aux.Start)
	if err != nil {
		return errors.NewInvalidConfigError(fmt.Sprintf("invalid start time format: %s", aux.Start), err)
	}
	c.Start = startTime

	delta, err := parseHMS(aux.StartDelta)
	if err != nil {
		return errors.NewInvalidConfigError(fmt.Sprintf("invalid startDelta format: %s", aux.StartDelta), err)
	}
	c.StartDelta = delta

	return nil
}

func parseHMS(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("duration %q: expected format HH:MM:SS", s)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}

	dur := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds*float64(time.Second))
	return dur, nil
}
