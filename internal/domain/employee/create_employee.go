package employee

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

func (e *employee) CreateEmployee(ctx context.Context, info port.MessageInfo) error {
	verificationCode := uuid.NewString()

	if err := e.repository.CreateEmployee(ctx, role.Employee, verificationCode); err != nil {
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

	return nil
}
