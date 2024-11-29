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
	bot         *bot.Bot
	logger      logger.Logger
	middlewares []Middleware
}

type Middleware func(next port.Handler) port.Handler

type Register interface {
	AddExactTextRoute(pattern string, handler port.Handler)
	AddMiddleware(middleware Middleware)
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

func (t *tgbot) AddExactTextRoute(pattern string, handler port.Handler) {
	t.bot.RegisterHandler(
		bot.HandlerTypeMessageText,
		pattern,
		bot.MatchTypeExact,
		t.wrapTextHandler(pattern, handler),
	)
}

func (t *tgbot) AddMiddleware(m Middleware) {
	t.middlewares = append(t.middlewares, m)
}

func (t *tgbot) wrapTextHandler(pattern string, h port.Handler) bot.HandlerFunc {
	t.applyMiddleware(h)

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		ctx, span := tracing.StartSpanName(ctx, pattern)
		defer span.End()

		if err := h(
			ctx,
			port.MessageInfo{
				MessageID: update.Message.ID,
				ChatID:    update.Message.Chat.ID,
				Text:      update.Message.Text,
			},
		); err != nil {
			t.logger.WithError(err).
				WithField("endpoint", pattern).
				Error("error while processing message")
		}
	}
}

func (t *tgbot) applyMiddleware(h port.Handler) port.Handler {
	for i := len(t.middlewares) - 1; i >= 0; i-- {
		h = t.middlewares[i](h)
	}

	return h
}
