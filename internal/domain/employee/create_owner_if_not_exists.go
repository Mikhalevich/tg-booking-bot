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
)

func (e *employee) CreateOwnerIfNotExists(ctx context.Context, info port.MessageInfo) error {
	_, ok := ctxdata.Employee(ctx)
	if ok {
		if err := e.sender.ReplyText(
			ctx,
			info.ChatID,
			info.MessageID,
			"not allowed",
		); err != nil {
			return fmt.Errorf("not allowed reply: %w", err)
		}

		return nil
	}

	if err := e.repository.Transaction(
		ctx,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			if err := createOwnerIfNotExistsTx(ctx, tx, info.ChatID); err != nil {
				return fmt.Errorf("create owner tx: %w", err)
			}

			return nil
		},
	); err != nil {
		if !e.repository.IsOwnerAlreadyExists(err) {
			return fmt.Errorf("transaction: %w", err)
		}

		if err := e.sender.ReplyText(
			ctx,
			info.ChatID,
			info.MessageID,
			"owner already exists",
		); err != nil {
			return fmt.Errorf("owner exists reply: %w", err)
		}

		return nil
	}

	if err := e.sender.ReplyText(
		ctx,
		info.ChatID,
		info.MessageID,
		"Owner created. Please specify first name",
		button.CancelButton("Cancel", button.ActionData{
			ID:   123,
			Type: button.ActionTypeCancel,
		}),
	); err != nil {
		return fmt.Errorf("created reply: %w", err)
	}

	return nil
}

func createOwnerIfNotExistsTx(
	ctx context.Context,
	tx port.EmployeeRepository,
	chatID int64,
) error {
	createdOwnerID, err := tx.CreateOwnerIfNotExists(ctx, chatID)
	if err != nil {
		return fmt.Errorf("create owner: %w", err)
	}

	actionPayload, err := actionpayload.BytesFromEmployeeID(createdOwnerID)
	if err != nil {
		return fmt.Errorf("convert employee id to action payload: %w", err)
	}

	if err := tx.AddAction(ctx, &action.ActionInfo{
		EmployeeID: createdOwnerID,
		Action:     action.EditEmployeeFirstName,
		Payload:    actionPayload,
		CreatedAt:  time.Now(),
	}); err != nil {
		return fmt.Errorf("edit employee first name action: %w", err)
	}

	return nil
}
