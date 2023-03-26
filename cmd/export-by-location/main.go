package main

import (
	"fmt"
	"log"
	"time"

	"google-calendar-utils/internal/agent"
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

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
}
