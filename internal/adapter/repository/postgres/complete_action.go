package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) CompleteAction(ctx context.Context, id action.ActionID, completedAt time.Time) error {
	return p.completeAction(ctx, id, action.StateCompleted, completedAt)
}

func (p *Postgres) completeAction(
	ctx context.Context,
	id action.ActionID,
	state action.State,
	completedAt time.Time,
) error {
	res, err := sqlx.NamedExecContext(
		ctx,
		p.db,
		`
			UPDATE actions SET
				state = :state,
				completed_at = :completed_at
			WHERE
				id = :id AND
				state = :state_in_progress
		`, map[string]any{
			"id":                id.Int(),
			"state":             state,
			"completed_at":      completedAt,
			"state_in_progress": action.StateInProgress,
		},
	)

	if err != nil {
		return fmt.Errorf("named: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errNotUpdated
	}

	return nil
}
