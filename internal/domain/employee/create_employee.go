package employee

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (e *employee) CreateEmployee(ctx context.Context, info port.MessageInfo) error {
	currentEmployee, ok := ctxdata.Employee(ctx)
	if !ok {
		return errors.New("unable to get current employee")
	}

	var (
		verificationCode      string
		editFirstNameActionID action.ActionID
	)

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			verificationCode = uuid.NewString()

			createdEmplID, err := tx.CreateEmployee(ctx, role.Employee, verificationCode)
			if err != nil {
				return fmt.Errorf("create employee: %w", err)
			}

			actionPayload, err := actionpayload.BytesFromEmployeeID(createdEmplID)
			if err != nil {
				return fmt.Errorf("convert employee id to action payload: %w", err)
			}

			editFirstNameActionID, err = tx.AddAction(ctx, &action.ActionInfo{
				EmployeeID: currentEmployee.ID,
				Action:     action.EditEmployeeFirstName,
				Payload:    actionPayload,
				CreatedAt:  time.Now(),
			})

			if err != nil {
				return fmt.Errorf("edit employee first name action: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("create employee transaction: %w", err)
	}

	e.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID,
		fmt.Sprintf("Sent this code to user to complete registration: `%s`", e.sender.EscapeMarkdown(verificationCode)),
	)

	e.sender.ReplyText(ctx, info.ChatID, info.MessageID,
		"Employee created. Please enter employee first name",
		button.CancelButton("Cancel", button.ActionData{
			ID:   editFirstNameActionID,
			Type: button.ActionTypeCancel,
		}),
	)

	return nil
}
