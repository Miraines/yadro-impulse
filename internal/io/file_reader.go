package io

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/util"
)

func ReadEvents(filename string) ([]domain.Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open events.txt.txt file %q: %w", filename, err)
	}
	defer file.Close()

	var events []domain.Event
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		tm, etypeInt, compID, extra, err := util.ParseEventLine(line)
		if err != nil {
			if strings.HasPrefix(err.Error(), "skipped") {
				continue
			}
			return nil, fmt.Errorf("error parsing line %d %q: %w", lineNum, line, err)
		}

		events = append(events, domain.Event{
			Time:         tm,
			Type:         domain.EventType(etypeInt),
			CompetitorID: compID,
			ExtraParams:  extra,
			Generated:    false,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %q: %w", filename, err)
	}
	return events, nil
}
