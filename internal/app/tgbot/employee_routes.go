package tgbot

import (
	"context"
	"fmt"

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

			r.AddExactTextRoute("/start", makeStartHandler(r))

			r.AddExactTextRoute("/makemeowner", e.CreateOwnerIfNotExists)

			r.AddDefaultTextHandler(e.ProcessNextAction)
			r.AddDefaultCallbackQueryHander(e.ProcessCallbackQuery)
		})
	}
}

func makeStartHandler(r router.Register) msginfo.Handler {
	return func(ctx context.Context, info msginfo.Info) error {
		if err := r.SendMessage(
			ctx,
			info.ChatID.Int64(),
			"To start using this bot, you should register by enter verification code",
		); err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return nil
	}
}
