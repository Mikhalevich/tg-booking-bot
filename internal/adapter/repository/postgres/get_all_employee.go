package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

type employee struct {
	ID               int                `db:"id"`
	FirstName        string             `db:"first_name"`
	LastName         string             `db:"last_name"`
	Role             role.Role          `db:"role_name"`
	ChatID           sql.NullInt64      `db:"chat_id"`
	State            empl.EmployeeState `db:"state"`
	VerificationCode sql.NullString     `db:"verification_code"`
	CreatedAt        time.Time          `db:"created_at"`
	UpdatedAt        time.Time          `db:"updated_at"`
}

func (p *Postgres) GetAllEmployee(ctx context.Context) ([]empl.Employee, error) {
	var empls []employee
	if err := sqlx.SelectContext(
		ctx,
		p.db,
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

func convertToEmployee(e employee) empl.Employee {
	return empl.Employee{
		ID:               empl.EmployeeIDFromInt(e.ID),
		FirstName:        e.FirstName,
		LastName:         e.LastName,
		Role:             e.Role,
		ChatID:           msginfo.ChatIDFromInt(e.ChatID.Int64),
		State:            e.State,
		VerificationCode: e.VerificationCode.String,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func convertToEmployees(empls []employee) []empl.Employee {
	if len(empls) == 0 {
		return nil
	}

	portEmpls := make([]empl.Employee, 0, len(empls))
	for _, e := range empls {
		portEmpls = append(portEmpls, convertToEmployee(e))
	}

	return portEmpls
}
