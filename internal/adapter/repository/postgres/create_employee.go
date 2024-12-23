package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) CreateEmployee(
	ctx context.Context,
	r role.Role,
	verificationCode string,
) (empl.EmployeeID, error) {
	var roleID int
	if err := sqlx.GetContext(
		ctx,
		p.db,
		&roleID,
		`SELECT id
		FROM role
		WHERE name = $1
	`, r.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("role name %s doesn't exist", r.String())
		}

		return 0, fmt.Errorf("get role_id by role name")
	}

	query, args, err := sqlx.Named(
		`INSERT INTO employee(
			role_id,
			state,
			verification_code
		) VALUES (
			:role_id,
			:state,
			:verification_code
		)
		RETURNING id
		`, map[string]any{
			"role_id":           roleID,
			"state":             empl.EmployeeStateVerificationRequired,
			"verification_code": verificationCode,
		})

	if err != nil {
		return 0, fmt.Errorf("create named query: %w", err)
	}

	var employeeID int
	if err := sqlx.GetContext(ctx, p.db, &employeeID, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("insert employee: %w", err)
	}

	return empl.EmployeeIDFromInt(employeeID), nil
}
