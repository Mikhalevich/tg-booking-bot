package empl

import (
	"time"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
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
	ChatID           msginfo.ChatID
	State            EmployeeState
	VerificationCode string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
