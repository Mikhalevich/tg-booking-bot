package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type scheduleTemplate struct {
	ID                   int         `db:"id"`
	Name                 string      `db:"name"`
	Description          string      `db:"description"`
	SchedulePayload      jsonb.JSONB `db:"schedule_payload"`
	PayloadVersion       int         `db:"payload_version"`
	PayloadSchemaVersion int         `db:"payload_schema_version"`
	CreatedAt            time.Time   `db:"created_at"`
	UpdatedAt            time.Time   `db:"updated_at"`
}

func (p *Postgres) GetAllTemplates(ctx context.Context) ([]port.ScheduleTemplate, error) {
	var tmpls []scheduleTemplate
	if err := p.db.SelectContext(
		ctx,
		&tmpls,
		`SELECT
			id,
			name,
			description,
			schedule_payload,
			payload_version,
			payload_schema_version,
			created_at,
			updated_at
		FROM
			schedule_template
	`); err != nil {
		return nil, fmt.Errorf("select schedule template: %w", err)
	}

	return nil, nil
}

func ConvertToScheduleTemplate(t scheduleTemplate) (port.ScheduleTemplate, error) {
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
