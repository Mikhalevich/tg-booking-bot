package noop

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

var _ port.ScheduleRepository = (*noop)(nil)

type noop struct {
}

func New() *noop {
	return &noop{}
}

func (n *noop) GetAllTemplates(ctx context.Context) ([]port.ScheduleTemplate, error) {
	return nil, nil
}

func (n *noop) CreateTestScheduleTemplate(ctx context.Context, tmpl port.ScheduleTemplate) error {
	return nil
}
