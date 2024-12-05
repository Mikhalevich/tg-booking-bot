package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func (p *Postgres) UpdateFirstName(
	ctx context.Context,
	employeeID int,
	name string,
	updatedAt time.Time,
) error {
	res, err := sqlx.NamedExecContext(ctx, p.db,
		`
			UPDATE employee
			SET
				first_name = :first_name,
				updated_at = :updated_at
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
		return errNotUpdated
	}

	return nil
}
