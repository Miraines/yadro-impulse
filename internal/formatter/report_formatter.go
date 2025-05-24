package formatter

import (
	"fmt"
	"strings"

	"yadro-impulse/internal/config"
	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/util"
)

type ReportFormatter struct {
	laps int
}

func NewReportFormatter(cfg *config.Config) *ReportFormatter {
	return &ReportFormatter{laps: cfg.Laps}
}

func (f *ReportFormatter) Format(comps []domain.Competitor) []string {
	var lines []string
	for _, c := range comps {
		var prefix string
		switch c.Status {
		case domain.StatusNotStarted:
			prefix = "[NotStarted]"
		case domain.StatusNotFinished:
			prefix = "[NotFinished]"
		case domain.StatusDisqualified:
			prefix = "[Disqualified]"
		default:
			totalTime := c.TotalTime() + c.StartDelay()
			prefix = util.FormatDuration(totalTime)
		}

		laps := c.CompletedLaps()
		var lapStats []string
		for i := 1; i <= f.laps; i++ {
			if i <= len(laps) {
				dur := util.FormatDuration(laps[i-1].Duration())
				spd := fmt.Sprintf("%.3f", laps[i-1].Speed())
				lapStats = append(lapStats, fmt.Sprintf("{%s, %s}", dur, spd))
			} else {
				lapStats = append(lapStats, "{,}")
			}
		}

		var penalty string
		if c.TotalPenaltyTime() == 0 || c.AvgPenaltySpeed() == 0 {
			penalty = "{,}"
		} else {
			penDur := util.FormatDuration(c.TotalPenaltyTime())
			penSpd := fmt.Sprintf("%.3f", c.AvgPenaltySpeed())
			penalty = fmt.Sprintf("{%s, %s}", penDur, penSpd)
		}

		h, s := c.TotalHits(), c.TotalShots()
		line := fmt.Sprintf("%s %d [%s] %s %d/%d",
			prefix, c.ID, strings.Join(lapStats, ", "), penalty, h, s,
		)
		lines = append(lines, line)
	}
	return lines
}
