package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

// CreateOwnerIfNotExists returns id of created owner.
func (p *Postgres) CreateOwnerIfNotExists(
	ctx context.Context,
	chatID int64,
) (int, error) {
	employeeID, err := p.createEmployeeWithoutVerification(ctx, role.Unspecified, chatID)
	if err != nil {
		return 0, fmt.Errorf("create employee with invalid role: %w", err)
	}

	if err := p.updateRoleFromUnspecifiedToOwner(ctx, employeeID); err != nil {
		return 0, fmt.Errorf("upate role from unspecified to owner: %w", err)
	}

	return employeeID, nil
}

func (p *Postgres) createEmployeeWithoutVerification(
	ctx context.Context,
	r role.Role,
	chatID int64,
) (int, error) {
	roleID, err := p.roleIDByName(ctx, r)
	if err != nil {
		return 0, fmt.Errorf("get role id: %w", err)
	}

	query, args, err := sqlx.Named(
		`INSERT INTO employee(
				role_id,
				state,
				chat_id
			) VALUES (
				:role_id,
				:state,
				:chat_id
			)
			RETURNING id
			`, map[string]any{
			"role_id": roleID,
			"state":   port.EmployeeStateRegistered,
			"chat_id": chatID,
		})

	if err != nil {
		return 0, fmt.Errorf("create named query: %w", err)
	}

	var employeeID int
	if err := sqlx.GetContext(ctx, p.db, &employeeID, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("insert employee: %w", err)
	}

	return employeeID, nil
}

func (p *Postgres) updateRoleFromUnspecifiedToOwner(ctx context.Context, employeeID int) error {
	ownerID, err := p.roleIDByName(ctx, role.Owner)
	if err != nil {
		return fmt.Errorf("get owner role id: %w", err)
	}

	res, err := p.db.NamedExecContext(ctx, `
		UPDATE employee SET
			role_id = :role_id
		WHERE
			id = :employee_id AND
			NOT EXISTS (
				SELECT *
				FROM employee
				WHERE role_id = :role_id
			)
	`, map[string]any{
		"role_id":     ownerID,
		"employee_id": employeeID,
	})

	if err != nil {
		return fmt.Errorf("named exec: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("unable to update invalid to owner")
	}

	return nil
}
