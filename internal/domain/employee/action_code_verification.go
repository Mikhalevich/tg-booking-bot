package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) actionCodeVerification(
	ctx context.Context,
	chatID msginfo.ChatID,
	msgID msginfo.MessageID,
	code string,
) error {
	var (
		editFirstNameActionID action.ActionID
		verifiedEmployee      *empl.Employee
		err                   error
	)

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			verifiedEmployee, err = tx.CodeVerification(ctx, code, chatID)
			if err != nil {
				return fmt.Errorf("code verification: %w", err)
			}

			editFirstNameActionPayload, err := actionpayload.BytesFromEmployeeID(verifiedEmployee.ID)
			if err != nil {
				return fmt.Errorf("convert employee id to action payload: %w", err)
			}

			editFirstNameActionID, err = tx.AddAction(ctx, &action.ActionInfo{
				EmployeeID: verifiedEmployee.ID.Int(),
				Action:     action.EditEmployeeFirstName,
				Payload:    editFirstNameActionPayload,
				CreatedAt:  time.Now(),
			})

			if err != nil {
				return fmt.Errorf("edit first name action: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("code verification transaction: %w", err)
	}

	e.sender.ReplyText(ctx, chatID, msgID, "Verification completed successfully")

	e.sender.ReplyText(ctx, chatID, msgID,
		fmt.Sprintf("Edit your first name: %s", verifiedEmployee.FirstName),
		button.CancelButton("Cancel", button.ActionData{
			ID:   editFirstNameActionID,
			Type: button.ActionTypeCancel,
		}),
	)

	return nil
}
