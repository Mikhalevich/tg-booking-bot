package employee

import (
	"context"
	"fmt"
	"strings"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (e *employee) GetAllEmployee(ctx context.Context, info port.MessageInfo) error {
	empls, err := e.repository.GetAllEmployee(ctx)
	if err != nil {
		return fmt.Errorf("get all employee: %w", err)
	}

	if err := e.sender.ReplyTextMarkdown(
		ctx,
		info.ChatID,
		info.MessageID,
		e.formatEmployeeMsg(empls),
	); err != nil {
		return fmt.Errorf("reply text markdown: %w", err)
	}

	return nil
}

func (e *employee) formatEmployeeMsg(empls []port.Employee) string {
	if len(empls) == 0 {
		return "*no employes found*"
	}

	output := make([]string, 0, len(empls))

	for _, emp := range empls {
		empLine := fmt.Sprintf("*role:* %s *verification code:* %s", emp.Role, e.sender.EscapeMarkdown(emp.VerificationCode))
		output = append(output, empLine)
	}

	return strings.Join(output, "\n")
}
