package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

type employee struct {
	repository port.EmployeeRepository
	sender     port.MessageSender
}

func New(
	repository port.EmployeeRepository,
	sender port.MessageSender,
) *employee {
	return &employee{
		repository: repository,
		sender:     sender,
	}
}

func (e *employee) nextNotCompletedAction(
	ctx context.Context,
	employeeID empl.EmployeeID,
) (*action.ActionInfo, bool, error) {
	info, err := e.repository.GetNextInProgressAction(ctx, employeeID)
	if err != nil {
		if e.repository.IsNotFoundError(err) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("get next action from repo: %w", err)
	}

	return info, true, nil
}
