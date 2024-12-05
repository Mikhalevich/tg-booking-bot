package postgres

import (
	"context"
	"database/sql"
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

func toSQLIsolationLevel(level port.TransactionLevel) sql.IsolationLevel {
	switch level {
	case port.TransactionLevelDefault:
		return sql.LevelDefault
	case port.TransactionLevelSerializable:
		return sql.LevelSerializable
	}

	return sql.LevelDefault
}

func (p *Postgres) Transaction(
	ctx context.Context,
	level port.TransactionLevel,
	fn func(ctx context.Context, tx port.EmployeeRepository) error,
) error {
	if err := transaction.TransactionWithLevel(
		ctx,
		p.db,
		toSQLIsolationLevel(level),
		func(ctx context.Context, tx *sqlx.Tx) error {
			return fn(ctx, newFromTransaction(tx))
		},
	); err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	return nil
}
