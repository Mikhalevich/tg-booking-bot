package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) NotCompletedActionMiddleware(next msginfo.Handler) msginfo.Handler {
	return func(ctx context.Context, msgInfo msginfo.Info) error {
		empl, ok := ctxdata.Employee(ctx)
		if !ok {
			return next(ctx, msgInfo)
		}

		_, ok, err := e.nextNotCompletedAction(ctx, empl.ID)
		if err != nil {
			return fmt.Errorf("get next not completed action: %w", err)
		}

		if !ok {
			return next(ctx, msgInfo)
		}

		e.sender.ReplyText(ctx, msgInfo.ChatID, msgInfo.MessageID, "You have not completed action")

		return nil
	}
}
