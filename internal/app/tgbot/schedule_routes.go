package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/app/tgbot/router"
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port/msginfo"
)

type Scheduler interface {
	GetAllTemplates(ctx context.Context, info msginfo.Info) error
	CreateTestScheduleTemplate(ctx context.Context, info msginfo.Info) error
}

func ScheduleRoutes(scheduler Scheduler) RouteRegisterFunc {
	return func(register router.Register) {
		register.AddExactTextRoute("/getalltemplates", scheduler.GetAllTemplates)
		register.AddExactTextRoute("/createtesttemplate", scheduler.CreateTestScheduleTemplate)
	}
}
