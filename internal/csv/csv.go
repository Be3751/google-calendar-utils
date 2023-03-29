package csv

import (
	myCSV "encoding/csv"
	"fmt"

	"google.golang.org/api/calendar/v3"
)

type CSV interface {
	WriteFromEvents(evnets *calendar.Events) error
}

type csv struct {
	writer *myCSV.Writer
}

func NewCSV(w *myCSV.Writer) CSV {
	return &csv{
		writer: w,
	}
}

func (c *csv) WriteFromEvents(events *calendar.Events) error {
	defer c.writer.Flush()

	header := []string{"title", "datetime", "creator"}
	if err := c.writer.Write(header); err != nil {
		return fmt.Errorf("faild to write record %v: %w", header, err)
	}

	for _, e := range events.Items {
		title := e.Summary
		date := e.Start.DateTime
		creator := remEMailDomain(e.Creator.Email)

		record := []string{title, date, creator}
		if err := c.writer.Write(record); err != nil {
			return fmt.Errorf("faild to write record %v: %w", record, err)
		}
	}
	if err := c.writer.Error(); err != nil {
		return err
	}
	return nil
}

func remEMailDomain(email string) string {
	for i, s := range email {
		if s == '@' {
			return email[:i]
		}
	}
	return ""
}
