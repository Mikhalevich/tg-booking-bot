package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

type EmployeeState string

const (
	EmployeeStateVerificationRequired EmployeeState = "verification_required"
	EmployeeStateRegistered           EmployeeState = "registered"
)

type Employee struct {
	ID               int
	FirstName        string
	LastName         string
	Role             role.Role
	ChatID           int64
	State            EmployeeState
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateOwnerIfNotExistsOutput struct {
	CreatedOwnerID  int
	IsAlreadyExists bool
}

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) (int, error)
	CreateOwnerIfNotExists(ctx context.Context, chatID int64) (CreateOwnerIfNotExistsOutput, error)
	GetAllEmployee(ctx context.Context) ([]Employee, error)
	GetEmployeeByChatID(ctx context.Context, chatID int64) (Employee, error)
	IsEmployeeNotFoundError(err error) bool
	AddAction(ctx context.Context, info action.AddActionInfo) error
	GetNextNotCompletedAction(ctx context.Context, employeeID int) (action.ActionInfo, error)
	IsActionNotFoundError(err error) bool
}
