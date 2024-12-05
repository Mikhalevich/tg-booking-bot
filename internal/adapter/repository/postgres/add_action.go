package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) AddAction(ctx context.Context, info *action.ActionInfo) (int, error) {
	return addAction(ctx, p.db, info)
}

func addAction(
	ctx context.Context,
	e sqlx.ExtContext,
	info *action.ActionInfo,
) (int, error) {
	query, args, err := sqlx.Named(
		`
			INSERT INTO actions(
				employee_id,
				action,
				payload,
				state,
				created_at
			) VALUES (
				:employee_id,
				:action,
				:payload,
				:state,
				:created_at
			)
			RETURNING id
		`,
		map[string]any{
			"employee_id": info.EmployeeID,
			"action":      info.Action,
			"payload":     info.Payload,
			"state":       action.StateInProgress,
			"created_at":  info.CreatedAt,
		})

	if err != nil {
		return 0, fmt.Errorf("prepare named: %w", err)
	}

	var actionID int
	if err := sqlx.GetContext(ctx, e, &actionID, e.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("insert action: %w", err)
	}

	return actionID, nil
}
