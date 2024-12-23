package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/jsonb"
	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (p *Postgres) CreateTestScheduleTemplate(ctx context.Context, tmpl port.ScheduleTemplate) error {
	payload, err := jsonb.NewFromMarshaler(tmpl.Schedule)
	if err != nil {
		return fmt.Errorf("create jsonb from schedule: %w", err)
	}

	if _, err := sqlx.NamedExecContext(
		ctx,
		p.db,
		`INSERT INTO schedule_template(
			name,
			description,
			schedule_payload
		) VALUES (
			:name,
			:description,
			:schedule_payload
		) ON CONFLICT(name) DO NOTHING`,
		model.ScheduleTemplate{
			Name:            tmpl.Name,
			Description:     tmpl.Description,
			SchedulePayload: payload,
		},
	); err != nil {
		return fmt.Errorf("insert schedule template: %w", err)
	}

	return nil
}
