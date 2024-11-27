package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

type employee struct {
	ID               int                `db:"id"`
	FirstName        string             `db:"first_name"`
	LastName         string             `db:"last_name"`
	Role             role.Role          `db:"role_name"`
	ChatID           sql.NullInt64      `db:"chat_id"`
	State            port.EmployeeState `db:"state"`
	VerificationCode sql.NullString     `db:"verification_code"`
	CreatedAt        time.Time          `db:"created_at"`
	UpdatedAt        time.Time          `db:"updated_at"`
}

func (p *Postgres) GetAllEmployee(ctx context.Context) ([]port.Employee, error) {
	var empls []employee
	if err := p.db.SelectContext(
		ctx,
		&empls,
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
	`); err != nil {
		return nil, fmt.Errorf("select employee: %w", err)
	}

	return convertToEmployees(empls), nil
}

func convertToEmployee(empl employee) port.Employee {
	return port.Employee{
		ID:               empl.ID,
		FirstName:        empl.FirstName,
		LastName:         empl.LastName,
		Role:             empl.Role,
		ChatID:           empl.ChatID.Int64,
		State:            empl.State,
		VerificationCode: empl.VerificationCode.String,
		CreatedAt:        empl.CreatedAt,
		UpdatedAt:        empl.UpdatedAt,
	}
}

func convertToEmployees(empls []employee) []port.Employee {
	if len(empls) == 0 {
		return nil
	}

	portEmpls := make([]port.Employee, 0, len(empls))
	for _, e := range empls {
		portEmpls = append(portEmpls, convertToEmployee(e))
	}

	return portEmpls
}
