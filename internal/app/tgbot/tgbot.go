package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
)

type Scheduler interface {
	GetAllTemplates(ctx context.Context, info port.MessageInfo) error
}

type tgbot struct {
	logger    logger.Logger
	scheduler Scheduler
}

func Start(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	scheduler Scheduler,
) error {
	tbot := &tgbot{
		logger:    logger,
		scheduler: scheduler,
	}

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/getalltemplates",
		bot.MatchTypeExact,
		tbot.handlerWrapper(tbot.scheduler.GetAllTemplates),
	)

	b.Start(ctx)

	return nil
}

type processorFunc func(ctx context.Context, info port.MessageInfo) error

func (t *tgbot) handlerWrapper(processor processorFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if err := processor(ctx, port.MessageInfo{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
			Text:      update.Message.Text,
		}); err != nil {
			t.logger.WithError(err).Error("error while processing message")
		}
	}
}
