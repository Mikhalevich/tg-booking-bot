package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (e *employee) EmployeeMiddleware(next port.Handler) port.Handler {
	return func(ctx context.Context, info port.MessageInfo) error {
		empl, err := e.repository.GetEmployeeByChatID(ctx, info.ChatID)

		if err != nil {
			if !e.repository.IsEmployeeNotFoundError(err) {
				return fmt.Errorf("get employee by chat_id: %w", err)
			}

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

		if err := next(ctxdata.WithEmployee(ctx, empl), info); err != nil {
			return fmt.Errorf("employee middleware: %w", err)
		}

		return nil
	}
}
