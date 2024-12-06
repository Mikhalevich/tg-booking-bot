package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (e *employee) cancelAction(ctx context.Context, actionID int, msgInfo port.MessageInfo) error {
	if err := e.repository.CancelAction(ctx, actionID, time.Now()); err != nil {
		if !e.repository.IsNotUpdatedError(err) {
			return fmt.Errorf("cancel action: %w", err)
		}

		if err := e.sender.ReplyText(
			ctx,
			msgInfo.ChatID,
			msgInfo.MessageID,
			"Already processed or expired",
		); err != nil {
			return fmt.Errorf("reply text: %w", err)
		}

		return nil
	}

	if err := e.sender.ReplyText(
		ctx,
		msgInfo.ChatID,
		msgInfo.MessageID,
		"Canceled successfully",
	); err != nil {
		return fmt.Errorf("reply text: %w", err)
	}

	return nil
}
