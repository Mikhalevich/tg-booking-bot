package port

import (
	"context"
	"errors"
	"fmt"
	"time"
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

func (d DayOfWeek) String() string {
	switch d {
	case Mon:
		return "Monday"
	case Tue:
		return "Tuesday"
	case Wed:
		return "Wednesday"
	case Thu:
		return "Thursday"
	case Fri:
		return "Friday"
	case Sat:
		return "Saturday"
	case Sun:
		return "Sunday"
	default:
		return "Invalid"
	}
}

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

func (t Time) String() string {
	return fmt.Sprintf("%d:%02d", t.Hour, t.Minute)
}

type TimeInterval struct {
	Start Time
	End   Time
}

func (ti TimeInterval) Validate() error {
	if err := ti.Start.Validate(); err != nil {
		return fmt.Errorf("ivalid start time: %w", err)
	}

	if err := ti.End.Validate(); err != nil {
		return fmt.Errorf("invalid end time: %w", err)
	}

	if !ti.End.After(ti.Start) {
		return errors.New("start time greater than end time")
	}

	return nil
}

type HoursByDay struct {
	Days  []DayOfWeek
	Hours []TimeInterval
}

type Schedule struct {
	WorkingHours []HoursByDay
}

type ScheduleTemplate struct {
	Name        string
	Description string
	Schedule    Schedule
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ScheduleRepository interface {
	GetAllTemplates(ctx context.Context) ([]ScheduleTemplate, error)
	CreateTestScheduleTemplate(ctx context.Context, tmpl ScheduleTemplate) error
}
