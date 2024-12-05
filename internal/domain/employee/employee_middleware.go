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
			if !e.repository.IsNotFoundError(err) {
				return fmt.Errorf("get employee by chat_id: %w", err)
			}
		} else {
			ctx = ctxdata.WithEmployee(ctx, empl)
		}

		return next(ctx, info)
	}
}
