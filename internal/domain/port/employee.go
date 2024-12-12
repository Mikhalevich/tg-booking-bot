package port

import (
	"context"
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/empl"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/role"
)

//nolint:interfacebloat
type EmployeeRepository interface {
	Transaction(ctx context.Context, fn func(ctx context.Context, tx EmployeeRepository) error) error
	IsNotFoundError(err error) bool
	IsNotUpdatedError(err error) bool
	IsAlreadyExistsError(err error) bool
	CreateOwnerIfNotExists(ctx context.Context, chatID msginfo.ChatID) (empl.EmployeeID, error)
	CreateEmployee(ctx context.Context, r role.Role, verificationCode string) (empl.EmployeeID, error)
	UpdateFirstName(ctx context.Context, id empl.EmployeeID, name string, updatedAt time.Time) error
	UpdateLastName(ctx context.Context, id empl.EmployeeID, name string, updatedAt time.Time) error
	GetAllEmployee(ctx context.Context) ([]empl.Employee, error)
	GetEmployeeByChatID(ctx context.Context, chatID msginfo.ChatID) (*empl.Employee, error)
	AddAction(ctx context.Context, info *action.ActionInfo) (action.ActionID, error)
	GetNextInProgressAction(ctx context.Context, employeeID empl.EmployeeID) (*action.ActionInfo, error)
	CodeVerification(ctx context.Context, code string, chatID msginfo.ChatID) (*empl.Employee, error)
	CompleteAction(ctx context.Context, id action.ActionID, completedAt time.Time) error
	CancelAction(ctx context.Context, id action.ActionID, completedAt time.Time) error
}
