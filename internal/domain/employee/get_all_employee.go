package employee

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

func (e *employee) GetAllEmployee(ctx context.Context, info msginfo.Info) error {
	empls, err := e.repository.GetAllEmployee(ctx)
	if err != nil {
		return fmt.Errorf("get all employee: %w", err)
	}

	e.sender.ReplyTextMarkdown(ctx, info.ChatID, info.MessageID, e.formatEmployeeMsg(empls))

	return nil
}

func (e *employee) formatEmployeeMsg(empls []empl.Employee) string {
	if len(empls) == 0 {
		return "*no employes found*"
	}

	output := make([]string, 0, len(empls))

	for _, emp := range empls {
		empLine := fmt.Sprintf(
			"*first name:* %s *last name:* %s *role:* %s *code:* %s *created at:* %s",
			e.sender.EscapeMarkdown(emp.FirstName),
			e.sender.EscapeMarkdown(emp.LastName),
			emp.Role,
			e.sender.EscapeMarkdown(emp.VerificationCode),
			e.sender.EscapeMarkdown(emp.CreatedAt.Format(time.RFC3339)),
		)
		output = append(output, empLine)
	}

	return strings.Join(output, "\n")
}
