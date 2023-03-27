package csv

import (
	myCSV "encoding/csv"
	"fmt"

	"google.golang.org/api/calendar/v3"
)

type CSV interface {
	Write(evnets *calendar.Events) error
}

type csv struct {
	writer *myCSV.Writer
}

func NewCSV(w *myCSV.Writer) CSV {
	return &csv{
		writer: w,
	}
}

func (c *csv) Write(events *calendar.Events) error {
	defer c.writer.Flush()

	header := []string{"title", "datetime", "attendees"}
	if err := c.writer.Write(header); err != nil {
		return fmt.Errorf("faild to write record %v: %w", header, err)
	}

	for _, e := range events.Items {
		title := e.Summary
		date := e.Start.DateTime
		var attendees string
		record := []string{title, date, attendees}
		if err := c.writer.Write(record); err != nil {
			return fmt.Errorf("faild to write record %v: %w", record, err)
		}
	}
	if err := c.writer.Error(); err != nil {
		return err
	}
	return nil
}

func commaSeparatedStr(attendees []*calendar.EventAttendee) string {
	result := ""
	for idx, elem := range attendees {
		fmt.Println(elem)
		if idx == 0 {
			if elem.DisplayName == "" {
				result = fmt.Sprint(elem.Email)
			} else {
				result = fmt.Sprint(elem.DisplayName)
			}
		} else {
			if elem.DisplayName == "" {
				result = result + fmt.Sprintf(", %v", elem.Email)
			} else {
				result = result + fmt.Sprintf(", %v", elem.DisplayName)
			}
		}
	}
	return result
}
