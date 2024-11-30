package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (e *employee) ProcessNextAction(ctx context.Context, info port.MessageInfo) error {
	empl, ok := ctxdata.Employee(ctx)
	if !ok {
		if err := e.processVerificationCode(ctx, info.ChatID, info.Text); err != nil {
			return fmt.Errorf("process verification code: %w", err)
		}

		return nil
	}

	actionInfo, ok, err := e.nextNotCompletedAction(ctx, empl.ID)
	if err != nil {
		return fmt.Errorf("next not completed action: %w", err)
	}

	if !ok {
		if err := e.sender.ReplyText(
			ctx, info.ChatID,
			info.MessageID,
			"no action required",
		); err != nil {
			return fmt.Errorf("reply text no action required: %w", err)
		}
	}

	if err := e.processNextAction(ctx, actionInfo); err != nil {
		return fmt.Errorf("process next action: %w", err)
	}

	return nil
}

func (e *employee) processVerificationCode(ctx context.Context, chatID int64, code string) error {
	return nil
}

func (e *employee) processNextAction(ctx context.Context, a action.ActionInfo) error {
	return nil
}
