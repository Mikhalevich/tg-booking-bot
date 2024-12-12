package employee

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) RegistrationMiddleware(next msginfo.Handler) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		_, ok := ctxdata.Employee(ctx)

		if !ok {
			e.sender.ReplyText(ctx, info.ChatID, info.MessageID,
				"please enter verification code for registration")
			return nil
		}

		return next(ctx, info)
	}
}
