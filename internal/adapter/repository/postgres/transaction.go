package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func newFromTransaction(tx *sqlx.Tx) *Postgres {
	return &Postgres{
		db: tx,
	}
}

func (p *Postgres) Transaction(
	ctx context.Context,
	fn func(ctx context.Context, tx port.EmployeeRepository) error,
) error {
	if err := transaction.Transaction(
		ctx,
		p.db,
		func(ctx context.Context, tx *sqlx.Tx) error {
			return fn(ctx, newFromTransaction(tx))
		},
	); err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	return nil
}
