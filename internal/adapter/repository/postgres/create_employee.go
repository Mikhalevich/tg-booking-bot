package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) CreateEmployee(ctx context.Context, r role.Role, verificationCode string) error {
	var roleID int
	if err := p.db.GetContext(
		ctx,
		&roleID,
		`SELECT id
		FROM role
		WHERE name = $1
	`, r.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("role name %s doesn't exist", r.String())
		}

		return fmt.Errorf("get role_id by role name")
	}

	res, err := p.db.NamedExecContext(
		ctx,
		`INSERT INTO employee(
			role_id,
			state,
			verification_code
		) VALUES (
			:role_id,
			:state,
			:verification_code
		)`, map[string]any{
			"role_id":           roleID,
			"state":             port.EmployeeStateVerificationRequired,
			"verification_code": verificationCode,
		})

	if err != nil {
		return fmt.Errorf("exec create employee: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected: %w", err)
	}

	return nil
}
