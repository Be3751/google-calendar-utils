package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"google-calendar-utils/internal/agent"
	myCSV "google-calendar-utils/internal/csv"
	"google-calendar-utils/pkg/calendar"
)

type arg struct {
	CalendarName  string
	DurationStart time.Time
	DurationEnd   time.Time
}

const (
	calendarNameFlag  = "-c"
	durationStartFlag = "-ds"
	durationEndFlag   = "-de"
)

func main() {
	arg, err := parseArgs(os.Args)
	if err != nil {
		log.Fatalf("failed to parse args: %v", err)
	}

	srv, err := calendar.NewService("secrets/credentials.json")
	if err != nil {
		log.Fatalf("unable to retrieve Calendar client: %v", err)
	}

	req := agent.FindRequest{
		TimeMin: arg.DurationStart,
		TimeMax: arg.DurationEnd,
	}

	agent := agent.NewAgent(srv)
	events, err := agent.FindEventsByCalName(arg.CalendarName, req)
	if err != nil {
		log.Fatalf("unable to find events by %s %v", arg.CalendarName, err)
	}

	savedPath := "out/result.csv"
	f, err := os.Create(savedPath)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	csv := myCSV.NewCSV(writer)
	if err := csv.WriteFromEvents(events); err != nil {
		log.Fatalf("failed to write csv: %v", err)
	}
	fmt.Printf("saved %s \n", savedPath)
}

func parseArgs(args []string) (*arg, error) {
	arg := arg{}
	for i, a := range args {
		switch a {
		case calendarNameFlag:
			arg.CalendarName = os.Args[i+1]
		case durationStartFlag:
			ds, err := time.Parse(time.RFC3339, os.Args[i+1]+":00+09:00")
			if err != nil {
				return nil, fmt.Errorf("failed to parse the ds value %s: %v", os.Args[i+1], err)
			}
			arg.DurationStart = ds
		case durationEndFlag:
			de, err := time.Parse(time.RFC3339, os.Args[i+1]+":00+09:00")
			if err != nil {
				return nil, fmt.Errorf("failed to parse the de value %s: %v", os.Args[i+1], err)
			}
			arg.DurationEnd = de
		}
	}
	return &arg, nil
}
