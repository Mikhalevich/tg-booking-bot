package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

func (p *Postgres) GetEmployeeByChatID(ctx context.Context, chatID int64) (port.Employee, error) {
	var emp employee
	if err := p.db.GetContext(
		ctx,
		&emp,
		`SELECT
			employee.id,
			employee.first_name,
			employee.last_name,
			role.name AS role_name,
			employee.chat_id,
			employee.state,
			employee.verification_code,
			employee.created_at,
			employee.updated_at
		FROM
			employee INNER JOIN role on employee.role_id = role.id
		WHERE
			employee.chat_id = $1
	`, chatID); err != nil {
		return port.Employee{}, fmt.Errorf("select employee: %w", err)
	}

	return convertToEmployee(emp), nil
}