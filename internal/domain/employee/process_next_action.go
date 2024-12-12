package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) ProcessNextAction(ctx context.Context, msgInfo msginfo.Info) error {
	empl, ok := ctxdata.Employee(ctx)
	if !ok {
		if err := e.actionCodeVerification(ctx, msgInfo.ChatID, msgInfo.MessageID, msgInfo.Text); err != nil {
			return fmt.Errorf("process verification code: %w", err)
		}

		return nil
	}

	actionInfo, ok, err := e.nextNotCompletedAction(ctx, empl.ID)
	if err != nil {
		return fmt.Errorf("next not completed action: %w", err)
	}

	if !ok {
		e.sender.ReplyText(ctx, msgInfo.ChatID, msgInfo.MessageID, "no action required")
		return nil
	}

	if err := e.processAction(ctx, msgInfo, *actionInfo); err != nil {
		return fmt.Errorf("process next action: %w", err)
	}

	return nil
}

func (e *employee) processAction(
	ctx context.Context,
	msgInfo msginfo.Info,
	actionInfo action.ActionInfo,
) error {
	var err error

	switch actionInfo.Action {
	case action.EditEmployeeFirstName:
		err = e.actionEditFirstName(ctx, msgInfo, actionInfo)
	case action.EditEmployeeLastName:
		err = e.actionEditLastName(ctx, msgInfo, actionInfo)
	default:
		err = errors.New("invalid action")
	}

	if err != nil {
		return fmt.Errorf("action %s: %w", actionInfo.Action.String(), err)
	}

	return nil
}
