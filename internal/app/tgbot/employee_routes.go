package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type Employer interface {
	CreateEmployee(ctx context.Context, info port.MessageInfo) error
	CreateOwnerIfNotExists(ctx context.Context, info port.MessageInfo) error
	GetAllEmployee(ctx context.Context, info port.MessageInfo) error
}

func EmployeeRoutes(e Employer) RouteRegisterFunc {
	return func(register Register) {
		register.AddExactTextRoute("/createemployee", e.CreateEmployee)
		register.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)
		register.AddExactTextRoute("/getallemployee", e.GetAllEmployee)
	}
}
