package tgbot

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/Mikhalevich/tg-booking-bot/internal/domain/port"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/logger"
	"github.com/Mikhalevich/tg-booking-bot/internal/infra/tracing"
)

type tgbot struct {
	bot    *bot.Bot
	logger logger.Logger
}

type TextHandlerFunc func(ctx context.Context, info port.MessageInfo) error

type Register interface {
	AddExactTextRoute(pattern string, handler TextHandlerFunc)
}

type RouteRegisterFunc func(register Register)

func Start(
	ctx context.Context,
	b *bot.Bot,
	logger logger.Logger,
	routesRegisterFn RouteRegisterFunc,
) error {
	tbot := &tgbot{
		bot:    b,
		logger: logger,
	}

	routesRegisterFn(tbot)

	b.Start(ctx)

	return nil
}

func (t *tgbot) AddExactTextRoute(pattern string, handler TextHandlerFunc) {
	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		t.wrapTextHandler(pattern, handler),
	)
}

func (t *tgbot) wrapTextHandler(pattern string, handler TextHandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		if err := handler(ctx, port.MessageInfo{
			MessageID: update.Message.ID,
			ChatID:    update.Message.Chat.ID,
			Text:      update.Message.Text,
		}); err != nil {
			t.logger.WithError(err).Error("error while processing message")
		}
	}
}
