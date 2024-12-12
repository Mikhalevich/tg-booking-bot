package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (p *Postgres) GetAllTemplates(ctx context.Context) ([]port.ScheduleTemplate, error) {
	var tmpls []model.ScheduleTemplate
	if err := sqlx.SelectContext(
		ctx,
		p.db,
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

	portTmpls, err := model.ToPortScheduleTemplates(tmpls)
	if err != nil {
		return nil, fmt.Errorf("convert to port templates: %w", err)
	}

	return portTmpls, nil
}
