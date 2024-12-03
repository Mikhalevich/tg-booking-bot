package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot/router"
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
	return func(r router.Register) {
		r.MiddlewareGroup(func(r router.Register) {
			r.AddMiddleware(e.EmployeeMiddleware)

			r.MiddlewareGroup(func(r router.Register) {
				r.AddMiddleware(e.RegistrationMiddleware)

				r.AddExactTextRoute("/createemployee", e.CreateEmployee)
				r.AddExactTextRoute("/getallemployee", e.GetAllEmployee)
			})

			r.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)

			r.AddDefaultTextHandler(e.ProcessNextAction)
		})
	}
}
