package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"yadro-impulse/pkg/errors"
)

const (
	timeFmtMS = "15:04:05.000"
	timeFmt   = "15:04:05"
)

func ParseTime(value string) (time.Time, error) {
	if t, err := time.Parse(timeFmtMS, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(timeFmt, value); err == nil {
		return t, nil
	}
	return time.Time{}, errors.NewInvalidTimeError(
		fmt.Sprintf("invalid time format: %q", value), nil)
}

func ParseDuration(value string) (time.Duration, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, errors.NewInvalidTimeError(
			fmt.Sprintf("invalid duration format: %q", value), nil)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, errors.NewInvalidTimeError(
			fmt.Sprintf("invalid hours in duration: %q", parts[0]), err)
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.NewInvalidTimeError(
			fmt.Sprintf("invalid minutes in duration: %q", parts[1]), err)
	}
	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, errors.NewInvalidTimeError(
			fmt.Sprintf("invalid seconds in duration: %q", parts[2]), err)
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds*float64(time.Second)), nil
}

func FormatTime(t time.Time) string {
	return t.Format(timeFmtMS)
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := d.Seconds() - float64(h*3600+m*60)
	return fmt.Sprintf("%02d:%02d:%06.3f", h, m, s)
}

func CalculateSpeed(distance int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return float64(distance) / duration.Seconds()
}
