package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/adapter/repository/postgres/internal/model"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
)

func (p *Postgres) GetAllEmployee(ctx context.Context) ([]empl.Employee, error) {
	var empls []model.Employee
	if err := sqlx.SelectContext(
		ctx,
		p.db,
		&empls,
		`SELECT
			employee.id,
			employee.first_name,
			employee.last_name,
			employee.role_id,
			role.name AS role_name,
			employee.chat_id,
			employee.state,
			employee.verification_code,
			employee.created_at,
			employee.updated_at
		FROM
			employee INNER JOIN role on employee.role_id = role.id
	`); err != nil {
		return nil, fmt.Errorf("select employee: %w", err)
	}

	return model.ToPortEmployees(empls), nil
}
