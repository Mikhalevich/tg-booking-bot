package employee

import (
	"context"
	"errors"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/internal/ctxdata"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

func (e *employee) ProcessNextAction(ctx context.Context, msgInfo port.MessageInfo) error {
	empl, ok := ctxdata.Employee(ctx)
	if !ok {
		if err := e.processVerificationCode(ctx, msgInfo.ChatID, msgInfo.Text); err != nil {
			return fmt.Errorf("process verification code: %w", err)
		}

		return nil
	}

	actionInfo, ok, err := e.nextNotCompletedAction(ctx, empl.ID)
	if err != nil {
		return fmt.Errorf("next not completed action: %w", err)
	}

	if !ok {
		if err := e.sender.ReplyText(
			ctx, msgInfo.ChatID,
			msgInfo.MessageID,
			"no action required",
		); err != nil {
			return fmt.Errorf("reply text no action required: %w", err)
		}
	}

	if err := e.processAction(ctx, msgInfo, actionInfo); err != nil {
		return fmt.Errorf("process next action: %w", err)
	}

	return nil
}

func (e *employee) processVerificationCode(ctx context.Context, chatID int64, code string) error {
	return nil
}

func (e *employee) processAction(
	ctx context.Context,
	msgInfo port.MessageInfo,
	actionInfo action.ActionInfo,
) error {
	var err error

	switch actionInfo.Action {
	case action.EditEmployeeFirstName:
		err = e.actionEditFirstName(ctx, msgInfo, actionInfo)
	case action.EditEmployeeLastName:
		err = errors.New("not implemented")
	default:
		err = errors.New("invalid action")
	}

	if err != nil {
		return fmt.Errorf("action %s: %w", actionInfo.Action.String(), err)
	}

	return nil
}
