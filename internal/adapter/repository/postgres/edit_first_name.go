package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) EditFirstName(
	ctx context.Context,
	nameInfo port.EditNameInput,
	nextAction *action.ActionInfo,
) error {
	if err := transaction.Transaction(ctx, p.db,
		func(ctx context.Context, tx *sqlx.Tx) error {
			if err := p.updateFirstName(
				ctx,
				tx,
				nameInfo.EmployeeID,
				nameInfo.Name,
				nameInfo.OperationTime,
			); err != nil {
				return fmt.Errorf("update first name: %w", err)
			}

			if err := p.completeAction(
				ctx,
				tx,
				nameInfo.TriggeredActionID,
				nameInfo.OperationTime,
			); err != nil {
				return fmt.Errorf("complete action: %w", err)
			}

			if nextAction != nil {
				if err := p.addAction(ctx, tx, nextAction); err != nil {
					return fmt.Errorf("add next action: %w", err)
				}
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}

func (p *Postgres) updateFirstName(
	ctx context.Context,
	e sqlx.ExtContext,
	employeeID int,
	name string,
	updatedAt time.Time,
) error {
	res, err := sqlx.NamedExecContext(ctx, e,
		`
			UPDATE employee
			SET
				first_name = :first_name,
				updated_at = :updated_at,
			WHERE
				id = :employee_id
		`,
		map[string]any{
			"first_name":  name,
			"employee_id": employeeID,
			"updated_at":  updatedAt,
		},
	)

	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

func (p *Postgres) completeAction(
	ctx context.Context,
	e sqlx.ExtContext,
	actionID int,
	completedAt time.Time,
) error {
	res, err := sqlx.NamedExecContext(ctx, e,
		`
			UPDATE actions
			SET
				is_completed = TRUE,
				completed_at = :completed_at,
			WHERE
				id = :action_id AND
				is_completed = FALSE
		`,
		map[string]any{
			"completed_at": completedAt,
			"action_id":    actionID,
		},
	)

	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("no rows updated")
	}

	return nil
}
