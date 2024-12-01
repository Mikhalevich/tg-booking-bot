package employee

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/employee/internal/actionpayload"
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

	if err := e.repository.EditFirstName(
		ctx,
		port.EditNameInput{
			EmployeeID:        emplID.ID,
			Name:              msgInfo.Text,
			TriggeredActionID: actionInfo.ActionID,
		},
		nil); err != nil {
		return fmt.Errorf("edit first name: %w", err)
	}

	return nil
}
