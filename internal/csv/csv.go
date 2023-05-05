package csv

import (
	myCSV "encoding/csv"
	"fmt"
	"time"

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

	header := []string{"title", "datetime-start", "datetime-end", "creator"}
	if err := c.writer.Write(header); err != nil {
		return fmt.Errorf("faild to write record %v: %w", header, err)
	}

	for _, e := range events.Items {
		title := e.Summary
		if title == "" {
			title = "private"
		}
		sDateTime := e.Start.DateTime
		convSDateTime, err := convToFormat(sDateTime)
		if err != nil {
			return fmt.Errorf("Unable to conv: %w", err)
		}
		eDateTime := e.End.DateTime
		convEDateTime, err := convToFormat(eDateTime)
		if err != nil {
			return fmt.Errorf("Unable to conv: %w", err)
		}

		creator := remEMailDomain(e.Creator.Email)

		record := []string{title, convSDateTime, convEDateTime, creator}
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

func convToFormat(rfc3339Str string) (string, error) {
	t, err := time.Parse(time.RFC3339, rfc3339Str)
	if err != nil {
		return "", fmt.Errorf("Unable to parse: %w", err)
	}
	return t.Format("2006/01/02, 15:04"), nil
}
