package tgbot

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type Scheduler interface {
	GetAllTemplates(ctx context.Context, info port.MessageInfo) error
	CreateTestScheduleTemplate(ctx context.Context, info port.MessageInfo) error
}

func ScheduleRoutes(scheduler Scheduler) RouteRegisterFunc {
	return func(register Register) {
		register.AddExactTextRoute("/getalltemplates", scheduler.GetAllTemplates)
		register.AddExactTextRoute("/createtesttemplate", scheduler.CreateTestScheduleTemplate)
	}
}
