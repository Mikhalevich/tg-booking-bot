package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) AddAction(ctx context.Context, info *action.ActionInfo) error {
	return p.addAction(ctx, p.db, info)
}

func (p *Postgres) addAction(
	ctx context.Context,
	e sqlx.ExtContext,
	info *action.ActionInfo,
) error {
	res, err := sqlx.NamedExecContext(
		ctx,
		e,
		`
			INSERT INTO actions(
				employee_id,
				action,
				payload,
				is_completed,
				created_at
			) VALUES (
				:employee_id,
				:action,
				:payload,
				:is_completed,
				:created_at
			)
		`,
		map[string]any{
			"employee_id":  info.EmployeeID,
			"action":       info.Action,
			"payload":      info.Payload,
			"is_completed": false,
			"created_at":   info.CreatedAt,
		})

	if err != nil {
		return fmt.Errorf("exec insert action: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
