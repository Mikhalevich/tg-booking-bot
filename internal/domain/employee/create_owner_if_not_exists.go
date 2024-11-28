package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (e *employee) CreateOwnerIfNotExists(ctx context.Context, info port.MessageInfo) error {
	createdOwnerID, err := e.repository.CreateOwnerIfNotExists(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("create owner: %w", err)
	}

	if err := e.sender.ReplyText(
		ctx,
		info.ChatID,
		info.MessageID,
		"created",
	); err != nil {
		return fmt.Errorf("reply text: %w", err)
	}

	actionPayload, err := employeeIDToActionPayload(createdOwnerID)
	if err != nil {
		return fmt.Errorf("convert employee id to action payload: %w", err)
	}

	if err := e.repository.AddAction(ctx, action.AddActionInfo{
		EmployeeID: createdOwnerID,
		Action:     action.EditEmployeeFirstName,
		Payload:    actionPayload,
		CreatedAt:  time.Now(),
	}); err != nil {
		return fmt.Errorf("edit employee first name action: %w", err)
	}

	return nil
}
