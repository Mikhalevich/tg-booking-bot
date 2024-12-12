package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) ProcessCallbackQuery(ctx context.Context, msgInfo msginfo.Info) error {
	actionData, err := button.DecodeFromString(msgInfo.Data)
	if err != nil {
		return fmt.Errorf("decode callback query data: %w", err)
	}

	if actionData.Type == button.ActionTypeCancel {
		if err := e.cancelAction(ctx, actionData.ID, msgInfo); err != nil {
			return fmt.Errorf("cancel action: %w", err)
		}
	}

	return nil
}
