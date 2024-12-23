package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) actionEditLastName(
	ctx context.Context,
	msgInfo msginfo.Info,
	actionInfo action.ActionInfo,
) error {
	emplID, err := actionpayload.EmployeeIDFromBytes(actionInfo.Payload)
	if err != nil {
		return fmt.Errorf("convert payload to employee_id: %w", err)
	}

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			now := time.Now()

			if err := tx.UpdateLastName(ctx, emplID.ID, msgInfo.Text, now); err != nil {
				return fmt.Errorf("update last name: %w", err)
			}

			if err := tx.CompleteAction(ctx, actionInfo.ActionID, now); err != nil {
				return fmt.Errorf("complete action: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	e.sender.ReplyText(ctx, msgInfo.ChatID, msgInfo.MessageID,
		"Last name edited successfully")

	return nil
}
