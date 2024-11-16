package repository

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

var _ port.ScheduleRepository = (*noop)(nil)

type noop struct {
}

func NewNoop() *noop {
	return &noop{}
}

func (n *noop) GetAllTemplates(ctx context.Context) ([]port.Schedule, error) {
	return nil, nil
}
