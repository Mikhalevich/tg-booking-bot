package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) CreateOwnerIfNotExists(ctx context.Context, info msginfo.Info) error {
	_, ok := ctxdata.Employee(ctx)
	if ok {
		e.sender.ReplyText(ctx, info.ChatID, info.MessageID, "not allowed")
		return nil
	}

	var actionID action.ActionID

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			createdOwnerID, err := tx.CreateOwnerIfNotExists(ctx, info.ChatID)
			if err != nil {
				return fmt.Errorf("create owner: %w", err)
			}

			actionPayload, err := actionpayload.BytesFromEmployeeID(createdOwnerID)
			if err != nil {
				return fmt.Errorf("convert employee id to action payload: %w", err)
			}

			id, err := tx.AddAction(ctx, &action.ActionInfo{
				EmployeeID: createdOwnerID.Int(),
				Action:     action.EditEmployeeFirstName,
				Payload:    actionPayload,
				CreatedAt:  time.Now(),
			})

			actionID = id

			if err != nil {
				return fmt.Errorf("edit employee first name action: %w", err)
			}

			return nil
		},
	); err != nil {
		if e.repository.IsAlreadyExistsError(err) {
			e.sender.ReplyText(ctx, info.ChatID, info.MessageID, "owner already exists")
			return nil
		}

		return fmt.Errorf("transaction: %w", err)
	}

	e.sender.ReplyText(ctx, info.ChatID, info.MessageID,
		"Owner created. Please enter your first name",
		button.CancelButton("Cancel", button.ActionData{
			ID:   actionID,
			Type: button.ActionTypeCancel,
		}),
	)

	return nil
}
