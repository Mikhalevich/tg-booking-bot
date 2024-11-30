package action

import (
	"time"
)

type Action string

func (a Action) String() string {
	return string(a)
}

type ActionInfo struct {
	ActionID   int
	EmployeeID int
	Action     Action
	Payload    []byte
	CreatedAt  time.Time
}

type CompleteActionInfo struct {
	ActionID    int
	CompletedAt time.Time
}

const (
	EditEmployeeFirstName Action = "edit_employee_first_name"
	EditEmployeeLastName  Action = "edit_employee_last_name"
)
