package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) CreateOwnerIfNotExists(
	ctx context.Context,
	chatID int64,
) (port.CreateOwnerIfNotExistsOutput, error) {
	var output port.CreateOwnerIfNotExistsOutput
	if err := transaction.TransactionWithLevel(ctx, p.db, sql.LevelSerializable,
		func(ctx context.Context, tx *sqlx.Tx) error {
			ownerRoleID, err := p.roleIDByName(ctx, role.Owner, tx)
			if err != nil {
				return fmt.Errorf("get role id: %w", err)
			}

			isExists, err := p.IsEmployeeWithRoleExists(ctx, ownerRoleID, tx)
			if err != nil {
				return fmt.Errorf("check for owner existence: %w", err)
			}

			if isExists {
				output = port.CreateOwnerIfNotExistsOutput{
					IsAlreadyExists: true,
				}
				return nil
			}

			ownerID, err := p.createEmployeeWithoutVerification(ctx, ownerRoleID, chatID, tx)
			if err != nil {
				return fmt.Errorf("create owner without verification: %w", err)
			}

			output = port.CreateOwnerIfNotExistsOutput{
				CreatedOwnerID: ownerID,
			}

			return nil
		}); err != nil {
		return output, fmt.Errorf("transaction: %w", err)
	}

	return output, nil
}

func (p *Postgres) IsEmployeeWithRoleExists(
	ctx context.Context,
	roleID int,
	tx *sqlx.Tx,
) (bool, error) {
	var isExists bool
	if err := sqlx.GetContext(
		ctx,
		tx,
		&isExists,
		`
			SELECT
				TRUE
			FROM
				employee
			WHERE
				role_id = $1
		`,
		roleID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("get context: %w", err)
	}

	return true, nil
}

func (p *Postgres) createEmployeeWithoutVerification(
	ctx context.Context,
	roleID int,
	chatID int64,
	tx *sqlx.Tx,
) (int, error) {
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
	if err := sqlx.GetContext(ctx, tx, &employeeID, p.db.Rebind(query), args...); err != nil {
		return 0, fmt.Errorf("insert employee: %w", err)
	}

	return employeeID, nil
}
