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
	State    action.State  `db:"state"`
}

func (p *Postgres) GetNextInProgressAction(ctx context.Context, employeeID int) (action.ActionInfo, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			action,
			payload,
			state
		FROM
			actions
		WHERE
			employee_id = :employee_id AND
			state = :state_in_progress
		ORDER BY created_at
		LIMIT 1
	`, map[string]any{
		"employee_id":       employeeID,
		"state_in_progress": action.StateInProgress,
	})

	if err != nil {
		return action.ActionInfo{}, fmt.Errorf("named: %w", err)
	}

	var info actionInfo
	if err := sqlx.GetContext(ctx, p.db, &info, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return action.ActionInfo{}, errNotFound
		}

		return action.ActionInfo{}, fmt.Errorf("select context: %w", err)
	}

	return action.ActionInfo{
		ActionID:   info.ActionID,
		EmployeeID: employeeID,
		Action:     info.Action,
		Payload:    info.Payload,
		State:      info.State,
	}, nil
}
