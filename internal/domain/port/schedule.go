package port

import (
	"context"
	"errors"
	"fmt"
)

type DayOfWeek int

const (
	Mon DayOfWeek = iota + 1
	Tue
	Wed
	Thu
	Fri
	Sat
	Sun
)

type Time struct {
	Hour   int
	Minute int
}

func (t Time) After(u Time) bool {
	if t.Hour == u.Hour {
		return t.Minute > u.Minute
	}

	return t.Hour > u.Hour
}

func (t Time) Validate() error {
	if t.Hour < 0 || t.Hour > 23 {
		return fmt.Errorf("invalid hour %d", t.Hour)
	}

	if t.Minute < 0 || t.Minute > 60 {
		return fmt.Errorf("invalid minute %d", t.Minute)
	}

	return nil
}

type TimeInterval struct {
	Start Time
	End   Time
}

func (ti TimeInterval) Validate() error {
	if !ti.End.After(ti.Start) {
		return errors.New("start time greater than end time")
	}

	return nil
}

type HoursByDay struct {
	Day   DayOfWeek
	Hours []TimeInterval
}

type Schedule struct {
	WorkingHours []HoursByDay
}

type ScheduleRepository interface {
	GetAllTemplates(ctx context.Context) ([]Schedule, error)
}
