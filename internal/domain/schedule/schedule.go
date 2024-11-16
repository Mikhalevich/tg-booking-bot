package schedule

import (
	"context"
	"fmt"

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

func (s *schedule) GetAllTemplates(ctx context.Context, info port.MessageInfo) error {
	_, err := s.repository.GetAllTemplates(ctx)
	if err != nil {
		return fmt.Errorf("get all templates: %w", err)
	}

	if err := s.sender.Reply(
		ctx,
		info.ChatID,
		info.MessageID,
		"no schedule",
	); err != nil {
		return fmt.Errorf("reply to message: %w", err)
	}

	return nil
}
