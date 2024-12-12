package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (e *employee) actionEditFirstName(
	ctx context.Context,
	msgInfo port.MessageInfo,
	actionInfo action.ActionInfo,
) error {
	emplID, err := actionpayload.EmployeeIDFromBytes(actionInfo.Payload)
	if err != nil {
		return fmt.Errorf("convert payload to employee_id: %w", err)
	}

	var actionID action.ActionID

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			now := time.Now()

			if err := tx.UpdateFirstName(ctx, emplID.ID, msgInfo.Text, now); err != nil {
				return fmt.Errorf("update first name: %w", err)
			}

			if err := tx.CompleteAction(ctx, actionInfo.ActionID, now); err != nil {
				return fmt.Errorf("complete action: %w", err)
			}

			actionID, err = tx.AddAction(ctx, &action.ActionInfo{
				EmployeeID: actionInfo.EmployeeID,
				Action:     action.EditEmployeeLastName,
				Payload:    actionInfo.Payload,
				CreatedAt:  time.Now(),
			})

			return nil
		},
	); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	e.sender.ReplyText(ctx, msgInfo.ChatID, msgInfo.MessageID,
		"First name edited successfully. Please enter your first name",
		button.CancelButton("Cancel", button.ActionData{
			ID:   actionID,
			Type: button.ActionTypeCancel,
		}),
	)

	return nil
}
