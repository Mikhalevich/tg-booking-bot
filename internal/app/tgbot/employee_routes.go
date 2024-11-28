package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type Employer interface {
	CreateEmployee(ctx context.Context, info port.MessageInfo) error
	CreateOwnerIfNotExists(ctx context.Context, info port.MessageInfo) error
	GetAllEmployee(ctx context.Context, info port.MessageInfo) error
	EmployeeMiddleware(next port.Handler) port.Handler
}

func EmployeeRoutes(e Employer) RouteRegisterFunc {
	return func(register Register) {
		register.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)

		register.AddMiddleware(e.EmployeeMiddleware)

		register.AddExactTextRoute("/createemployee", e.CreateEmployee)
		register.AddExactTextRoute("/getallemployee", e.GetAllEmployee)
	}
}
