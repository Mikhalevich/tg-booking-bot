package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (p *Postgres) CreateOwnerIfNotExists(ctx context.Context, chatID msginfo.ChatID) (empl.EmployeeID, error) {
	var ownerID empl.EmployeeID

	if err := transaction.Transaction(ctx, p.db, true, func(ctx context.Context, tx sqlx.ExtContext) error {
		if _, err := tx.ExecContext(ctx, `LOCK TABLE employee IN SHARE ROW EXCLUSIVE MODE`); err != nil {
			return fmt.Errorf("lock table: %w", err)
		}

		isExists, err := p.isOwnerExists(ctx)
		if err != nil {
			return fmt.Errorf("is owner exists: %w", err)
		}

		if isExists {
			return errAlreadyExists
		}

		ownerID, err = p.CreateEmployeeWithoutVerification(ctx, role.Owner, chatID)
		if err != nil {
			return fmt.Errorf("create employee without verifiation")
		}

		return nil
	},
	); err != nil {
		return 0, fmt.Errorf("transaction: %w", err)
	}

	return ownerID, nil
}

func (p *Postgres) isOwnerExists(ctx context.Context) (bool, error) {
	var isExists bool
	if err := sqlx.GetContext(
		ctx,
		p.db,
		&isExists,
		`
			SELECT
				TRUE
			FROM
				employee INNER JOIN role ON employee.role_id = role.id
			WHERE
				role.name = $1
		`,
		role.Owner.String(),
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("get context: %w", err)
	}

	return true, nil
}
