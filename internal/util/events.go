package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var eventLineRegex = regexp.MustCompile(
	`^\s*\[(\d{2}:\d{2}:\d{2}(?:\.\d{3})?)\]\s+(\d+)\s+(\d+)(?:\s+(.*))?$`,
)

func ParseEventLine(line string) (
	eventTime time.Time,
	eventType int,
	competitorID int,
	extraParams []string,
	err error,
) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		err = fmt.Errorf("skipped (empty or comment)")
		return
	}

	m := eventLineRegex.FindStringSubmatch(line)
	if m == nil {
		err = fmt.Errorf("invalid event line format: %q", line)
		return
	}

	if eventTime, err = ParseTime(m[1]); err != nil {
		return
	}
	if eventType, err = strconv.Atoi(m[2]); err != nil {
		err = fmt.Errorf("invalid event type %q: %w", m[2], err)
		return
	}
	if competitorID, err = strconv.Atoi(m[3]); err != nil {
		err = fmt.Errorf("invalid competitor ID %q: %w", m[3], err)
		return
	}
	extraParams = []string{}
	if m[4] != "" {
		extraParams = strings.Fields(m[4])
	}
	return
}

func SplitFields(s string) []string {
	return strings.Fields(strings.TrimSpace(s))
}
