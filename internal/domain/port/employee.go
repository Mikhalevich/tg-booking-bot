package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

type EmployeeState string

const (
	EmployeeStateVerificationRequired EmployeeState = "verification_required"
	EmployeeStateRegistered           EmployeeState = "registered"
)

type Employee struct {
	FirstName        string
	LastName         string
	Role             role.Role
	ChatID           int64
	State            EmployeeState
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) error
	GetAllEmployee(ctx context.Context) ([]Employee, error)
}
