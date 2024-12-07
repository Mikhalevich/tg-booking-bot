package action

import (
	"time"
)

type ActionInfo struct {
	ActionID   int
	EmployeeID int
	Action     Action
	Payload    []byte
	State      State
	CreatedAt  time.Time
}

type State string

func (s State) String() string {
	return string(s)
}

const (
	StateInProgress State = "in_progress"
	StateCompleted  State = "completed"
	StateCanceled   State = "canceled"
)

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	EditEmployeeFirstName Action = "edit_employee_first_name"
	EditEmployeeLastName  Action = "edit_employee_last_name"
)
