package service

import (
	"fmt"

	"yadro-impulse/internal/config"
	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/formatter"
	ioPkg "yadro-impulse/internal/io"
)

type Processor struct {
	cfg     *config.Config
	events  []domain.Event
	comp    *domain.Competition
	logFmt  *formatter.LogFormatter
	reportF *formatter.ReportFormatter
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		cfg:     cfg,
		comp:    domain.NewCompetition(cfg),
		logFmt:  formatter.NewLogFormatter(),
		reportF: formatter.NewReportFormatter(cfg),
	}
}

func (p *Processor) LoadEvents(path string) error {
	evs, err := ioPkg.ReadEvents(path)
	if err != nil {
		return fmt.Errorf("failed to load events.txt.txt: %w", err)
	}
	p.events = evs
	return nil
}

func (p *Processor) Run() ([]string, []string, error) {
	var logLines []string

	for _, e := range p.events {
		outgoing, err := p.comp.ProcessEvent(e)
		if err != nil {
			return nil, nil, fmt.Errorf("error processing event %v: %w", e, err)
		}

		logLines = append(logLines, p.logFmt.Format(e))
		for _, oe := range outgoing {
			logLines = append(logLines, p.logFmt.Format(oe))
		}
	}

	for _, oe := range p.comp.AtEnd() {
		logLines = append(logLines, p.logFmt.Format(oe))
	}

	competitors := p.comp.FinalReport()
	reportLines := p.reportF.Format(competitors)

	return logLines, reportLines, nil
}
