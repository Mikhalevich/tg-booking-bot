package employee

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

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

	actionPayload, err := employeeIDToActionPayload(createdEmpID)
	if err != nil {
		return fmt.Errorf("convert employee id to action payload: %w", err)
	}

	if err := e.repository.AddAction(ctx, action.AddActionInfo{
		EmployeeID: currentEmployee.ID,
		Action:     action.EditEmployeeFirstName,
		Payload:    actionPayload,
		CreatedAt:  time.Now(),
	}); err != nil {
		return fmt.Errorf("edit employee first name action: %w", err)
	}

	return nil
}

func employeeIDToActionPayload(id int) ([]byte, error) {
	b, err := json.Marshal(map[string]int{
		"employee_id": id,
	})

	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return b, nil
}

//nolint:unused
func actionPayloadToEmployeeID(payload []byte) (int, error) {
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		return 0, fmt.Errorf("json unmarshal: %w", err)
	}

	rawID, ok := m["employee_id"]
	if !ok {
		return 0, errors.New("no employee id in payload")
	}

	id, ok := rawID.(int)

	if !ok {
		return 0, errors.New("invalid employee_id type")
	}

	return id, nil
}
