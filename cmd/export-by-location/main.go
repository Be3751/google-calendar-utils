package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"google-calendar-utils/internal/agent"
	myCSV "google-calendar-utils/internal/csv"
	"google-calendar-utils/pkg/calendar"
)

func main() {
	// TODO: コマンドライン引数を受け取る処理
	calendarName := "KC111"

	srv, err := calendar.NewService("secrets/credentials.json")
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	req := agent.FindRequest{
		TimeMin: time.Date(2020, time.April, 1, 0, 0, 0, 0, time.FixedZone("JST", 0)),
		TimeMax: time.Date(2023, time.April, 1, 0, 0, 0, 0, time.FixedZone("JST", 0)),
	}

	agent := agent.NewAgent(srv)
	events, err := agent.FindEventsByCalName(calendarName, req)
	if err != nil {
		log.Fatalf("Unable to find events by %s %v", calendarName, err)
	}

	f, err := os.Create("out/sample.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	csv := myCSV.NewCSV(writer)
	if err := csv.Write(events); err != nil {
		panic(err)
	}
}
