package model

import (
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/action"
)

type ActionInfo struct {
	ActionID   int           `db:"id"`
	EmployeeID int           `db:"employee_id"`
	Action     action.Action `db:"action"`
	Payload    []byte        `db:"payload"`
	State      action.State  `db:"state"`
}

func ToPortActionInfo(a ActionInfo) action.ActionInfo {
	return action.ActionInfo{
		ActionID:   action.ActionIDFromInt(a.ActionID),
		EmployeeID: a.EmployeeID,
		Action:     a.Action,
		Payload:    a.Payload,
		State:      a.State,
	}
}
