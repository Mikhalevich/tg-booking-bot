package employee

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/button"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

var (
	errOwnerAlreadyExists = errors.New("owner already exists")
)

func (e *employee) CreateOwnerIfNotExists(ctx context.Context, info port.MessageInfo) error {
	_, ok := ctxdata.Employee(ctx)
	if ok {
		e.sender.ReplyText(ctx, info.ChatID, info.MessageID, "not allowed")
		return nil
	}

	var (
		actionID int
		err      error
	)

	if err := e.repository.Transaction(
		ctx,
		port.TransactionLevelSerializable,
		func(ctx context.Context, tx port.EmployeeRepository) error {
			actionID, err = createOwnerIfNotExistsTx(ctx, tx, info.ChatID)
			if err != nil {
				return fmt.Errorf("create owner tx: %w", err)
			}

			return nil
		},
	); err != nil {
		if !errors.Is(err, errOwnerAlreadyExists) {
			return fmt.Errorf("transaction: %w", err)
		}

		e.sender.ReplyText(ctx, info.ChatID, info.MessageID, "owner already exists")

		return nil
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

func createOwnerIfNotExistsTx(
	ctx context.Context,
	tx port.EmployeeRepository,
	chatID int64,
) (int, error) {
	isExists, err := tx.IsEmployeeWithRoleExists(ctx, role.Owner)
	if err != nil {
		return 0, fmt.Errorf("is owner exists: %w", err)
	}

	if isExists {
		return 0, errOwnerAlreadyExists
	}

	createdOwnerID, err := tx.CreateEmployeeWithoutVerification(ctx, role.Owner, chatID)
	if err != nil {
		return 0, fmt.Errorf("create owner: %w", err)
	}

	actionPayload, err := actionpayload.BytesFromEmployeeID(createdOwnerID)
	if err != nil {
		return 0, fmt.Errorf("convert employee id to action payload: %w", err)
	}

	actionID, err := tx.AddAction(ctx, &action.ActionInfo{
		EmployeeID: createdOwnerID,
		Action:     action.EditEmployeeFirstName,
		Payload:    actionPayload,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		return 0, fmt.Errorf("edit employee first name action: %w", err)
	}

	return actionID, nil
}
