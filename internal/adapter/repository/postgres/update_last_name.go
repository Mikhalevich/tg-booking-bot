package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

func (p *Postgres) UpdateLastName(
	ctx context.Context,
	employeeID empl.EmployeeID,
	name string,
	updatedAt time.Time,
) error {
	res, err := sqlx.NamedExecContext(ctx, p.db,
		`
			UPDATE employee
			SET
				last_name = :last_name,
				updated_at = :updated_at
			WHERE
				id = :employee_id
		`,
		map[string]any{
			"last_name":   name,
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
		return errNotUpdated
	}

	return nil
}
