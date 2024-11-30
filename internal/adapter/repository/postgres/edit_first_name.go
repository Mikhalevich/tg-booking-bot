package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) EditFirstName(
	ctx context.Context,
	nameInfo port.EditNameInput,
	nextAction *action.ActionInfo,
) error {
	if err := transaction.Transaction(ctx, p.db,
		func(ctx context.Context, tx *sqlx.Tx) error {
			return nil
		}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
