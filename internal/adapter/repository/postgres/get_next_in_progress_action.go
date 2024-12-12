package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

func (p *Postgres) GetNextInProgressAction(
	ctx context.Context,
	employeeID empl.EmployeeID,
) (action.ActionInfo, error) {
	query, args, err := sqlx.Named(`
		SELECT
			id,
			employee_id,
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
		"employee_id":       employeeID.Int(),
		"state_in_progress": action.StateInProgress,
	})

	if err != nil {
		return action.ActionInfo{}, fmt.Errorf("named: %w", err)
	}

	var info model.ActionInfo
	if err := sqlx.GetContext(ctx, p.db, &info, p.db.Rebind(query), args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return action.ActionInfo{}, errNotFound
		}

		return action.ActionInfo{}, fmt.Errorf("select context: %w", err)
	}

	return model.ToPortActionInfo(info), nil
}
