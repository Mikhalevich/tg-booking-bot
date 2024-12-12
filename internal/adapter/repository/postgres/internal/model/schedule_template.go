package model

import (
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type ScheduleTemplate struct {
	ID                   int         `db:"id"`
	Name                 string      `db:"name"`
	Description          string      `db:"description"`
	SchedulePayload      jsonb.JSONB `db:"schedule_payload"`
	PayloadVersion       int         `db:"payload_version"`
	PayloadSchemaVersion int         `db:"payload_schema_version"`
	CreatedAt            time.Time   `db:"created_at"`
	UpdatedAt            time.Time   `db:"updated_at"`
}

func ToPortScheduleTemplates(tmpls []ScheduleTemplate) ([]port.ScheduleTemplate, error) {
	if len(tmpls) == 0 {
		return nil, nil
	}

	output := make([]port.ScheduleTemplate, 0, len(tmpls))

	for _, t := range tmpls {
		portTempl, err := ToPortScheduleTemplate(t)
		if err != nil {
			return nil, fmt.Errorf("convert to port template: %w", err)
		}

		output = append(output, portTempl)
	}

	return output, nil
}

func ToPortScheduleTemplate(t ScheduleTemplate) (port.ScheduleTemplate, error) {
	var sch port.Schedule
	if err := jsonb.ConvertTo(t.SchedulePayload, &sch); err != nil {
		return port.ScheduleTemplate{}, fmt.Errorf("convert payload: %w", err)
	}

	return port.ScheduleTemplate{
		Name:        t.Name,
		Description: t.Description,
		Schedule:    sch,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}
