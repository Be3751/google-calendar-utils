package agent

import (
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

type Agent interface {
	FindEventsByCalName(calName string, req FindRequest) (foundEvents *calendar.Events, err error)
}

type agent struct {
	srv *calendar.Service
}

func NewAgent(s *calendar.Service) Agent {
	return &agent{
		srv: s,
	}
}

type FindRequest struct {
	TimeMin time.Time
	TimeMax time.Time
}

func (a *agent) FindEventsByCalName(calName string, req FindRequest) (*calendar.Events, error) {
	calId, err := a.findCalIdByName(calName)
	if err != nil {
		return nil, fmt.Errorf("Unable to find Calendar ID: %w", err)
	}

	events, err := a.srv.Events.List(calId).ShowDeleted(false).
		SingleEvents(true).TimeMin(req.TimeMin.Format(time.RFC3339)).
		TimeMax(req.TimeMax.Format(time.RFC3339)).OrderBy("startTime").MaxResults(1000).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve next ten of the user's events: %w", err)
	}
	return events, nil
}

func (a *agent) findCalIdByName(name string) (string, error) {
	calendarList, err := a.srv.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	for _, cal := range calendarList.Items {
		if cal.Summary == name {
			return cal.Id, nil
		}
	}
	return "", errors.New("No such a location")
}
