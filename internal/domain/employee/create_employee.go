package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (e *employee) CreateEmployee(ctx context.Context, info port.MessageInfo) error {
	verificationCode := uuid.NewString()

	createdEmpID, err := e.repository.CreateEmployee(ctx, role.Employee, verificationCode)
	if err != nil {
		return fmt.Errorf("create employee: %w", err)
	}

	if err := e.sender.ReplyTextMarkdown(
		ctx,
		info.ChatID,
		info.MessageID,
		fmt.Sprintf("sent this code to user to complete registration: _%s_", e.sender.EscapeMarkdown(verificationCode)),
	); err != nil {
		return fmt.Errorf("reply text markdown: %w", err)
	}

	currentEmployee, err := e.repository.GetEmployeeByChatID(ctx, info.ChatID)
	if err != nil {
		return fmt.Errorf("get employee by chat_id: %w", err)
	}

	actionPayload, err := actionpayload.BytesFromEmployeeID(createdEmpID)
	if err != nil {
		return fmt.Errorf("convert employee id to action payload: %w", err)
	}

	if _, err := e.repository.AddAction(ctx, &action.ActionInfo{
		EmployeeID: currentEmployee.ID,
		Action:     action.EditEmployeeFirstName,
		Payload:    actionPayload,
		CreatedAt:  time.Now(),
	}); err != nil {
		return fmt.Errorf("edit employee first name action: %w", err)
	}

	return nil
}
