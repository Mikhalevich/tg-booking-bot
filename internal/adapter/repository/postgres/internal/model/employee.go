package model

import (
	"database/sql"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

type Employee struct {
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

func ToPortEmployee(e Employee) *empl.Employee {
	return &empl.Employee{
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

func ToPortEmployees(empls []Employee) []empl.Employee {
	if len(empls) == 0 {
		return nil
	}

	portEmpls := make([]empl.Employee, 0, len(empls))
	for _, e := range empls {
		portEmpls = append(portEmpls, *ToPortEmployee(e))
	}

	return portEmpls
}
