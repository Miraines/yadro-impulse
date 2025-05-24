package main

import (
	"flag"
	"log"
	"yadro-impulse/internal/config"
	"yadro-impulse/internal/io"
	"yadro-impulse/internal/service"
)

func main() {
	var (
		configFile = flag.String("config", "config.json", "Path to configuration file")
		eventsFile = flag.String("events", "events.txt", "Path to events file")
		logFile    = flag.String("log", "", "Path to output log file (stdout if empty)")
		reportFile = flag.String("report", "", "Path to output report file (stdout if empty)")
	)
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	processor := service.NewProcessor(cfg)

	if err := processor.LoadEvents(*eventsFile); err != nil {
		log.Fatalf("Failed to load events: %v", err)
	}

	logLines, reportLines, err := processor.Run()
	if err != nil {
		log.Fatalf("Failed to process competition: %v", err)
	}

	if err := io.WriteLines(*logFile, logLines); err != nil {
		log.Fatalf("Failed to write log: %v", err)
	}

	if err := io.WriteLines(*reportFile, reportLines); err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}
}
