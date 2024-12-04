package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(ctx context.Context, tx *sqlx.Tx) error

func Transaction(
	ctx context.Context,
	s sqlx.ExtContext,
	fn TxFunc,
) error {
	return TransactionWithLevel(ctx, s, sql.LevelDefault, fn)
}

func TransactionWithLevel(
	ctx context.Context,
	s sqlx.ExtContext,
	level sql.IsolationLevel,
	fn TxFunc,
) error {
	db, ok := s.(*sqlx.DB)
	if !ok {
		return errors.New("not sqlx db object")
	}

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: level,
	})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		//nolint:errcheck
		tx.Rollback()
	}()

	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(fmt.Errorf("tx body: %w", err), fmt.Errorf("rollback: %w", err))
		}

		return fmt.Errorf("tx body: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
