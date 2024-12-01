package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type Employer interface {
	CreateEmployee(ctx context.Context, msgInfo port.MessageInfo) error
	CreateOwnerIfNotExists(ctx context.Context, msgInfo port.MessageInfo) error
	GetAllEmployee(ctx context.Context, msgInfo port.MessageInfo) error
	ProcessNextAction(ctx context.Context, msgInfo port.MessageInfo) error
	EmployeeMiddleware(next port.Handler) port.Handler
	RegistrationMiddleware(next port.Handler) port.Handler
}

func EmployeeRoutes(e Employer) RouteRegisterFunc {
	return func(register Register) {
		register.AddMiddleware(e.EmployeeMiddleware)

		register.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)

		register.AddMiddleware(e.RegistrationMiddleware)

		register.AddExactTextRoute("/createemployee", e.CreateEmployee)
		register.AddExactTextRoute("/getallemployee", e.GetAllEmployee)

		register.AddDefaultTextHandler(e.ProcessNextAction, e.EmployeeMiddleware)
	}
}
