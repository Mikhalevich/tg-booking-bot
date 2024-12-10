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

//nolint:interfacebloat
type EmployeeRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context, tx EmployeeRepository) error) error
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
	IsAlreadyExistsError(err error) bool
	CreateOwnerIfNotExists(ctx context.Context, chatID int64) (int, error)
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) (int, error)
	UpdateFirstName(ctx context.Context, id int, name string, updatedAt time.Time) error
	UpdateLastName(ctx context.Context, id int, name string, updatedAt time.Time) error
	GetAllEmployee(ctx context.Context) ([]Employee, error)
	GetEmployeeByChatID(ctx context.Context, chatID int64) (Employee, error)
	AddAction(ctx context.Context, info *action.ActionInfo) (int, error)
	GetNextInProgressAction(ctx context.Context, employeeID int) (action.ActionInfo, error)
	CompleteAction(ctx context.Context, id int, completedAt time.Time) error
	CancelAction(ctx context.Context, id int, completedAt time.Time) error
}
