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

func (e EmployeeState) String() string {
	return string(e)
}

type EmployeeID int

func (e EmployeeID) Int() int {
	return int(e)
}

func EmployeeIDFromInt(id int) EmployeeID {
	return EmployeeID(id)
}

type Employee struct {
	ID               EmployeeID
	FirstName        string
	LastName         string
	Role             role.Role
	ChatID           ChatID
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
	CreateOwnerIfNotExists(ctx context.Context, chatID ChatID) (EmployeeID, error)
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) (EmployeeID, error)
	UpdateFirstName(ctx context.Context, id EmployeeID, name string, updatedAt time.Time) error
	UpdateLastName(ctx context.Context, id EmployeeID, name string, updatedAt time.Time) error
	GetAllEmployee(ctx context.Context) ([]Employee, error)
	GetEmployeeByChatID(ctx context.Context, chatID ChatID) (Employee, error)
	AddAction(ctx context.Context, info *action.ActionInfo) (action.ActionID, error)
	GetNextInProgressAction(ctx context.Context, employeeID EmployeeID) (action.ActionInfo, error)
	CodeVerification(ctx context.Context, code string, chatID ChatID) (*Employee, error)
	CompleteAction(ctx context.Context, id action.ActionID, completedAt time.Time) error
	CancelAction(ctx context.Context, id action.ActionID, completedAt time.Time) error
}
