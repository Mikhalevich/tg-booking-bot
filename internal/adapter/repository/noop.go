package repository

import (
	"context"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
)

type noop struct {
}

func NewNoop() *noop {
	return &noop{}
}

func (n *noop) GetAllTemplates(ctx context.Context) ([]port.Schedule, error) {
	return nil, nil
}
