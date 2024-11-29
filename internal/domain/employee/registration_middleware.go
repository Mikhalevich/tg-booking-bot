package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (e *employee) RegistrationMiddleware(next port.Handler) port.Handler {
	return func(ctx context.Context, info port.MessageInfo) error {
		_, ok := ctxdata.Employee(ctx)

		if !ok {
			if err := e.sender.ReplyText(
				ctx,
				info.ChatID,
				info.MessageID,
				"please enter verification code for registration",
			); err != nil {
				return fmt.Errorf("reply text: %w", err)
			}

			return nil
		}

		return next(ctx, info)
	}
}
