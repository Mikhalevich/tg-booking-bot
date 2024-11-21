package schedule

import (
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type schedule struct {
	repository port.ScheduleRepository
	sender     port.MessageSender
}

func New(
	repository port.ScheduleRepository,
	sender port.MessageSender,
) *schedule {
	return &schedule{
		repository: repository,
		sender:     sender,
	}
}
