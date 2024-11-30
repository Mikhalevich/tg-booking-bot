package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

type actionInfo struct {
	ActionID int           `db:"id"`
	Action   action.Action `db:"action"`
	Payload  []byte        `db:"payload"`
}

func (p *Postgres) GetNextNotCompletedAction(ctx context.Context, employeeID int) (action.ActionInfo, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			action,
			payload
		FROM
			actions
		WHERE
			employee_id = :employee_id AND
			is_completed = :is_completed
		ORDER BY created_at
		LIMIT 1
	`, map[string]any{
		"employee_id":  employeeID,
		"is_completed": false,
	})

	if err != nil {
		return action.ActionInfo{}, fmt.Errorf("named: %w", err)
	}

	var info actionInfo
	if err := sqlx.GetContext(ctx, p.db, &info, p.db.Rebind(query), args...); err != nil {
		return action.ActionInfo{}, fmt.Errorf("select context: %w", err)
	}

	return action.ActionInfo{
		ActionID: info.ActionID,
		Action:   info.Action,
		Payload:  info.Payload,
	}, nil
}

func (p *Postgres) IsActionNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
