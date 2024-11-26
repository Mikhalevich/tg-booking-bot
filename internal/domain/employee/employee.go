package employee

import (
	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type employee struct {
	repository port.EmployeeRepository
	sender     port.MessageSender
}

func New(
	repository port.EmployeeRepository,
	sender port.MessageSender,
) *employee {
	return &employee{
		repository: repository,
		sender:     sender,
	}
}
