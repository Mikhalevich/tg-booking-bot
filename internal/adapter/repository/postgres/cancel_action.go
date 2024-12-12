package postgres

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (p *Postgres) CancelAction(ctx context.Context, id action.ActionID, completedAt time.Time) error {
	return p.completeAction(ctx, id, action.StateCanceled, completedAt)
}
