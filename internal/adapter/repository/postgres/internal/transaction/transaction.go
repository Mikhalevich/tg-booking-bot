package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(ctx context.Context, tx *sqlx.Tx) error

func Transaction(
	ctx context.Context,
	s *sqlx.DB,
	fn TxFunc,
) error {
	tx, err := s.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			//nolint:errcheck
			tx.Rollback()
			panic(r)
		}
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
