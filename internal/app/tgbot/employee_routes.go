package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

type Employer interface {
	CreateEmployee(ctx context.Context, msgInfo msginfo.Info) error
	CreateOwnerIfNotExists(ctx context.Context, msgInfo msginfo.Info) error
	GetAllEmployee(ctx context.Context, msgInfo msginfo.Info) error
	ProcessNextAction(ctx context.Context, msgInfo msginfo.Info) error
	ProcessCallbackQuery(ctx context.Context, msgInfo msginfo.Info) error
	EmployeeMiddleware(next msginfo.Handler) msginfo.Handler
	RegistrationMiddleware(next msginfo.Handler) msginfo.Handler
	NotCompletedActionMiddleware(next msginfo.Handler) msginfo.Handler
}

func EmployeeRoutes(e Employer) RouteRegisterFunc {
	return func(r router.Register) {
		r.MiddlewareGroup(func(r router.Register) {
			r.AddMiddleware(e.EmployeeMiddleware)

			r.MiddlewareGroup(func(r router.Register) {
				r.AddMiddleware(e.RegistrationMiddleware)
				r.AddMiddleware(e.NotCompletedActionMiddleware)

				r.AddExactTextRoute("/createemployee", e.CreateEmployee)
				r.AddExactTextRoute("/getallemployee", e.GetAllEmployee)
			})

			r.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)

			r.AddDefaultTextHandler(e.ProcessNextAction)
			r.AddDefaultCallbackQueryHander(e.ProcessCallbackQuery)
		})
	}
}
