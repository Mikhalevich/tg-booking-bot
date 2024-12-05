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

type EditNameInput struct {
	EmployeeID        int
	Name              string
	TriggeredActionID int
	OperationTime     time.Time
}

type TransactionLevel int

const (
	TransactionLevelDefault TransactionLevel = iota + 1
	TransactionLevelSerializable
)

//nolint:interfacebloat
type EmployeeRepository interface {
	Transaction(ctx context.Context, level TransactionLevel,
		fn func(ctx context.Context, tx EmployeeRepository) error) error
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) (int, error)
	CreateEmployeeWithoutVerification(ctx context.Context, r role.Role, chatID int64) (int, error)
	IsEmployeeWithRoleExists(ctx context.Context, role role.Role) (bool, error)
	EditFirstName(ctx context.Context, nameInfo EditNameInput, nextAction *action.ActionInfo) error
	GetAllEmployee(ctx context.Context) ([]Employee, error)
	GetEmployeeByChatID(ctx context.Context, chatID int64) (Employee, error)
	IsEmployeeNotFoundError(err error) bool
	AddAction(ctx context.Context, info *action.ActionInfo) error
	GetNextNotCompletedAction(ctx context.Context, employeeID int) (action.ActionInfo, error)
	IsActionNotFoundError(err error) bool
}
